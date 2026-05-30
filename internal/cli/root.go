package cli

import (
	"io"
	"strings"

	"github.com/awuzag/opendart"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
)

// NewRootCommand builds the opendart CLI command tree.
func NewRootCommand(out io.Writer, errOut io.Writer, getenv getenvFunc) *cobra.Command {
	options := &rootOptions{
		baseURL: defaultBaseURL,
		output:  outputJSON,
		out:     out,
		errOut:  errOut,
		getenv:  getenv,
	}

	root := &cobra.Command{
		Use:           "opendart",
		Short:         "OpenDART API command line client",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return options.validateOutput()
		},
	}
	root.SetOut(out)
	root.SetErr(errOut)
	root.PersistentFlags().StringVar(&options.apiKey, "api-key", "", "OpenDART API key. Defaults to OPENDART_API_KEY.")
	root.PersistentFlags().StringVar(&options.baseURL, "base-url", defaultBaseURL, "OpenDART base URL.")
	root.PersistentFlags().StringVarP(&options.output, "output", "o", outputJSON, "Output format: json, raw, table, or csv.")

	addAPIVerbs(root, options)
	addHiddenDeprecatedCommand(root, newSummarizeCommand(options), "use `opendart get <business-resource> --view summary` instead")
	addHiddenDeprecatedCommand(root, newCompareCommand(options), "use `opendart get <business-resource> --corp-codes ...` instead")
	addHiddenDeprecatedCommand(root, newInspectCommand(options), "use `opendart get financial-metric --view source` instead")
	addLegacyAliases(root, options)

	return root
}

func (options *rootOptions) validateOutput() error {
	switch options.output {
	case outputJSON, outputRaw, outputTable, outputCSV:
		return nil
	default:
		return oops.In("opendart_cli").
			With("output", options.output).
			Errorf("opendart cli: unsupported output %q", options.output)
	}
}

func newCorpCodesCommand(options *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "corp-codes",
		Short: "DART 고유번호 목록을 조회합니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCorpCodes(cmd, options)
		},
	}
}

func newFinancialStatementCommand(options *rootOptions) *cobra.Command {
	var corpCode string
	var businessYear string
	var reportCode string

	cmd := &cobra.Command{
		Use:   "financial-statement",
		Short: "단일회사 주요계정 재무제표를 조회합니다.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(corpCode) == "" {
				return requiredFlagError("financial-statement", "corp-code")
			}
			if strings.TrimSpace(businessYear) == "" {
				return requiredFlagError("financial-statement", "business-year")
			}
			if strings.TrimSpace(reportCode) == "" {
				return requiredFlagError("financial-statement", "report-code")
			}

			client, err := newSDKClient(options)
			if err != nil {
				return err
			}
			result, err := client.FnlttSinglAcnt(cmd.Context(), opendart.FnlttSinglAcntParams{
				CorpCode:  corpCode,
				BsnsYear:  businessYear,
				ReprtCode: reportCode,
			})
			if err != nil {
				return err
			}
			return writeJSON(options.out, result.List)
		},
	}

	cmd.Flags().StringVar(&corpCode, "corp-code", "", "OpenDART corp_code.")
	cmd.Flags().StringVar(&businessYear, "business-year", "", "Business year, for example 2025.")
	cmd.Flags().StringVar(&reportCode, "report-code", "", "Report code, for example "+opendart.ReportCodeAnnual+".")
	return cmd
}

func runCorpCodes(cmd *cobra.Command, options *rootOptions) error {
	client, err := newSDKClient(options)
	if err != nil {
		return err
	}
	file, err := client.CorpCode(cmd.Context())
	if err != nil {
		return err
	}
	if options.output == outputRaw {
		_, err := options.out.Write(file.Body)
		return err
	}
	codes, err := decodeCorpCodeZIP(file.Body)
	if err != nil {
		return err
	}
	return writeJSON(options.out, codes)
}

