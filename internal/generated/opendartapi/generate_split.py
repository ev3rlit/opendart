#!/usr/bin/env python3
import json
import re
import subprocess
import tempfile
from dataclasses import dataclass
from pathlib import Path


API_PACKAGE = "opendartapi"
SCHEMA_PACKAGE = "opendartschema"
OAPI_CODEGEN = "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.7.0"
SCHEMA_IMPORT = "github.com/ev3rlit/opendart/internal/generated/opendartschema"
ROOT_METHOD_OVERRIDES = {
    "company": "CompanyRaw",
    "document": "DocumentRaw",
}


@dataclass(frozen=True)
class CLIParamDef:
    name: str
    description: str
    required: bool


@dataclass(frozen=True)
class APIDef:
    api_id: str
    group_code: str
    path: str
    method: str
    operation_id: str
    summary: str
    response_content_type: str
    response_type: str
    params_type: str
    cli_params: list[CLIParamDef]
    file_stem: str


def snake_case(value: str) -> str:
    value = re.sub(r"[^0-9A-Za-z]+", "_", value)
    value = re.sub(r"([a-z0-9])([A-Z])", r"\1_\2", value)
    value = re.sub(r"_+", "_", value).strip("_").lower()
    return value or "api"


def kebab_case(value: str) -> str:
    return snake_case(value).replace("_", "-")


def load_json(path: Path) -> dict:
    return json.loads(path.read_text(encoding="utf-8"))


def exported_name(value: str) -> str:
    parts = re.findall(r"[A-Z]?[a-z0-9]+|[A-Z]+(?=[A-Z]|$)", value)
    if not parts:
        return "API"
    return "".join(part[:1].upper() + part[1:] for part in parts)


def strip_security(path_item: dict) -> dict:
    result = json.loads(json.dumps(path_item))
    for operation in result.values():
        if isinstance(operation, dict):
            operation.pop("security", None)
    return result


def add_bun_tags(components: dict) -> dict:
    result = json.loads(json.dumps(components))
    for schema in result.get("schemas", {}).values():
        if not isinstance(schema, dict) or schema.get("type") != "object":
            continue
        for name, prop in schema.get("properties", {}).items():
            if not isinstance(prop, dict):
                continue
            if prop.get("type") in {"array", "object"}:
                continue
            extra_tags = prop.setdefault("x-oapi-codegen-extra-tags", {})
            extra_tags.setdefault("bun", name)
    return result


def response_info(operation: dict) -> tuple[str, str]:
    content = operation.get("responses", {}).get("200", {}).get("content", {})
    if not isinstance(content, dict) or not content:
        return "", ""
    if "application/json" in content:
        schema = content["application/json"].get("schema", {})
        ref = schema.get("$ref", "")
        return "application/json", ref.removeprefix("#/components/schemas/")
    content_type = next(iter(content.keys()))
    return content_type, ""


def params_type(operation_id: str, operation: dict) -> str:
    for param in operation.get("parameters", []):
        if not isinstance(param, dict):
            continue
        if param.get("name") == "crtfc_key":
            continue
        return f"{exported_name(operation_id)}Params"
    return ""


def cli_params(operation: dict) -> list[CLIParamDef]:
    params: list[CLIParamDef] = []
    for param in operation.get("parameters", []):
        if not isinstance(param, dict):
            continue
        name = param.get("name", "")
        if name == "crtfc_key":
            continue
        description = param.get("x-opendart-title") or param.get("description") or name
        params.append(
            CLIParamDef(
                name=name,
                description=description,
                required=bool(param.get("required", False)),
            )
        )
    return params


def api_defs(repo_root: Path, docs_dir: Path, index: dict) -> list[APIDef]:
    definitions: list[APIDef] = []
    for item in index["x-opendart-split-files"]:
        split_doc = load_json(docs_dir / item["file"])
        path_item = split_doc["path"]
        method, operation = next(
            (method, operation)
            for method, operation in path_item.items()
            if isinstance(operation, dict)
        )
        operation_id = operation.get("operationId", item["api_id"])
        response_content_type, response_type = response_info(operation)
        definitions.append(
            APIDef(
                api_id=item["api_id"],
                group_code=operation.get("x-opendart-group-code", ""),
                path=item["path"],
                method=method.upper(),
                operation_id=operation_id,
                summary=operation.get("summary", ""),
                response_content_type=response_content_type,
                response_type=response_type,
                params_type=params_type(operation_id, operation),
                cli_params=cli_params(operation),
                file_stem=f"{item['api_id']}_{snake_case(operation_id)}",
            )
        )
    return definitions


