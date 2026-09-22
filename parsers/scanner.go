package parsers

import (
	"bufio"
	"io"
	"strings"

	"github.com/rkoval/alfred-aws-console-services-workflow/aliases"
)

const (
	eof            = rune(0)
	openAllKeyword = "OPEN_ALL"
)

// Scanner splits the raw query into words and whitespace runs, and recognises
// the aliases that give a word a special meaning.
type Scanner struct {
	reader *bufio.Reader
}

// NewScanner returns a scanner of reader.
func NewScanner(reader io.Reader) *Scanner {
	return &Scanner{reader: bufio.NewReader(reader)}
}

// Scan returns the next token, its literal value and whether whitespace
// follows it.
func (s *Scanner) Scan() (TokenType, string, bool) {
	ch := s.read()
	if ch == eof {
		return EOF, "", false
	}
	s.unread()
	if isWhitespace(ch) {
		return WHITESPACE, s.scanWhile(isWhitespace), true
	}
	return s.scanWord()
}

func (s *Scanner) scanWord() (TokenType, string, bool) {
	word := s.scanWhile(isNotWhitespace)
	hasTrailingWhitespace := s.nextIsWhitespace()
	tokenType, value := classifyWord(word)
	return tokenType, value, hasTrailingWhitespace
}

// classifyWord recognises the aliases and keywords, stripping an alias from
// the value it introduces.
func classifyWord(word string) (TokenType, string) {
	if value, found := strings.CutPrefix(word, aliases.Search); found {
		return SEARCH_ALIAS, value
	}
	if value, found := strings.CutPrefix(word, aliases.OverrideAwsRegion); found {
		return REGION_OVERRIDE, value
	}
	if value, found := strings.CutPrefix(word, aliases.OverrideAwsProfile); found {
		return PROFILE_OVERRIDE, value
	}
	if word == openAllKeyword {
		return OPEN_ALL, word
	}
	return WORD, word
}

// scanWhile consumes the runes satisfying accept and returns them.
func (s *Scanner) scanWhile(accept func(rune) bool) string {
	var literal strings.Builder
	for {
		ch := s.read()
		if ch == eof {
			return literal.String()
		}
		if !accept(ch) {
			s.unread()
			return literal.String()
		}
		literal.WriteRune(ch)
	}
}

func (s *Scanner) nextIsWhitespace() bool {
	ch := s.read()
	if ch == eof {
		return false
	}
	s.unread()
	return isWhitespace(ch)
}

// read returns the next rune, or eof once the input is exhausted.
func (s *Scanner) read() rune {
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	_ = s.reader.UnreadRune()
}

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isNotWhitespace(ch rune) bool {
	return !isWhitespace(ch)
}
