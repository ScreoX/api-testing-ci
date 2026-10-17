package tests

import (
	"net/http"
	"testing"
)

const (
	defaultGraphQLURL    = "https://graphql.postman-echo.com/graphql"
	helloQuery           = `query Hello($person: PersonInput) { hello(person: $person) }`
	createPersonMutation = `mutation CreatePerson($person: PersonInput!) { createPerson(person: $person) { id name age } }`
	invalidFieldName     = "fieldThatDoesNotExist"
	invalidFieldQuery    = "query InvalidField { " + invalidFieldName + " }"
)

type graphQLRequest struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type graphQLError struct {
	Message string `json:"message"`
}

type graphQLResponse[T any] struct {
	Data   *T             `json:"data"`
	Errors []graphQLError `json:"errors"`
}

type graphQLPersonInput struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type helloData struct {
	Hello string `json:"hello"`
}

type graphQLPerson struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type createPersonData struct {
	CreatePerson graphQLPerson `json:"createPerson"`
}

func newGraphQLClient() *apiClient {
	return newAPIClient(envOrDefault("GRAPHQL_URL", defaultGraphQLURL))
}

func executeGraphQL(t *testing.T, query string, variables map[string]any) apiResponse {
	t.Helper()
	return newGraphQLClient().doJSON(t, http.MethodPost, "", graphQLRequest{
		Query:     query,
		Variables: variables,
	}, nil)
}

func requireNoGraphQLErrors(t *testing.T, errors []graphQLError) {
	t.Helper()
	if len(errors) != 0 {
		t.Fatalf("GraphQL response contains errors: %+v", errors)
	}
}
