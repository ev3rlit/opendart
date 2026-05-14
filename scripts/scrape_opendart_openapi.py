#!/usr/bin/env python3
"""Scrape the official OpenDART guide into an OpenAPI document.

The OpenDART guide is published as static HTML tables. This script intentionally
uses only the Python standard library so the generated API contract can be
refreshed without adding repository tooling dependencies.
"""

from __future__ import annotations

import argparse
import json
import re
import shutil
from dataclasses import dataclass
from datetime import date
from html.parser import HTMLParser
from pathlib import Path
from typing import Any
from urllib.parse import urljoin, urlparse
from urllib.request import urlopen


BASE_URL = "https://opendart.fss.or.kr"
GROUPS = [
    ("DS001", "공시정보"),
    ("DS002", "정기보고서 주요정보"),
    ("DS003", "정기보고서 재무정보"),
    ("DS004", "지분공시 종합정보"),
    ("DS005", "주요사항보고서 주요정보"),
    ("DS006", "증권신고서 주요정보"),
]


@dataclass(frozen=True)
class Link:
    href: str
    text: str


class GuideParser(HTMLParser):
    def __init__(self) -> None:
        super().__init__()
        self.tables: list[dict[str, Any]] = []
        self.links: list[Link] = []
        self._table: dict[str, Any] | None = None
        self._row: list[str] | None = None
        self._cell: list[str] | None = None
        self._caption: list[str] | None = None
        self._link_href: str | None = None
        self._link_text: list[str] | None = None

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        attrs_dict = dict(attrs)
        if tag == "a":
            self._link_href = attrs_dict.get("href")
            self._link_text = []
        if tag == "table":
            self._table = {"caption": "", "rows": []}
        elif self._table is not None and tag == "caption":
            self._caption = []
        elif self._table is not None and tag == "tr":
            self._row = []
        elif self._table is not None and tag in {"td", "th"}:
            self._cell = []
        elif self._cell is not None and tag == "br":
            self._cell.append("\n")

    def handle_endtag(self, tag: str) -> None:
        if tag == "a" and self._link_text is not None:
            text = compact("".join(self._link_text))
            if self._link_href:
                self.links.append(Link(self._link_href, text))
            self._link_href = None
            self._link_text = None
        if tag in {"td", "th"} and self._cell is not None:
            if self._row is not None:
                self._row.append(compact("".join(self._cell)))
            self._cell = None
        elif tag == "tr" and self._row is not None:
            if any(self._row) and self._table is not None:
                self._table["rows"].append(self._row)
            self._row = None
        elif tag == "caption" and self._caption is not None:
            if self._table is not None:
                self._table["caption"] = compact("".join(self._caption))
            self._caption = None
        elif tag == "table" and self._table is not None:
            self.tables.append(self._table)
            self._table = None

    def handle_data(self, data: str) -> None:
        if self._cell is not None:
            self._cell.append(data)
        elif self._caption is not None:
            self._caption.append(data)
        if self._link_text is not None:
            self._link_text.append(data)


def compact(value: str) -> str:
    return " ".join(value.replace("\r", " ").split())


def fetch(url: str) -> GuideParser:
    with urlopen(url, timeout=30) as response:
        html = response.read().decode("utf-8", "replace")
    parser = GuideParser()
    parser.feed(html)
    return parser


def table(parser: GuideParser, caption: str) -> dict[str, Any] | None:
    for item in parser.tables:
        if item["caption"] == caption:
            return item
    return None


def parse_main_group(group_code: str, group_name: str) -> list[dict[str, str]]:
    url = f"{BASE_URL}/guide/main.do?apiGrpCd={group_code}"
    parser = fetch(url)
    detail_links = [
        urljoin(BASE_URL, link.href)
        for link in parser.links
        if "/guide/detail.do" in link.href and f"apiGrpCd={group_code}" in link.href
    ]
    seen: set[str] = set()
    unique_links: list[str] = []
    for link in detail_links:
        api_id = api_id_from_url(link)
        if api_id not in seen:
            seen.add(api_id)
            unique_links.append(link)

    main_table = next((item for item in parser.tables if item["rows"] and item["rows"][0][:2] == ["번호", "API명"]), None)
    rows = main_table["rows"][1:] if main_table else []

    apis: list[dict[str, str]] = []
    for index, detail_url in enumerate(unique_links):
        api_id = api_id_from_url(detail_url)
        row = rows[index] if index < len(rows) else []
        apis.append(
            {
                "group_code": group_code,
                "group_name": group_name,
                "api_id": api_id,
                "name": row[1] if len(row) > 1 else api_id,
                "summary": row[2] if len(row) > 2 else "",
                "guide_url": detail_url,
                "group_url": url,
            }
        )
    return apis