def load_cli_names(path: Path) -> dict[str, dict[str, str]]:
    sections: dict[str, dict[str, str]] = {}
    if not path.exists():
        return sections

    section = ""
    for raw_line in path.read_text(encoding="utf-8").splitlines():
        line = raw_line.split("#", 1)[0].rstrip()
        if not line.strip():
            continue
        if not line.startswith(" ") and line.endswith(":"):
            section = line[:-1].strip()
            sections.setdefault(section, {})
            continue
        if section == "" or not line.startswith("  "):
            continue
        key, separator, value = line.strip().partition(":")
        if separator == "":
            continue
        sections[section][key.strip()] = value.strip().strip('"').strip("'")
    return sections


def write_cli_catalog(repo_root: Path, apis: list[APIDef], cli_names: dict[str, dict[str, str]]) -> None:
    output = repo_root / "internal" / "cli" / "catalog_gen.go"
    groups = cli_names.get("groups", {})
    commands = cli_names.get("commands", {})
    lines = [
        "// Code generated by internal/generated/opendartapi/generate_split.py; DO NOT EDIT.",
        "",
        "package cli",
        "",
        "var apiCatalog = []apiSpec{",
    ]
    for api in apis:
        lines.extend(
            [
                "\t{",
                f'\t\tGroup:       "{groups.get(api.group_code, kebab_case(api.group_code))}",',
                f'\t\tCommand:     "{commands.get(api.operation_id, kebab_case(api.operation_id))}",',
                f'\t\tAPIID:       "{api.api_id}",',
                f'\t\tOperationID: "{api.operation_id}",',
                f"\t\tName:        {json.dumps(api.summary, ensure_ascii=False)},",
                f'\t\tEndpoint:    "{api.path}",',
            ]
        )
        if api.cli_params:
            lines.append("\t\tParams: []paramSpec{")
            for param in api.cli_params:
                lines.append(
                    "\t\t\t{"
                    f'Name: "{param.name}", '
                    f"Description: {json.dumps(param.description, ensure_ascii=False)}, "
                    f"Required: {str(param.required).lower()}"
                    "},"
                )
            lines.append("\t\t},")
        lines.extend(["\t},"])
    lines.extend(["}", ""])
    output.write_text("\n".join(lines), encoding="utf-8")


def generated_type_names(package_dir: Path) -> list[str]:
    names: list[str] = []
    for path in sorted(package_dir.glob("*.gen.go")):
        for line in path.read_text(encoding="utf-8").splitlines():
            match = re.match(r"type ([A-Z][A-Za-z0-9_]*)\b", line)
            if match:
                names.append(match.group(1))
    if len(names) != len(set(names)):
        duplicates = sorted({name for name in names if names.count(name) > 1})
        raise RuntimeError(f"duplicate generated type names: {duplicates}")
    return names


def existing_client_methods(repo_root: Path) -> set[str]:
    names: set[str] = set()
    pattern = re.compile(r"func \(client \*Client\) ([A-Z][A-Za-z0-9_]*)\(")
    for path in sorted(repo_root.glob("*.go")):
        if path.name == "api_methods_gen.go":
            continue
        for line in path.read_text(encoding="utf-8").splitlines():
            match = pattern.search(line)
            if match:
                names.add(match.group(1))
    return names


def write_root_aliases(repo_root: Path, names: list[str]) -> None:
    output = repo_root / "api_types_gen.go"
    lines = [
        "// Code generated by internal/generated/opendartapi/generate_split.py; DO NOT EDIT.",
        "",
        "package opendart",
        "",
        'import opendartschema "github.com/ev3rlit/opendart/internal/generated/opendartschema"',
        "",
    ]
    for name in names:
        lines.extend(
            [
                f"// {name} is a generated OpenDART API schema type.",
                f"type {name} = opendartschema.{name}",
                "",
            ]
        )
    output.write_text("\n".join(lines), encoding="utf-8")


