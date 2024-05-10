package getfilter

import (
	"github.com/gin-gonic/gin"
)

type FilterItem struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type FilterSort struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type FilterMiddleware struct {
	ArrQuery []FilterItem
	ArrSort  []FilterSort
	Message  string
}

func Initiate(c *gin.Context) *FilterMiddleware {
	filtermiddleware := &FilterMiddleware{}

	filter, exists := c.GetQuery("filter")
	if exists {
		filterItems := make([]FilterItem, 0)
		for field, q := range ArrQuery(filter) {
			filterItems = append(filterItems, FilterItem{
				Field:    field,
				Operator: q.Operator,
				Value:    q.Value,
			})
		}
		filtermiddleware.ArrQuery = filterItems
	}
	sort, exists := c.GetQuery("sort")
	if exists {
		filterSort := make([]FilterSort, 0)
		for field, q := range ArrSort(sort) {
			filterSort = append(filterSort, FilterSort{
				Field: field,
				Value: q.Value,
			})
		}
		filtermiddleware.ArrSort = filterSort
	}

	return filtermiddleware
}

func Validation(query *FilterMiddleware) bool {
	for _, item := range query.ArrQuery {
		if !contains(item.Operator, QueryParserSymbols) {
			return false
		}
	}
	return true
}

func Handle(c *gin.Context) bool {
	if Validation(Initiate(c)) {
		c.Set("ArrQuery", Initiate(c).ArrQuery)
		c.Set("ArrSort", Initiate(c).ArrSort)
		return false
	}
	return true
}

func contains(item string, arr []string) bool {
	for _, v := range arr {
		if v == item {
			return true
		}
	}
	return false
}
