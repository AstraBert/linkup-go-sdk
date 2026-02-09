package linkup

import (
	"bytes"
	"encoding/json"
	"os"
	"slices"
	"testing"
	"time"
)

type TestUser struct {
	ID          int            `json:"id"`
	Name        string         `json:"name" jsonschema:"title=the name,description=The name of a friend,example=joe,example=lucy,default=alex"`
	Friends     []int          `json:"friends,omitempty" jsonschema_description:"The list of IDs, omitted when empty"`
	Tags        map[string]any `json:"tags,omitempty" jsonschema_extras:"a=b,foo=bar,foo=bar1"`
	BirthDate   time.Time      `json:"birth_date,omitempty" jsonschema:"oneof_required=date"`
	YearOfBirth string         `json:"year_of_birth,omitempty" jsonschema:"oneof_required=year"`
	Metadata    any            `json:"metadata,omitempty" jsonschema:"oneof_type=string;array"`
	FavColor    string         `json:"fav_color,omitempty" jsonschema:"enum=red,enum=green,enum=blue"`
}

func getCorrectSchema() (string, error) {
	content, err := os.ReadFile("testfiles/schema.json")
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := json.Compact(&buf, content); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestGenerateJSONSchema(t *testing.T) {
	schema, err := GenerateJSONSchema[TestUser]()
	if err != nil {
		t.Fatalf("Unexpected error when generating JSON schema: %s", err.Error())
	}
	correctSchema, err := getCorrectSchema()
	if err != nil {
		t.Fatalf("Unexpected error when getting correct JSON schema: %s", err.Error())
	}
	if string(schema) != correctSchema {
		t.Fatalf("Expected JSON schema to be %s, got %s", correctSchema, schema)
	}
}

func TestGetResultFromRawJson(t *testing.T) {
	user := TestUser{
		ID:      1,
		Name:    "Lucy",
		Friends: []int{2, 3, 4},
	}
	serialized, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Unexpecting error while marshaling user: %s", err.Error())
	}
	val, err := GetResultFromRawJSON[TestUser](string(serialized))
	if err != nil {
		t.Fatalf("Unexpected error while getting user from serialized representation: %s", err.Error())
	}
	typedVal, ok := val.(TestUser)
	if !ok {
		t.Fatal("Expecting TestUser type, but type assertion was unsuccessfull.")
	}
	if typedVal.ID != user.ID || typedVal.Name != user.Name || !slices.Equal(typedVal.Friends, user.Friends) {
		t.Fatalf("Expecting function to yield %v as return value, got %v", user, typedVal)
	}
}

func TestGetResultFromSourcedOutput(t *testing.T) {
	user := TestUser{
		ID:      1,
		Name:    "Lucy",
		Friends: []int{2, 3, 4},
	}
	userMap := map[string]any{"id": 1, "name": "Lucy", "friends": []int{2, 3, 4}}
	sourcedOutput := &StructuredWithSourcesDto{
		Data:    userMap,
		Sources: nil,
	}
	val, err := GetResultFromSourcedOutput[TestUser](sourcedOutput)
	if err != nil {
		t.Fatalf("Unexpected error while getting user from serialized representation: %s", err.Error())
	}
	typedVal, ok := val.(TestUser)
	if !ok {
		t.Fatal("Expecting TestUser type, but type assertion was unsuccessfull.")
	}
	if typedVal.ID != user.ID || typedVal.Name != user.Name || !slices.Equal(typedVal.Friends, user.Friends) {
		t.Fatalf("Expecting function to yield %v as return value, got %v", user, typedVal)
	}
}

func TestGetResultFromSourcedOutputEmptyData(t *testing.T) {
	sourcedOutput := &StructuredWithSourcesDto{
		Data:    nil,
		Sources: nil,
	}
	_, err := GetResultFromSourcedOutput[TestUser](sourcedOutput)
	if err == nil {
		t.Fatal("Expecting an error, got none")
	} else {
		if err.Error() != "the Data field is null" {
			t.Fatalf("Unexpected error message: %s", err.Error())
		}
	}
}
