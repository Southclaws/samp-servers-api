package main

import (
	"net/http"
)

// QueryRequest represents a GraphQL POST request payload
type QueryRequest struct {
	Query     string
	Operation string
	Variables map[string]interface{}
}

// QueryResponse represents a GraphQL response according to the HTTP spec
type QueryResponse struct {
	Data   map[string]interface{}
	Errors []string
}

// GraphQL handles GraphQL queries and conforms to http://graphql.org/learn/serving-over-http/
func (app *App) GraphQL(w http.ResponseWriter, r *http.Request) {

}
