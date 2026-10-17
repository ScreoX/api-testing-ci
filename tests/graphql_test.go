package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGraphQLQuery(t *testing.T) {
	const name = "Go Student"

	response := executeGraphQL(t, helloQuery, map[string]any{
		"person": graphQLPersonInput{Name: name, Age: 21},
	})
	requireStatus(t, response, http.StatusOK)

	result := decodeJSON[graphQLResponse[helloData]](t, response)
	requireNoGraphQLErrors(t, result.Errors)
	if result.Data == nil {
		t.Fatal("GraphQL query returned null data")
	}

	if result.Data.Hello != "Hello "+name {
		t.Fatalf("data.hello=%q, want %q", result.Data.Hello, "Hello "+name)
	}

}

func TestGraphQLMutation(t *testing.T) {
	name := fmt.Sprintf("Go Student %d", time.Now().UnixNano())
	const age = 23

	response := executeGraphQL(t, createPersonMutation, map[string]any{
		"person": graphQLPersonInput{Name: name, Age: age},
	})
	requireStatus(t, response, http.StatusOK)

	result := decodeJSON[graphQLResponse[createPersonData]](t, response)
	requireNoGraphQLErrors(t, result.Errors)
	if result.Data == nil {
		t.Fatal("GraphQL mutation returned null data")
	}

	created := result.Data.CreatePerson
	if created.ID == "" {
		t.Error("mutation returned an empty person id")
	}
	if created.Name != name {
		t.Errorf("mutation name=%q, want %q", created.Name, name)
	}
	if created.Age != age {
		t.Errorf("mutation age=%d, want %d", created.Age, age)
	}
}

func TestGraphQLValidationError(t *testing.T) {
	response := executeGraphQL(t, invalidFieldQuery, nil)
	if response.statusCode != http.StatusOK && response.statusCode != http.StatusBadRequest {
		t.Fatalf("validation status=%d, want 200 or 400; body=%s", response.statusCode, response.body)
	}

	result := decodeJSON[graphQLResponse[struct{}]](t, response)
	if len(result.Errors) == 0 || !strings.Contains(result.Errors[0].Message, invalidFieldName) {
		t.Fatalf("invalid query returned no error about the unknown field: %s", response.body)
	}
	if result.Data != nil {
		t.Fatalf("invalid query unexpectedly returned data: %+v", result.Data)
	}
}
