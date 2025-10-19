// Package n4l provides a modern N4L compiler with proper error reporting
// Treats N4L as a compiled language with lexer, parser, and semantic analysis
package n4l

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// TokenType represents different types of N4L tokens
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenChapter
	TokenContext
	TokenNode
	TokenArrow
	TokenAlias
	TokenString
	TokenDitto // " (standalone quote - repeats previous node)
	TokenIndent
	TokenNewline
	TokenComment
	TokenError
)

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenChapter:
		return "CHAPTER"
	case TokenContext:
		return "CONTEXT"
	case TokenNode:
		return "NODE"
	case TokenArrow:
		return "ARROW"
	case TokenAlias:
		return "ALIAS"
	case TokenString:
		return "STRING"
	case TokenDitto:
		return "DITTO"
	case TokenIndent:
		return "INDENT"
	case TokenNewline:
		return "NEWLINE"
	case TokenComment:
		return "COMMENT"
	case TokenError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Token represents a lexical token in N4L
type Token struct {
	Type    TokenType
	Value   string
	Line    int
	Column  int
	Literal string // Original text
}

// Position tracks location in source file
type Position struct {
	Line   int
	Column int
	File   string
}

func (p Position) String() string {
	if p.File != "" {
		return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Column)
	}
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Error represents a compilation error with proper context
type Error struct {
	Pos     Position
	Message string
	Code    string // The actual line of code
	Hint    string // Suggestion for fixing
}

func (e *Error) Error() string {
	var b strings.Builder

	// Location
	fmt.Fprintf(&b, "\n%s: error: %s\n", e.Pos, e.Message)

	// Show the actual code line if available
	if e.Code != "" {
		fmt.Fprintf(&b, "  %4d | %s\n", e.Pos.Line, e.Code)

		// Point to the error position
		if e.Pos.Column > 0 {
			fmt.Fprintf(&b, "       | %s^\n", strings.Repeat(" ", e.Pos.Column-1))
		}
	}

	// Hint for fixing
	if e.Hint != "" {
		fmt.Fprintf(&b, "  hint: %s\n", e.Hint)
	}

	return b.String()
}

// ErrorList is a collection of compilation errors
type ErrorList []*Error

func (e ErrorList) Error() string {
	var b strings.Builder
	fmt.Fprintf(&b, "\nFound %d error(s):\n", len(e))
	for _, err := range e {
		b.WriteString(err.Error())
	}
	return b.String()
}

// Arrow represents an arrow definition (built-in or user-defined)
type Arrow struct {
	Name        string
	ShortName   string
	Inverse     string
	Description string
	BuiltIn     bool // Core arrows that always exist
}

// CoreArrows are the built-in arrows that always work
var CoreArrows = map[string]*Arrow{
	"contain": {
		Name:        "contain",
		ShortName:   "contain",
		Inverse:     "belong",
		Description: "Basic containment relationship",
		BuiltIn:     true,
	},
	"belong": {
		Name:        "belong",
		ShortName:   "belong",
		Inverse:     "contain",
		Description: "Inverse of contain",
		BuiltIn:     true,
	},
	"hasX": {
		Name:        "has value",
		ShortName:   "hasX",
		Inverse:     "isXof",
		Description: "Has a property value",
		BuiltIn:     true,
	},
	"isXof": {
		Name:        "is the value of",
		ShortName:   "isXof",
		Inverse:     "hasX",
		Description: "Is value of property",
		BuiltIn:     true,
	},
	"next": {
		Name:        "next",
		ShortName:   "next",
		Inverse:     "prev",
		Description: "Sequential ordering",
		BuiltIn:     true,
	},
	"prev": {
		Name:        "prev",
		ShortName:   "prev",
		Inverse:     "next",
		Description: "Reverse sequential ordering",
		BuiltIn:     true,
	},
}

// Node represents a parsed N4L node
type Node struct {
	Name     string
	Alias    string
	Children []*Node
	Arrows   []*ArrowInstance
	Pos      Position
}

// ArrowInstance represents an arrow usage
type ArrowInstance struct {
	Arrow  *Arrow
	Target string
	Pos    Position
}

// AST represents the parsed N4L document
type AST struct {
	Chapter  string
	Contexts []string
	Nodes    []*Node
	Arrows   map[string]*Arrow // All arrows (core + user-defined)
	Errors   ErrorList
}

// Lexer tokenizes N4L input
type Lexer struct {
	reader  *bufio.Reader
	file    string
	line    int
	column  int
	current rune
	errors  ErrorList
}

// NewLexer creates a new lexer for the given input
func NewLexer(r io.Reader, filename string) *Lexer {
	lexer := &Lexer{
		reader: bufio.NewReader(r),
		file:   filename,
		line:   1,
		column: 0,
	}
	lexer.advance() // Read first character
	return lexer
}

// advance reads the next character
func (l *Lexer) advance() {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		l.current = 0 // EOF
		return
	}
	l.current = r
	l.column++
	if r == '\n' {
		l.line++
		l.column = 0
	}
}

// peek looks ahead without consuming
func (l *Lexer) peek() rune {
	bytes, err := l.reader.Peek(1)
	if err != nil || len(bytes) == 0 {
		return 0
	}
	return rune(bytes[0])
}

// position returns current position
func (l *Lexer) position() Position {
	return Position{
		Line:   l.line,
		Column: l.column,
		File:   l.file,
	}
}

