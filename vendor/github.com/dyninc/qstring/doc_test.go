package qstring_test

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/dyninc/qstring"
)

func ExampleUnmarshal() {
	// Query is the http request query struct.
	type Query struct {
		Names []string
		Limit int
		Page  int
	}

	query := &Query{}
	qValues, _ := url.ParseQuery("names=foo&names=bar&limit=50&page=1")
	err := qstring.Unmarshal(qValues, query)
	if err != nil {
		panic("Unable to Parse Query String")
	}

	os.Stdout.Write([]byte(fmt.Sprintf("%+v", query)))
	// Output: &{Names:[foo bar] Limit:50 Page:1}
}

func ExampleMarshalString() {
	// Query is the http request query struct.
	type Query struct {
		Names []string
		Limit int
		Page  int
	}

	query := &Query{
		Names: []string{"foo", "bar"},
		Limit: 50,
		Page:  1,
	}
	q, _ := qstring.MarshalString(query)
	os.Stdout.Write([]byte(q))
	// Output: limit=50&names=foo&names=bar&page=1
}

func ExampleUnmarshal_complex() {
	// PagingParams represents common pagination information for query strings
	type PagingParams struct {
		Page  int `qstring:"page"`
		Limit int `qstring:"limit"`
	}

	// Query is the http request query struct.
	type Query struct {
		Names    []string
		IDs      []int
		PageInfo *PagingParams
		Created  time.Time
	}
	query := &Query{}
	qValues, _ := url.ParseQuery("names=foo&names=bar&limit=50&page=1&ids=1&ids=2&created=2006-01-02T15:04:05Z")
	err := qstring.Unmarshal(qValues, query)
	if err != nil {
		panic("Unable to Parse Query String")
	}
}

func ExampleComparativeTime() {
	type DateQuery struct {
		Created  qstring.ComparativeTime
		Modified qstring.ComparativeTime
	}

	var query DateQuery
	qValues, _ := url.ParseQuery("created=>=2006-01-02T15:04:05Z&modified=<=2016-01-01T15:04Z")
	err := qstring.Unmarshal(qValues, &query)
	if err != nil {
		panic("Unable to Parse Query String")
	}
}