def write_api_support(api_dir: Path) -> None:
    output = api_dir / "types.gen.go"
    output.write_text(
        "\n".join(
            [
                "// Code generated by internal/generated/opendartapi/generate_split.py; DO NOT EDIT.",
                "",
                f"package {API_PACKAGE}",
                "",
                'import "context"',
                "",
                "// APISpec describes a generated OpenDART API operation.",
                "type APISpec struct {",
                "\tAPIID string",
                "\tOperationID string",
                "\tMethod string",
                "\tPath string",
                "\tSummary string",
                "\tResponseContentType string",
                "}",
                "",
                "// FileResponse contains bytes returned by generated OpenDART file APIs.",
                "type FileResponse struct {",
                "\tContentType string",
                "\tBody []byte",
                "}",
                "",
                "// Caller executes generated OpenDART API operations.",
                "type Caller interface {",
                "\tCallOpenDARTAPI(ctx context.Context, spec APISpec, params any, out any) error",
                "\tCallOpenDARTFile(ctx context.Context, spec APISpec, params any) (*FileResponse, error)",
                "}",
                "",
            ]
        ),
        encoding="utf-8",
    )


def write_api_function(api_dir: Path, api: APIDef) -> None:
    output = api_dir / f"{api.file_stem}.gen.go"
    name = exported_name(api.operation_id)
    imports = ['"context"']
    if api.params_type or api.response_type:
        imports.append(f'opendartschema "{SCHEMA_IMPORT}"')
    lines = [
        "// Code generated by internal/generated/opendartapi/generate_split.py; DO NOT EDIT.",
        "",
        f"package {API_PACKAGE}",
        "",
        "import (",
    ]
    lines.extend(f"\t{item}" for item in imports)
    lines.extend(
        [
            ")",
            "",
            f"// {name}Spec describes the generated {api.operation_id} API operation.",
            f"var {name}Spec = APISpec{{",
            f'\tAPIID: "{api.api_id}",',
            f'\tOperationID: "{api.operation_id}",',
            f'\tMethod: "{api.method}",',
            f'\tPath: "{api.path}",',
            f'\tSummary: {json.dumps(api.summary, ensure_ascii=False)},',
            f'\tResponseContentType: "{api.response_content_type}",',
            "}",
            "",
        ]
    )
    params_decl = ""
    params_arg = "nil"
    if api.params_type:
        params_decl = f", params opendartschema.{api.params_type}"
        params_arg = "params"

    if api.response_type:
        lines.extend(
            [
                f"// {name} calls the generated OpenDART {api.operation_id} API operation.",
                f"func {name}(ctx context.Context, caller Caller{params_decl}) (*opendartschema.{api.response_type}, error) {{",
                f"\tvar result opendartschema.{api.response_type}",
                f"\tif err := caller.CallOpenDARTAPI(ctx, {name}Spec, {params_arg}, &result); err != nil {{",
                "\t\treturn nil, err",
                "\t}",
                "\treturn &result, nil",
                "}",
                "",
            ]
        )
    else:
        lines.extend(
            [
                f"// {name} calls the generated OpenDART {api.operation_id} file API operation.",
                f"func {name}(ctx context.Context, caller Caller{params_decl}) (*FileResponse, error) {{",
                f"\treturn caller.CallOpenDARTFile(ctx, {name}Spec, {params_arg})",
                "}",
                "",
            ]
        )
    output.write_text("\n".join(lines), encoding="utf-8")


def root_method_name(api: APIDef, existing_methods: set[str], emitted: set[str]) -> str:
    if api.operation_id in ROOT_METHOD_OVERRIDES:
        name = ROOT_METHOD_OVERRIDES[api.operation_id]
        if name in existing_methods or name in emitted:
            raise RuntimeError(f"duplicate root method override: {name}")
        emitted.add(name)
        return name

    base = exported_name(api.operation_id)
    if base not in existing_methods and base not in emitted:
        emitted.add(base)
        return base
    name = f"{base}Raw"
    if name not in existing_methods and name not in emitted:
        emitted.add(name)
        return name
    index = 2
    while f"{name}{index}" in existing_methods or f"{name}{index}" in emitted:
        index += 1
    emitted.add(f"{name}{index}")
    return f"{name}{index}"


