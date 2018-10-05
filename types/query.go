package types

import (
	"net/url"

	"github.com/dyninc/qstring"
)

// -
// Pagination
// -

// PageSize represents a query page size parameter
type PageSize int

// PageSizeDefault controls the default page size of listings
const PageSizeDefault PageSize = 5000

// -
// Sorting
// -

// SortOrder represents a query result sort order
type SortOrder string

// SortColumn represents a column to sort results by
type SortColumn string

// SortAsc is the ascending sort order for listings
const SortAsc SortOrder = "asc"

// SortDesc is the descending sort order for listings
const SortDesc SortOrder = "desc"

// ByPlayers means the list will use the amount of players as a sort key
const ByPlayers SortColumn = "player"

// -
// Filtering
// -

// FilterAttribute represents a filter to apply to results
type FilterAttribute string

// FilterPassword filters out servers with passwords
const FilterPassword FilterAttribute = "password"

// FilterEmpty filters out empty servers
const FilterEmpty FilterAttribute = "empty"

// FilterFull filters out full servers
const FilterFull FilterAttribute = "full"

// -
// URL Query
// -

// ServerListParams represents the URL query parameters for server listing
type ServerListParams struct {
	Page     int
	PageSize PageSize
	Sort     SortOrder
	By       SortColumn
	Filters  []FilterAttribute
}

// Example returns an example of ServerListParams in url.Values format
func (slp ServerListParams) Example() (result url.Values) {
	// nolint
	result, err := qstring.Marshal(&ServerListParams{
		Page:     2,
		PageSize: 100,
		Sort:     SortAsc,
		By:       ByPlayers,
		Filters:  []FilterAttribute{FilterFull, FilterPassword},
	})
	if err != nil {
		panic(err)
	}
	return
}