func addAPIVerbs(root *cobra.Command, options *rootOptions) {
	verbs := map[string]*cobra.Command{}
	for _, spec := range apiCatalog {
		verb, ok := verbs[spec.Verb]
		if !ok {
			verb = &cobra.Command{
				Use:   spec.Verb,
				Short: verbDescription(spec.Verb),
			}
			verbs[spec.Verb] = verb
			root.AddCommand(verb)
		}
		verb.AddCommand(newGenericAPICommand(options, spec, spec.Resource))
	}
	if getCommand := verbs["get"]; getCommand != nil {
		addFinancialViewCommands(getCommand, options)
	}
}

func addHiddenDeprecatedCommand(root *cobra.Command, cmd *cobra.Command, message string) {
	cmd.Hidden = true
	cmd.Deprecated = message
	root.AddCommand(cmd)
}

func addLegacyAliases(root *cobra.Command, options *rootOptions) {
	corpCodes := newCorpCodesCommand(options)
	corpCodes.Hidden = true
	root.AddCommand(corpCodes)

	financialStatement := newFinancialStatementCommand(options)
	financialStatement.Hidden = true
	root.AddCommand(financialStatement)

	groups := map[string]*cobra.Command{}
	for _, spec := range apiCatalog {
		group, ok := groups[spec.Group]
		if !ok {
			group = &cobra.Command{
				Use:    spec.Group,
				Short:  groupDescription(spec.Group),
				Hidden: true,
			}
			groups[spec.Group] = group
			root.AddCommand(group)
		}
		legacy := newGenericAPICommand(options, spec, spec.Command)
		legacy.Hidden = true
		group.AddCommand(legacy)
	}
}

func newGenericAPICommand(options *rootOptions, spec apiSpec, use string) *cobra.Command {
	values := make(map[string]*string, len(spec.Params))
	cmd := &cobra.Command{
		Use:   use,
		Short: spec.Name,
		RunE: func(cmd *cobra.Command, args []string) error {
			if spec.OperationID == "corpCode" {
				return runCorpCodes(cmd, options)
			}
			for _, param := range spec.Params {
				value := values[param.Name]
				if param.Required && (value == nil || strings.TrimSpace(*value) == "") {
					return requiredFlagError(commandLabel(spec), flagName(param.Name))
				}
			}

			body, contentType, err := requestGeneric(cmd.Context(), options, spec, derefValues(values))
			if err != nil {
				return err
			}
			if strings.HasSuffix(spec.Endpoint, ".json") {
				return writeRawJSON(options.out, body)
			}
			if options.output == outputRaw {
				_, err := options.out.Write(body)
				return err
			}
			return writeBinaryJSON(options.out, spec, contentType, body)
		},
	}

	for _, param := range spec.Params {
		name := param.Name
		description := param.Description
		value := ""
		values[name] = &value
		cmd.Flags().StringVar(values[name], flagName(name), "", description)
	}
	return cmd
}

func commandLabel(spec apiSpec) string {
	if spec.Verb != "" && spec.Resource != "" {
		return spec.Verb + " " + spec.Resource
	}
	if spec.Command != "" {
		return spec.Command
	}
	return spec.OperationID
}

func requiredFlagError(command string, flag string) error {
	return oops.In("opendart_cli").
		With("command", command, "flag", flag).
		Errorf("opendart cli: --%s is required", flag)
}

func derefValues(values map[string]*string) map[string]string {
	result := make(map[string]string, len(values))
	for name, value := range values {
		if value != nil {
			result[name] = *value
		}
	}
	return result
}

func flagName(name string) string {
	return strings.ReplaceAll(name, "_", "-")
}

func groupDescription(group string) string {
	switch group {
	case "disclosure":
		return "공시정보 API"
	case "company":
		return "정기보고서 주요정보 API"
	case "financial":
		return "정기보고서 재무정보 API"
	case "ownership":
		return "지분공시 종합정보 API"
	case "material":
		return "주요사항보고서 주요정보 API"
	case "registration":
		return "증권신고서 주요정보 API"
	default:
		return "OpenDART API"
	}
}

func verbDescription(verb string) string {
	switch verb {
	case "search":
		return "공시와 공개 자료를 검색합니다."
	case "get":
		return "OpenDART JSON API 원문 응답을 조회합니다."
	case "list":
		return "마스터 목록을 조회합니다."
	case "download":
		return "OpenDART 파일 응답을 다운로드합니다."
	default:
		return "OpenDART API"
	}
}
