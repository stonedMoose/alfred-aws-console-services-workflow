package parsers

// TokenType tells what a piece of the query is.
type TokenType int

const (
	ILLEGAL TokenType = iota

	EOF
	WHITESPACE
	WORD

	OPEN_ALL
	SEARCH_ALIAS
	REGION_OVERRIDE
	PROFILE_OVERRIDE
)

// Token is a piece of the query. For aliases, Value is what follows the alias.
type Token struct {
	Type  TokenType
	Value string
}
