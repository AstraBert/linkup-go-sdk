package linkup

import (
	"encoding/json"
	"errors"

	"github.com/invopop/jsonschema"
)

// Utility function to generate a JSON schema, provided a struct type
// as generic type
func GenerateJSONSchema[T any]() (json.RawMessage, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	schema.ID = ""
	schema.Version = ""
	serialized, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	return serialized, nil
}

// Utility function to get struct-typed results from a raw JSON output
func GetResultFromRawJSON[T any](jsonOutput string) (any, error) {
	var v T
	err := json.Unmarshal([]byte(jsonOutput), &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Utility function to get struct-typed results from the Data field of a StructuredWithSources output
func GetResultFromSourcedOutput[T any](sourcedOutput *StructuredWithSourcesDto) (any, error) {
	if sourcedOutput.Data != nil {
		data, err := json.Marshal(sourcedOutput.Data)
		if err != nil {
			return nil, err
		}
		var v T
		err = json.Unmarshal(data, &v)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
	return nil, errors.New("the Data field is null")
}
