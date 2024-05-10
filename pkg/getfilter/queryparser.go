package getfilter

var QueryParserSymbols = []string{
	"=",
	"<",
	">",
	"<=",
	">=",
	"in",
	"like",
	"is",
}

var QueryParserOperators = map[string]string{
	"eq":   "=",
	"lt":   "<",
	"gt":   ">",
	"lte":  "<=",
	"gte":  ">=",
	"in":   "in",
	"like": "like",
	"is":   "is",
}