// NextToken returns the next token
func (l *Lexer) NextToken() *Token {
	// Skip whitespace (except newlines and indentation)
	for l.current == ' ' || l.current == '\t' {
		l.advance()
	}

	pos := l.position()

	// EOF
	if l.current == 0 {
		return &Token{Type: TokenEOF, Line: pos.Line, Column: pos.Column}
	}

	// Newline
	if l.current == '\n' {
		l.advance()
		return &Token{Type: TokenNewline, Line: pos.Line, Column: pos.Column}
	}

	// Chapter marker
	if l.current == '-' && l.column == 1 {
		l.advance()
		if l.current == ' ' {
			l.advance()
			return l.readChapter(pos)
		}
	}

	// Context marker
	if l.current == '+' && l.column == 1 {
		l.advance()
		if l.current == ' ' {
			l.advance()
			return l.readContext(pos)
		}
	}

	// Alias marker
	if l.current == '@' {
		return l.readAlias(pos)
	}

	// Arrow (quoted string with parentheses)
	if l.current == '"' {
		// Check if it's a ditto (standalone quote at start of line/word)
		peek := l.peek()
		if peek == ' ' || peek == '\t' || peek == '\n' || peek == 0 {
			// This is a ditto mark - it means "repeat previous node"
			l.advance() // consume the quote
			return &Token{
				Type:   TokenDitto,
				Value:  "PREV",
				Line:   pos.Line,
				Column: pos.Column,
			}
		}
		return l.readQuotedString(pos)
	}

	// Node (unquoted identifier)
	return l.readIdentifier(pos)
}

func (l *Lexer) readChapter(pos Position) *Token {
	var value strings.Builder
	for l.current != '\n' && l.current != 0 {
		value.WriteRune(l.current)
		l.advance()
	}
	return &Token{
		Type:   TokenChapter,
		Value:  strings.TrimSpace(value.String()),
		Line:   pos.Line,
		Column: pos.Column,
	}
}

func (l *Lexer) readContext(pos Position) *Token {
	var value strings.Builder
	for l.current != '\n' && l.current != 0 {
		value.WriteRune(l.current)
		l.advance()
	}
	return &Token{
		Type:   TokenContext,
		Value:  strings.TrimSpace(value.String()),
		Line:   pos.Line,
		Column: pos.Column,
	}
}

func (l *Lexer) readAlias(pos Position) *Token {
	var value strings.Builder
	l.advance() // skip @

	for l.current != ' ' && l.current != '\n' && l.current != 0 {
		value.WriteRune(l.current)
		l.advance()
	}

	return &Token{
		Type:   TokenAlias,
		Value:  value.String(),
		Line:   pos.Line,
		Column: pos.Column,
	}
}

func (l *Lexer) readQuotedString(pos Position) *Token {
	var value strings.Builder
	l.advance() // skip opening quote

	for l.current != '"' && l.current != '\n' && l.current != 0 {
		if l.current == '\\' {
			l.advance()
			if l.current != 0 {
				value.WriteRune(l.current)
				l.advance()
			}
		} else {
			value.WriteRune(l.current)
			l.advance()
		}
	}

	if l.current != '"' {
		l.errors = append(l.errors, &Error{
			Pos:     pos,
			Message: "unterminated string",
			Hint:    "add closing quote",
		})
		return &Token{Type: TokenError, Line: pos.Line, Column: pos.Column}
	}

	l.advance() // skip closing quote

	return &Token{
		Type:   TokenString,
		Value:  value.String(),
		Line:   pos.Line,
		Column: pos.Column,
	}
}

func (l *Lexer) readIdentifier(pos Position) *Token {
	var value strings.Builder

	for l.current != ' ' && l.current != '\n' && l.current != 0 && l.current != '"' {
		value.WriteRune(l.current)
		l.advance()
	}

	return &Token{
		Type:   TokenNode,
		Value:  value.String(),
		Line:   pos.Line,
		Column: pos.Column,
	}
}

// Compiler compiles N4L source files
type Compiler struct {
	filename string
	errors   ErrorList
}

// NewCompiler creates a new N4L compiler
func NewCompiler(filename string) *Compiler {
	return &Compiler{
		filename: filename,
	}
}

// readSourceLines reads all source lines for error display
func (c *Compiler) readSourceLines() []string {
	file, err := os.Open(c.filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// Compile compiles an N4L file and returns the AST or errors
func (c *Compiler) Compile() (*AST, error) {
	file, err := os.Open(c.filename)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	// Read source for error display
	sourceLines := c.readSourceLines()

	// Lexical analysis
	file.Seek(0, 0) // Reset to beginning
	lexer := NewLexer(file, c.filename)

	// TODO: Parser implementation
	// For now, just tokenize and check for errors
	for {
		token := lexer.NextToken()
		if token.Type == TokenEOF {
			break
		}
		if token.Type == TokenError {
			break
		}
	}

	// Add source code to errors
	if len(lexer.errors) > 0 {
		for _, err := range lexer.errors {
			if err.Pos.Line > 0 && err.Pos.Line <= len(sourceLines) {
				err.Code = sourceLines[err.Pos.Line-1]
			}
		}
		return nil, lexer.errors
	}

	// Create AST with core arrows
	ast := &AST{
		Arrows: make(map[string]*Arrow),
	}

	// Load core arrows
	for name, arrow := range CoreArrows {
		ast.Arrows[name] = arrow
	}

	return ast, nil
}

// ValidateArrow checks if an arrow is valid (core or user-defined)
func (ast *AST) ValidateArrow(name string, pos Position) error {
	if _, ok := ast.Arrows[name]; !ok {
		return &Error{
			Pos:     pos,
			Message: fmt.Sprintf("undefined arrow: (%s)", name),
			Hint:    "define the arrow in SSTconfig or use a core arrow: " + strings.Join(coreArrowNames(), ", "),
		}
	}
	return nil
}

func coreArrowNames() []string {
	names := make([]string, 0, len(CoreArrows))
	for name := range CoreArrows {
		names = append(names, name)
	}
	return names
}
