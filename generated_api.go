package opendart

import (
	"context"
	"reflect"
	"strconv"
	"strings"

	opendartapi "github.com/awuzag/opendart/internal/generated/opendartapi"
	"github.com/go-resty/resty/v2"
	"github.com/samber/oops"
)

type apiCallerConfig struct {
	apiKey string
	resty  *resty.Client
}

type generatedAPICaller struct {
	config apiCallerConfig
}

func (caller generatedAPICaller) CallOpenDARTAPI(ctx context.Context, spec opendartapi.APISpec, params any, out any) error {
	query, err := generatedQueryParams(params)
	if err != nil {
		return err
	}
	return getJSON(ctx, caller.config, spec.Path, query, spec.OperationID, endpointOp(spec.Path), out)
}

func (caller generatedAPICaller) CallOpenDARTFile(ctx context.Context, spec opendartapi.APISpec, params any) (*opendartapi.FileResponse, error) {
	query, err := generatedQueryParams(params)
	if err != nil {
		return nil, err
	}
	result, err := getFile(ctx, caller.config, spec.Path, query, spec.OperationID, endpointOp(spec.Path))
	if err != nil {
		return nil, err
	}
	return &opendartapi.FileResponse{ContentType: result.ContentType, Body: result.Body}, nil
}

func generatedQueryParams(params any) (map[string]string, error) {
	result := make(map[string]string)
	if params == nil {
		return result, nil
	}

	value := reflect.ValueOf(params)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return result, nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil, oops.In("generated_query").
			With("kind", value.Kind().String()).
			New("opendart: generated API params must be a struct")
	}

	valueType := value.Type()
	for index := range value.NumField() {
		field := valueType.Field(index)
		name := strings.Split(field.Tag.Get("form"), ",")[0]
		if name == "" || name == "-" {
			continue
		}
		text, ok, err := generatedFieldValue(value.Field(index))
		if err != nil {
			return nil, oops.In("generated_query").
				With("field", field.Name).
				Wrap(err)
		}
		if ok && strings.TrimSpace(text) != "" {
			result[name] = text
		}
	}
	return result, nil
}

func generatedFieldValue(value reflect.Value) (string, bool, error) {
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return "", false, nil
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.String:
		return value.String(), value.String() != "", nil
	case reflect.Bool:
		if !value.Bool() {
			return "", false, nil
		}
		return strconv.FormatBool(value.Bool()), true, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() == 0 {
			return "", false, nil
		}
		return strconv.FormatInt(value.Int(), 10), true, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if value.Uint() == 0 {
			return "", false, nil
		}
		return strconv.FormatUint(value.Uint(), 10), true, nil
	case reflect.Float32, reflect.Float64:
		if value.Float() == 0 {
			return "", false, nil
		}
		return strconv.FormatFloat(value.Float(), 'f', -1, value.Type().Bits()), true, nil
	default:
		return "", false, oops.In("generated_query").
			With("kind", value.Kind().String()).
			New("opendart: unsupported generated API param field type")
	}
}