def api_id_from_url(url: str) -> str:
    match = re.search(r"apiId=(\d+)", url)
    if not match:
        raise ValueError(f"missing apiId in URL: {url}")
    return match.group(1)


def parse_detail(api: dict[str, str]) -> dict[str, Any]:
    parser = fetch(api["guide_url"])
    basic = table(parser, "기본 정보")
    request = table(parser, "요청 인자")
    response = table(parser, "응답 결과")

    endpoints = []
    if basic:
        for row in basic["rows"][1:]:
            if len(row) >= 4:
                endpoints.append(
                    {
                        "method": row[0],
                        "url": row[1],
                        "encoding": row[2],
                        "format": row[3],
                    }
                )

    params = []
    if request:
        for row in request["rows"][1:]:
            if len(row) >= 5:
                params.append(
                    {
                        "name": row[0],
                        "title": row[1],
                        "type": row[2],
                        "required": row[3] == "Y",
                        "description": row[4],
                    }
                )

    response_rows = []
    if response:
        for row in response["rows"][1:]:
            if len(row) >= 3:
                response_rows.append({"name": row[0], "title": row[1], "description": row[2]})

    return {**api, "endpoints": endpoints, "params": params, "response_rows": response_rows}


def schema_name(endpoint: str, suffix: str) -> str:
    words = endpoint_words(endpoint)
    if not words:
        words = [Path(endpoint).stem]
    return "".join(word[:1].upper() + word[1:] for word in words) + suffix


def endpoint_slug(endpoint: str) -> str:
    words = slug_words(endpoint_words(endpoint))
    return "-".join(word.lower() for word in words) if words else Path(endpoint).stem.lower()


def operation_id(endpoint: str) -> str:
    words = endpoint_words(endpoint)
    if not words:
        return Path(endpoint).stem
    return words[0].lower() + "".join(word[:1].upper() + word[1:] for word in words[1:])


def endpoint_words(endpoint: str) -> list[str]:
    stem = Path(endpoint).stem
    return re.findall(r"[A-Z]?[a-z]+|[A-Z]+(?=\d|[A-Z]|$)|\d+", stem)


def slug_words(words: list[str]) -> list[str]:
    result: list[str] = []
    index = 0
    while index < len(words):
        if index + 1 < len(words) and words[index].isalpha() and words[index + 1].isdigit() and len(words[index]) == 1:
            result.append(words[index] + words[index + 1])
            index += 2
            continue
        result.append(words[index])
        index += 1
    return result


def endpoint_path(url: str) -> str:
    return urlparse(url).path


def json_endpoint(api: dict[str, Any]) -> dict[str, str] | None:
    for endpoint in api["endpoints"]:
        if endpoint["format"].upper() == "JSON":
            return endpoint
    return api["endpoints"][0] if api["endpoints"] else None


def parameter_schema(raw_type: str) -> dict[str, Any]:
    match = re.fullmatch(r"STRING\((\d+)\)", raw_type)
    if match:
        return {"type": "string", "maxLength": int(match.group(1))}
    return {"type": "string"}


def response_field_schema(field: dict[str, str]) -> dict[str, Any]:
    schema: dict[str, Any] = {"type": "string"}
    if field["title"]:
        schema["title"] = field["title"]
    if field["description"]:
        schema["description"] = field["description"]
    return schema


def properties_from_fields(fields: list[dict[str, str]]) -> dict[str, Any]:
    return {field["name"]: response_field_schema(field) for field in fields if field["name"]}


