package linkup

import (
	"encoding/json"

	"github.com/invopop/jsonschema"
)

// Utility function to generate a JSON schema, provided a struct type
// as generic type
func GenerateJSONSchema[T any]() (json.RawMessage, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            false,
	}
	var v T
	schema := reflector.Reflect(v)
	serialized, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	return serialized, nil
}

func GetResultFromRawJSON[T any](jsonOutput string) (any, error) {
	var v T
	err := json.Unmarshal([]byte(jsonOutput), &v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