def write_root_methods(repo_root: Path, apis: list[APIDef]) -> None:
    output = repo_root / "api_methods_gen.go"
    existing_methods = existing_client_methods(repo_root)
    emitted: set[str] = set()
    lines = [
        "// Code generated by internal/generated/opendartapi/generate_split.py; DO NOT EDIT.",
        "",
        "package opendart",
        "",
        "import (",
        '\t"context"',
        "",
        '\topendartapi "github.com/ev3rlit/opendart/internal/generated/opendartapi"',
        ")",
        "",
    ]
    for api in apis:
        api_func = exported_name(api.operation_id)
        method_name = root_method_name(api, existing_methods, emitted)
        params_decl = ""
        params_arg = ""
        if api.params_type:
            params_decl = f", params {api.params_type}"
            params_arg = ", params"

        if method_name.endswith("Raw"):
            lines.append(
                f"// {method_name} calls the generated OpenDART {api.operation_id} API operation. The Raw suffix avoids a handwritten SDK method name collision."
            )
        else:
            lines.append(f"// {method_name} calls the generated OpenDART {api.operation_id} API operation.")

        if api.response_type:
            lines.extend(
                [
                    f"func (client *Client) {method_name}(ctx context.Context{params_decl}) (*{api.response_type}, error) {{",
                    f"\treturn opendartapi.{api_func}(ctx, client.apiCaller{params_arg})",
                    "}",
                    "",
                ]
            )
        else:
            lines.extend(
                [
                    f"func (client *Client) {method_name}(ctx context.Context{params_decl}) (*FileResponse, error) {{",
                    f"\tresult, err := opendartapi.{api_func}(ctx, client.apiCaller{params_arg})",
                    "\tif err != nil {",
                    "\t\treturn nil, err",
                    "\t}",
                    "\treturn &FileResponse{ContentType: result.ContentType, Body: result.Body}, nil",
                    "}",
                    "",
                ]
            )
    output.write_text("\n".join(lines), encoding="utf-8")


def gofmt_generated(repo_root: Path, api_dir: Path, schema_dir: Path) -> None:
    paths = [
        *(sorted(api_dir.glob("*.go"))),
        *(sorted(schema_dir.glob("*.go"))),
        repo_root / "api_types_gen.go",
        repo_root / "api_methods_gen.go",
        repo_root / "internal" / "cli" / "catalog_gen.go",
    ]
    subprocess.run(["gofmt", "-w", *(str(path) for path in paths)], cwd=repo_root, check=True)


def main() -> None:
    api_dir = Path(__file__).resolve().parent
    generated_dir = api_dir.parent
    schema_dir = generated_dir / SCHEMA_PACKAGE
    repo_root = api_dir.parents[2]
    docs_dir = repo_root / "docs" / "apis"
    index = load_json(docs_dir / "opendart.openapi.json")
    cli_names = load_cli_names(docs_dir / "cli-names.yaml")
    apis = api_defs(repo_root, docs_dir, index)

    schema_dir.mkdir(parents=True, exist_ok=True)
    for generated in schema_dir.glob("*.gen.go"):
        generated.unlink()
    for generated in api_dir.glob("*.gen.go"):
        generated.unlink()

    with tempfile.TemporaryDirectory(prefix="opendart-oapi-codegen-") as temp_name:
        temp_dir = Path(temp_name)
        for api, item in zip(apis, index["x-opendart-split-files"]):
            split_doc = load_json(docs_dir / item["file"])
            path_item = strip_security(split_doc["path"])
            output = schema_dir / f"{api.file_stem}.gen.go"
            spec = {
                "openapi": "3.0.3",
                "info": {
                    "title": f"OpenDART {api.operation_id}",
                    "version": index.get("info", {}).get("version", "0.0.0"),
                },
                "paths": {api.path: path_item},
                "components": add_bun_tags(split_doc.get("components", {})),
            }
            spec_path = temp_dir / f"{api.api_id}.json"
            config_path = temp_dir / f"{api.api_id}.yaml"
            spec_path.write_text(json.dumps(spec, ensure_ascii=False, indent=2), encoding="utf-8")
            config_path.write_text(
                "\n".join(
                    [
                        f"package: {SCHEMA_PACKAGE}",
                        f"output: {output}",
                        "generate:",
                        "  models: true",
                        "output-options:",
                        "  prefer-skip-optional-pointer: true",
                        "",
                    ]
                ),
                encoding="utf-8",
            )
            subprocess.run(
                ["go", "run", OAPI_CODEGEN, "-config", str(config_path), str(spec_path)],
                cwd=repo_root,
                check=True,
            )

    write_api_support(api_dir)
    for api in apis:
        write_api_function(api_dir, api)
    write_root_aliases(repo_root, generated_type_names(schema_dir))
    write_root_methods(repo_root, apis)
    write_cli_catalog(repo_root, apis, cli_names)
    gofmt_generated(repo_root, api_dir, schema_dir)


if __name__ == "__main__":
    main()