def split_response_rows(rows: list[dict[str, str]]) -> dict[str, Any]:
    top_fields: list[dict[str, str]] = []
    list_fields: list[dict[str, str]] = []
    groups: list[dict[str, Any]] = []
    mode = "top"
    current_group: dict[str, Any] | None = None

    for row in rows:
        name = row["name"]
        if name == "result":
            continue
        if name == "group":
            current_group = {"title": "", "fields": []}
            groups.append(current_group)
            mode = "group"
            continue
        if name == "list":
            mode = "group_list" if current_group is not None else "list"
            continue
        if name == "title" and current_group is not None and not current_group["title"]:
            current_group["title"] = row["description"] or row["title"]
            continue
        if mode == "list":
            list_fields.append(row)
        elif mode == "group_list" and current_group is not None:
            current_group["fields"].append(row)
        else:
            top_fields.append(row)

    return {"top_fields": top_fields, "list_fields": list_fields, "groups": groups}


def add_response_components(components: dict[str, Any], api: dict[str, Any], path: str) -> str:
    split = split_response_rows(api["response_rows"])
    response_name = schema_name(path, "Response")
    item_name = schema_name(path, "Item")
    properties = properties_from_fields(split["top_fields"])
    required = [name for name in ("status", "message") if name in properties]

    if split["list_fields"]:
        components["schemas"][item_name] = {
            "type": "object",
            "additionalProperties": False,
            "properties": properties_from_fields(split["list_fields"]),
        }
        properties["list"] = {
            "type": "array",
            "items": {"$ref": f"#/components/schemas/{item_name}"},
        }

    if split["groups"]:
        group_refs = []
        groups_metadata = []
        for index, group in enumerate(split["groups"], start=1):
            group_item_name = f"{schema_name(path, 'Group')}{index}Item"
            group_name = f"{schema_name(path, 'Group')}{index}"
            components["schemas"][group_item_name] = {
                "type": "object",
                "additionalProperties": False,
                "properties": properties_from_fields(group["fields"]),
            }
            components["schemas"][group_name] = {
                "type": "object",
                "additionalProperties": False,
                "properties": {
                    "title": {"type": "string", "const": group["title"]} if group["title"] else {"type": "string"},
                    "list": {
                        "type": "array",
                        "items": {"$ref": f"#/components/schemas/{group_item_name}"},
                    },
                },
            }
            group_refs.append({"$ref": f"#/components/schemas/{group_name}"})
            groups_metadata.append({"title": group["title"], "schema": group_item_name})
        properties["group"] = {"type": "array", "items": {"oneOf": group_refs}}
        components["x-opendart-groups"][response_name] = groups_metadata

    components["schemas"][response_name] = {
        "type": "object",
        "additionalProperties": False,
        "required": required,
        "properties": properties,
    }
    return response_name


def build_operation(
    api: dict[str, Any],
    endpoint: dict[str, str],
    components: dict[str, Any],
    schema_ref_prefix: str = "#/components/schemas/",
) -> dict[str, Any]:
    path = endpoint_path(endpoint["url"])
    response_schema_name = add_response_components(components, api, path)
    parameters = []
    for param in api["params"]:
        if param["name"] == "crtfc_key":
            continue
        parameters.append(
            {
                "name": param["name"],
                "in": "query",
                "required": param["required"],
                "description": param["description"],
                "schema": parameter_schema(param["type"]),
                "x-opendart-title": param["title"],
                "x-opendart-type": param["type"],
            }
        )

    content_type = "application/json" if endpoint["format"].upper() == "JSON" else "application/octet-stream"
    schema: dict[str, Any] = {"$ref": f"{schema_ref_prefix}{response_schema_name}"}
    if content_type != "application/json":
        schema = {"type": "string", "format": "binary"}

    return {
        "operationId": operation_id(path),
        "summary": api["name"],
        "description": api["summary"],
        "tags": [api["group_name"]],
        "security": [{"OpenDartApiKey": []}],
        "parameters": parameters,
        "responses": {
            "200": {
                "description": "OpenDART response. Business errors are returned with 200 and status/message fields for JSON APIs.",
                "content": {content_type: {"schema": schema}},
            }
        },
        "externalDocs": {"url": api["guide_url"]},
        "x-opendart-api-id": api["api_id"],
        "x-opendart-group-code": api["group_code"],
    }


