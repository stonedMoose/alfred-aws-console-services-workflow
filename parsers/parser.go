package parsers

import (
	"fmt"
	"strings"

	"github.com/rkoval/alfred-aws-console-services-workflow/awsconfig"
	"github.com/rkoval/alfred-aws-console-services-workflow/awsworkflow"
)

// Parser turns the raw Alfred query into a Query.
type Parser struct {
	rawQuery string
	scanner  *Scanner
}

// NewParser returns a parser of rawQuery.
func NewParser(rawQuery string) *Parser {
	return &Parser{
		rawQuery: rawQuery,
		scanner:  NewScanner(strings.NewReader(rawQuery)),
	}
}

// Parse reads the service catalogue at ymlPath and interprets the query
// against it.
func (p *Parser) Parse(ymlPath string) (*Query, []awsworkflow.AwsService) {
	awsServices := ParseConsoleServicesYml(ymlPath)
	tokens, endsWithWhitespace := p.scanIntoTokens()
	builder := newQueryBuilder(p.rawQuery, endsWithWhitespace, awsServices)
	for i, token := range tokens {
		builder.consume(token, i == len(tokens)-1)
	}
	return builder.build(), awsServices
}

// scanIntoTokens returns the tokens of the query, whitespace dropped, and
// whether the query ends with whitespace.
func (p *Parser) scanIntoTokens() ([]Token, bool) {
	var tokens []Token
	endsWithWhitespace := false
	for {
		tokenType, literal, hasTrailingWhitespace := p.scanner.Scan()
		if tokenType == EOF {
			return tokens, endsWithWhitespace
		}
		endsWithWhitespace = hasTrailingWhitespace
		if tokenType != WHITESPACE {
			tokens = append(tokens, Token{Type: tokenType, Value: literal})
		}
	}
}

// queryBuilder interprets the tokens one by one: words name a service, then a
// sub-service, and the rest of them is the search term; aliases select a
// region or a profile.
type queryBuilder struct {
	query              *Query
	awsServices        []awsworkflow.AwsService
	remainingQuery     string
	endsWithWhitespace bool
}

func newQueryBuilder(rawQuery string, endsWithWhitespace bool, awsServices []awsworkflow.AwsService) *queryBuilder {
	return &queryBuilder{
		query: &Query{
			RawQuery:              rawQuery,
			HasTrailingWhitespace: endsWithWhitespace,
		},
		awsServices:        awsServices,
		endsWithWhitespace: endsWithWhitespace,
	}
}

func (b *queryBuilder) consume(token Token, isLastToken bool) {
	switch token.Type {
	case WORD:
		b.consumeWord(token.Value)
	case OPEN_ALL:
		b.query.HasOpenAll = true
	case SEARCH_ALIAS:
		b.query.HasDefaultSearchAlias = true
		b.remainingQuery += token.Value
	case REGION_OVERRIDE:
		b.consumeRegionOverride(token.Value, isLastToken)
	case PROFILE_OVERRIDE:
		b.consumeProfileOverride(token.Value, isLastToken)
	default:
		panic(fmt.Errorf("no handler for token: %#v", token))
	}
}

// consumeWord appends the word to the search term, unless the words read so
// far name a service or, once a service is known, one of its sub-services.
func (b *queryBuilder) consumeWord(word string) {
	b.appendWord(word)
	switch {
	case b.query.Service == nil:
		if service := awsworkflow.FindServiceById(b.awsServices, b.remainingQuery); service != nil {
			b.query.Service = service
			b.remainingQuery = ""
		}
	case b.query.SubService == nil:
		if subService := awsworkflow.FindServiceById(b.query.Service.SubServices, b.remainingQuery); subService != nil {
			b.query.SubService = subService
			b.remainingQuery = ""
		}
	}
}

func (b *queryBuilder) appendWord(word string) {
	if b.remainingQuery != "" {
		b.remainingQuery += " "
	}
	b.remainingQuery += word
}

// consumeRegionOverride selects the region when it is known; a region still
// being typed, or unknown, turns the query into a region search instead.
func (b *queryBuilder) consumeRegionOverride(name string, isLastToken bool) {
	if b.query.ProfileQuery != nil {
		return
	}
	if region := awsconfig.FindRegion(name); region != nil {
		b.query.regionOverride = region
	}
	if b.query.regionOverride == nil || b.isBeingTyped(isLastToken) {
		b.query.RegionQuery = &name
	}
}

// consumeProfileOverride selects the profile when it is known; a profile
// still being typed, or unknown, turns the query into a profile search instead.
func (b *queryBuilder) consumeProfileOverride(name string, isLastToken bool) {
	if b.query.RegionQuery != nil {
		return
	}
	if profile := awsconfig.FindProfile(name); profile != nil {
		b.query.ProfileOverride = profile
	}
	if b.query.ProfileOverride == nil || b.isBeingTyped(isLastToken) {
		b.query.ProfileQuery = &name
	}
}

// isBeingTyped tells whether the token is the one the user is still typing:
// the last one, with no whitespace after it.
func (b *queryBuilder) isBeingTyped(isLastToken bool) bool {
	return isLastToken && !b.endsWithWhitespace
}

func (b *queryBuilder) build() *Query {
	b.query.RemainingQuery = b.remainingQuery
	return b.query
}