def root_components() -> dict[str, Any]:
    return {
        "securitySchemes": {
            "OpenDartApiKey": {
                "type": "apiKey",
                "in": "query",
                "name": "crtfc_key",
                "description": "OpenDART API 인증키",
            }
        }
    }


def schema_components() -> dict[str, Any]:
    return {"schemas": {}, "x-opendart-groups": {}}


def build_bundled_openapi(apis: list[dict[str, Any]]) -> dict[str, Any]:
    components: dict[str, Any] = {
        **root_components(),
        "schemas": {},
        "x-opendart-groups": {},
    }
    paths: dict[str, Any] = {}
    for api in apis:
        endpoint = json_endpoint(api)
        if endpoint is None:
            continue
        path = endpoint_path(endpoint["url"])
        paths.setdefault(path, {})["get"] = build_operation(api, endpoint, components)

    return openapi_base(paths, components)


def build_split_openapi(apis: list[dict[str, Any]], split_dir: Path) -> dict[str, Any]:
    paths: dict[str, Any] = {}
    files: list[dict[str, str]] = []
    if split_dir.exists():
        shutil.rmtree(split_dir)
    split_dir.mkdir(parents=True, exist_ok=True)

    for api in apis:
        endpoint = json_endpoint(api)
        if endpoint is None:
            continue
        path = endpoint_path(endpoint["url"])
        filename = f"{api['api_id']}-{endpoint_slug(path)}.json"
        path_file = split_dir / filename
        components = schema_components()
        operation = build_operation(api, endpoint, components)
        path_doc = {
            "path": {"get": operation},
            "components": components,
        }
        path_file.write_text(json.dumps(path_doc, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
        paths[path] = {"$ref": f"./openapi/apis/{filename}#/path"}
        files.append({"api_id": api["api_id"], "path": path, "file": f"openapi/apis/{filename}"})

    openapi = openapi_base(paths, root_components())
    openapi["x-opendart-split-files"] = files
    return openapi


def openapi_base(paths: dict[str, Any], components: dict[str, Any]) -> dict[str, Any]:
    return {
        "openapi": "3.1.0",
        "info": {
            "title": "OpenDART API",
            "version": date.today().isoformat(),
            "description": "Generated from the official OpenDART development guide HTML tables.",
            "x-generated-from": [f"{BASE_URL}/guide/main.do?apiGrpCd={code}" for code, _ in GROUPS],
        },
        "servers": [{"url": BASE_URL}],
        "paths": paths,
        "components": components,
        "tags": [{"name": name, "x-opendart-group-code": code} for code, name in GROUPS],
    }


def scrape() -> list[dict[str, Any]]:
    apis: list[dict[str, Any]] = []
    for group_code, group_name in GROUPS:
        for api in parse_main_group(group_code, group_name):
            apis.append(parse_detail(api))
    return apis


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", default="docs/apis/opendart.openapi.json")
    parser.add_argument("--split-dir", default="docs/apis/openapi/apis")
    parser.add_argument("--bundle-output", default="docs/apis/opendart.openapi.bundle.json")
    parser.add_argument("--metadata-output", default="docs/apis/opendart-api-metadata.json")
    args = parser.parse_args()

    apis = scrape()
    openapi = build_split_openapi(apis, Path(args.split_dir))

    Path(args.output).write_text(json.dumps(openapi, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    if args.bundle_output:
        bundled = build_bundled_openapi(apis)
        Path(args.bundle_output).write_text(json.dumps(bundled, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    Path(args.metadata_output).write_text(json.dumps(apis, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")
    print(f"wrote {args.output}: {len(openapi['paths'])} path refs from {len(apis)} guide pages")
    print(f"wrote {args.split_dir}: {len(openapi['paths'])} per-API OpenAPI path files")
    if args.bundle_output:
        print(f"wrote {args.bundle_output}: bundled OpenAPI document")
    print(f"wrote {args.metadata_output}: raw scraped guide tables")


if __name__ == "__main__":
    main()
