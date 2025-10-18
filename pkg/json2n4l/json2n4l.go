// Package json2n4l provides a Computer Science compliant parser for JSON to N4L conversion.
//
// This parser follows proper compiler design principles:
// - Lexical analysis via json.Decoder (tokenization)
// - Syntactic analysis via recursive descent parser
// - Semantic analysis via N4L graph construction
// - Code generation via buffered output writer
//
// N4L Semantic Mapping (using stable CN-2 Contains arrows):
// - Object keys → (contain) relationship with ditto notation
// - Array elements → (contain) indexed elements
// - Values → (contain) leaf node values
// - Proper N4L syntax: statement followed by " (arrow) target
package json2n4l

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ArrowStyle defines which semantic arrows to use for relationships
type ArrowStyle int

const (
	// ArrowStyleSimple uses only (contain) arrows for all relationships
	ArrowStyleSimple ArrowStyle = iota
	// ArrowStyleSemantic uses appropriate semantic arrows based on context
	ArrowStyleSemantic
	// ArrowStyleBidirectional adds inverse arrows (belong, pt-of, etc.)
	ArrowStyleBidirectional
)

// Config holds the conversion configuration
type Config struct {
	InputFile          string     // Path to input JSON file
	OutputFile         string     // Path to output N4L file (optional, defaults to input.n4l)
	RootName           string     // Name for the root node (defaults to filename)
	ChapterName        string     // Chapter name for the N4L document
	ContextTags        []string   // Context tags to add
	AutoContext        bool       // Automatically add context tags based on JSON structure
	IncludeTypes       bool       // Include type information as annotations
	MaxDepth           int        // Maximum nesting depth (0 = unlimited)
	PrettyPrint        bool       // Add comments for readability
	GenerateAliases    bool       // Generate @alias references for objects
	ArrowStyle         ArrowStyle // Which semantic arrows to use
	PreserveStructure  bool       // Use (pt-of) to preserve hierarchical structure
	UseSequenceForList bool       // Use _sequence_ mode for arrays (when order matters)
}

// Parser handles JSON to N4L conversion using proper CS parser design
type Parser struct {
	config       Config
	writer       *bufio.Writer   // Buffered output writer
	buf          *bytes.Buffer   // Internal buffer for string operations
	depth        int             // Current nesting depth
	nodeCount    int             // Total nodes generated
	outputSize   int             // Final output size in bytes
	aliasCounter int             // Counter for generating unique aliases
	contextStack []string        // Stack of active contexts
	parentStack  []string        // Stack of parent nodes for hierarchy
	detectedTags map[string]bool // Auto-detected context tags
}

// NewParser creates a new JSON to N4L parser with buffered I/O
func NewParser(config Config) *Parser {
	// Set defaults
	if config.OutputFile == "" {
		ext := filepath.Ext(config.InputFile)
		config.OutputFile = strings.TrimSuffix(config.InputFile, ext) + ".n4l"
	}
	if config.RootName == "" {
		base := filepath.Base(config.InputFile)
		config.RootName = strings.TrimSuffix(base, filepath.Ext(base))
	}
	if config.ChapterName == "" {
		config.ChapterName = "JSON Import"
	}

	buf := &bytes.Buffer{}
	return &Parser{
		config:       config,
		writer:       bufio.NewWriter(buf),
		buf:          buf,
		detectedTags: make(map[string]bool),
	}
}

// Parse reads a JSON file and converts it to N4L format using buffered I/O
func (p *Parser) Parse() error {
	// Open input file with buffered reader
	inFile, err := os.Open(p.config.InputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)

	// Decode JSON using streaming decoder
	decoder := json.NewDecoder(reader)

	var jsonData interface{}
	if err := decoder.Decode(&jsonData); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Generate N4L document
	if err := p.writeHeader(); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	if err := p.parseValue(p.config.RootName, jsonData, 0); err != nil {
		return fmt.Errorf("failed to parse JSON structure: %w", err)
	}

	// Flush writer to ensure all data is written to buffer
	if err := p.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush buffer: %w", err)
	}

	// Capture output size for stats
	p.outputSize = p.buf.Len()

	// Write output file with buffered writer
	outFile, err := os.Create(p.config.OutputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	if _, err := p.buf.WriteTo(writer); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}

	return nil
}

// ParseString parses JSON string directly to N4L string using buffered I/O
func (p *Parser) ParseString(jsonStr string) (string, error) {
	// Use streaming decoder for consistent parsing
	decoder := json.NewDecoder(strings.NewReader(jsonStr))

	var jsonData interface{}
	if err := decoder.Decode(&jsonData); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	if err := p.writeHeader(); err != nil {
		return "", fmt.Errorf("failed to write header: %w", err)
	}

	if err := p.parseValue(p.config.RootName, jsonData, 0); err != nil {
		return "", fmt.Errorf("failed to parse JSON structure: %w", err)
	}

	if err := p.writer.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush buffer: %w", err)
	}

	// Capture output size for stats
	p.outputSize = p.buf.Len()

	return p.buf.String(), nil
}

// writeHeader writes the N4L document header using buffered writer
func (p *Parser) writeHeader() error {
	// Simple section header (N4L style)
	if _, err := p.writer.WriteString(fmt.Sprintf("- %s\n\n", p.config.ChapterName)); err != nil {
		return err
	}

	// Build context tags
	contextTags := make([]string, 0, len(p.config.ContextTags)+len(p.detectedTags))
	contextTags = append(contextTags, p.config.ContextTags...)

	if p.config.AutoContext {
		// Add auto-detected tags
		for tag := range p.detectedTags {
			contextTags = append(contextTags, tag)
		}
	}

	if len(contextTags) > 0 {
		if _, err := p.writer.WriteString(fmt.Sprintf(":: %s ::\n\n", strings.Join(contextTags, ", "))); err != nil {
			return err
		}
	}

	if p.config.PrettyPrint {
		if _, err := p.writer.WriteString("# Converted from JSON to N4L\n"); err != nil {
			return err
		}
		if p.config.ArrowStyle == ArrowStyleSemantic {
			if _, err := p.writer.WriteString("# Using semantic arrows from SSTconfig CN-2\n"); err != nil {
				return err
			}
		}
		if _, err := p.writer.WriteString("\n"); err != nil {
			return err
		}
	}

	return nil
}

// getContainmentArrow returns the appropriate arrow for containment based on style
func (p *Parser) getContainmentArrow() string {
	switch p.config.ArrowStyle {
	case ArrowStyleSemantic, ArrowStyleBidirectional:
		return "contain"
	default:
		return "contain"
	}
}

// getBelongArrow returns the inverse arrow for containment
func (p *Parser) getBelongArrow() string {
	if p.config.ArrowStyle == ArrowStyleBidirectional {
		return "belong"
	}
	return ""
}

// getPartOfArrow returns the arrow for part-of relationships
func (p *Parser) getPartOfArrow() string {
	if p.config.PreserveStructure {
		return "pt-of"
	}
	return ""
}

// getHasPartArrow returns the arrow for has-part relationships
func (p *Parser) getHasPartArrow() string {
	if p.config.PreserveStructure {
		return "has-pt"
	}
	return ""
}

// getSetMemberArrow returns arrows for set membership
func (p *Parser) getSetMemberArrow() (string, string) {
	if p.config.ArrowStyle >= ArrowStyleSemantic {
		return "in-set", "setof"
	}
	return "belong", "contain"
}

// generateAlias creates a unique alias for a node
func (p *Parser) generateAlias(baseName string) string {
	if !p.config.GenerateAliases {
		return ""
	}
	p.aliasCounter++
	// Clean base name for valid alias
	cleaned := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		}
		return '_'
	}, baseName)
	return fmt.Sprintf("@%s_%d", cleaned, p.aliasCounter)
}

// detectContext adds semantic context tags based on JSON keys
func (p *Parser) detectContext(keys []string) {
	if !p.config.AutoContext {
		return
	}

	// Common semantic patterns
	semanticPatterns := map[string][]string{
		"data":     {"data", "dataset"},
		"user":     {"user", "profile", "account"},
		"config":   {"configuration", "settings"},
		"api":      {"api", "endpoint", "response"},
		"metadata": {"metadata", "info"},
		"id":       {"identifier"},
		"name":     {"naming"},
		"email":    {"contact"},
		"address":  {"location", "geo"},
		"phone":    {"contact"},
		"date":     {"temporal", "time"},
		"time":     {"temporal"},
		"status":   {"state"},
		"type":     {"classification"},
		"error":    {"error", "exception"},
		"message":  {"messaging"},
	}

	for _, key := range keys {
		lowerKey := strings.ToLower(key)
		for pattern, tags := range semanticPatterns {
			if strings.Contains(lowerKey, pattern) {
				for _, tag := range tags {
					p.detectedTags[tag] = true
				}
			}
		}
	}
}

// parseValue recursively parses JSON values using proper recursive descent
func (p *Parser) parseValue(name string, value interface{}, depth int) error {
	// Check max depth
	if p.config.MaxDepth > 0 && depth > p.config.MaxDepth {
		if p.config.PrettyPrint {
			_, err := p.writer.WriteString(fmt.Sprintf("# Max depth reached at: %s\n", name))
			return err
		}
		return nil
	}

	p.nodeCount++

	// Type dispatch for recursive descent parsing
	switch v := value.(type) {
	case map[string]interface{}:
		return p.parseObject(name, v, depth)
	case []interface{}:
		return p.parseArray(name, v, depth)
	case string:
		return p.parseString(name, v, depth)
	case float64:
		return p.parseNumber(name, v, depth)
	case bool:
		return p.parseBool(name, v, depth)
	case nil:
		return p.parseNull(name, depth)
	default:
		return p.parseUnknown(name, v, depth)
	}
}

// parseObject parses a JSON object using recursive descent
func (p *Parser) parseObject(name string, obj map[string]interface{}, depth int) error {
	indent := strings.Repeat(" ", depth)

	// Auto-detect context from object keys
	if p.config.AutoContext && depth == 1 {
		keys := make([]string, 0, len(obj))
		for k := range obj {
			keys = append(keys, k)
		}
		p.detectContext(keys)
	}

	// Create the object node
	if p.config.PrettyPrint {
		if _, err := p.writer.WriteString(fmt.Sprintf("%s# JSON Object: %s (%d keys)\n", indent, name, len(obj))); err != nil {
			return err
		}
	}

	// Generate alias if enabled
	alias := p.generateAlias(name)
	if alias != "" {
		if _, err := p.writer.WriteString(fmt.Sprintf("%s %s\n", alias, p.escapeN4L(name))); err != nil {
			return err
		}
	} else {
		if _, err := p.writer.WriteString(fmt.Sprintf(" %s\n", p.escapeN4L(name))); err != nil {
			return err
		}
	}

	// Add parent to stack for hierarchy tracking
	p.parentStack = append(p.parentStack, name)

	// Process each key-value pair with recursive descent
	arrow := p.getContainmentArrow()
	for key, val := range obj {
		// Object contains key (stable CN-2 arrow) using ditto format
		if _, err := p.writer.WriteString(fmt.Sprintf("      \" (%s) %s\n", arrow, p.escapeN4L(key))); err != nil {
			return err
		}

		// Recurse into the value
		if err := p.parseValue(key, val, depth+1); err != nil {
			return err
		}

		// Add inverse relationship if bidirectional style (key belongs to object)
		if p.config.ArrowStyle == ArrowStyleBidirectional {
			inverseArrow := p.getBelongArrow()
			if inverseArrow != "" {
				// Use explicit subject to avoid self-reference
				if _, err := p.writer.WriteString(fmt.Sprintf(" %s\n", p.escapeN4L(key))); err != nil {
					return err
				}
				if _, err := p.writer.WriteString(fmt.Sprintf("      \" (%s) %s\n", inverseArrow, p.escapeN4L(name))); err != nil {
					return err
				}
			}
		}
	}

	// Pop parent from stack
	if len(p.parentStack) > 0 {
		p.parentStack = p.parentStack[:len(p.parentStack)-1]
	}

	if _, err := p.writer.WriteString("\n"); err != nil {
		return err
	}

	return nil
}

// parseArray parses a JSON array using recursive descent
func (p *Parser) parseArray(name string, arr []interface{}, depth int) error {
	indent := strings.Repeat(" ", depth)

	if p.config.PrettyPrint {
		if _, err := p.writer.WriteString(fmt.Sprintf("%s# JSON Array: %s [%d items]\n", indent, name, len(arr))); err != nil {
			return err
		}
	}

	// Check if we should use sequence mode
	useSequence := p.config.UseSequenceForList && len(arr) > 0

	if useSequence {
		if _, err := p.writer.WriteString(":: _sequence_ ::\n\n"); err != nil {
			return err
		}
	}

	if _, err := p.writer.WriteString(fmt.Sprintf(" %s\n", p.escapeN4L(name))); err != nil {
		return err
	}

	// Add semantic description for arrays
	memberArrow, collectionArrow := p.getSetMemberArrow()
	if p.config.ArrowStyle >= ArrowStyleSemantic && len(arr) > 0 {
		if _, err := p.writer.WriteString(fmt.Sprintf("      \" (%s) collection\n", collectionArrow)); err != nil {
			return err
		}
	}

	// Process each array element with recursive descent
	arrow := p.getContainmentArrow()
	for i, val := range arr {
		elementName := fmt.Sprintf("%s[%d]", name, i)

		// Array contains element (stable CN-2 arrow) using ditto format
		if _, err := p.writer.WriteString(fmt.Sprintf("      \" (%s) %s\n", arrow, p.escapeN4L(elementName))); err != nil {
			return err
		}

		// Recurse into the element
		if err := p.parseValue(elementName, val, depth+1); err != nil {
			return err
		}

		// Add set membership relationship if semantic style (element belongs to collection)
		if p.config.ArrowStyle >= ArrowStyleSemantic {
			// Use explicit subject to avoid self-reference
			if _, err := p.writer.WriteString(fmt.Sprintf(" %s\n", p.escapeN4L(elementName))); err != nil {
				return err
			}
			if _, err := p.writer.WriteString(fmt.Sprintf("      \" (%s) %s\n", memberArrow, p.escapeN4L(name))); err != nil {
				return err
			}
		}
	}

	if useSequence {
		if _, err := p.writer.WriteString("\n-:: _sequence_ ::\n"); err != nil {
			return err
		}
	}

	if _, err := p.writer.WriteString("\n"); err != nil {
		return err
	}

	return nil
}

// parseString parses a JSON string value (terminal symbol)
func (p *Parser) parseString(name string, value string, depth int) error {
	// Use semantic arrow for properties if in semantic mode
	arrow := "contain"
	if p.config.ArrowStyle >= ArrowStyleSemantic {
		arrow = "propt" // Property of parent
	}

	// Add type annotation if requested
	typeInfo := ""
	if p.config.IncludeTypes {
		typeInfo = " {string}"
	}

	_, err := p.writer.WriteString(fmt.Sprintf("           \" (%s) %s%s\n", arrow, p.escapeN4L(value), typeInfo))
	return err
}

// parseNumber parses a JSON number value (terminal symbol)
func (p *Parser) parseNumber(name string, value float64, depth int) error {
	// Use semantic arrow for properties if in semantic mode
	arrow := "contain"
	if p.config.ArrowStyle >= ArrowStyleSemantic {
		arrow = "propt" // Property of parent
	}

	// Add type annotation if requested
	typeInfo := ""
	if p.config.IncludeTypes {
		typeInfo = " {number}"
	}

	_, err := p.writer.WriteString(fmt.Sprintf("           \" (%s) %.6g%s\n", arrow, value, typeInfo))
	return err
}

// parseBool parses a JSON boolean value (terminal symbol)
func (p *Parser) parseBool(name string, value bool, depth int) error {
	// Use semantic arrow for properties if in semantic mode
	arrow := "contain"
	if p.config.ArrowStyle >= ArrowStyleSemantic {
		arrow = "propt" // Property of parent
	}

	// Add type annotation if requested
	typeInfo := ""
	if p.config.IncludeTypes {
		typeInfo = " {bool}"
	}

	_, err := p.writer.WriteString(fmt.Sprintf("           \" (%s) %t%s\n", arrow, value, typeInfo))
	return err
}

// parseNull parses a JSON null value (terminal symbol)
func (p *Parser) parseNull(name string, depth int) error {
	// Use semantic arrow for properties if in semantic mode
	arrow := "contain"
	if p.config.ArrowStyle >= ArrowStyleSemantic {
		arrow = "propt" // Property of parent
	}

	// Add type annotation if requested
	typeInfo := ""
	if p.config.IncludeTypes {
		typeInfo = " {null}"
	}

	_, err := p.writer.WriteString(fmt.Sprintf("           \" (%s) null%s\n", arrow, typeInfo))
	return err
}

// parseUnknown handles unknown JSON types (error recovery)
func (p *Parser) parseUnknown(name string, value interface{}, depth int) error {
	// Use semantic arrow for properties if in semantic mode
	arrow := "contain"
	if p.config.ArrowStyle >= ArrowStyleSemantic {
		arrow = "propt" // Property of parent
	}

	_, err := p.writer.WriteString(fmt.Sprintf("           \" (%s) %v\n", arrow, value))
	return err
}

// escapeN4L escapes special characters for N4L format (lexical rules)
func (p *Parser) escapeN4L(s string) string {
	// Escape backslashes first
	s = strings.ReplaceAll(s, "\\", "\\\\")
	// Escape curly braces
	s = strings.ReplaceAll(s, "{", "\\{")
	s = strings.ReplaceAll(s, "}", "\\}")
	// Wrap in quotes if contains spaces or special chars
	if strings.ContainsAny(s, " \t\n![]") {
		s = fmt.Sprintf("\"%s\"", s)
	}
	return s
}

// GetStats returns parsing statistics
func (p *Parser) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"nodeCount":  p.nodeCount,
		"maxDepth":   p.depth,
		"outputSize": p.outputSize,
	}
}

// GetConfig returns the parser's configuration
func (p *Parser) GetConfig() Config {
	return p.config
}

// ParseReader parses JSON from an io.Reader to N4L format using streaming
func ParseReader(r io.Reader, config Config) (string, error) {
	parser := NewParser(config)

	// Use streaming decoder for efficient memory usage
	decoder := json.NewDecoder(r)

	var jsonData interface{}
	if err := decoder.Decode(&jsonData); err != nil {
		return "", fmt.Errorf("failed to decode JSON: %w", err)
	}

	if err := parser.writeHeader(); err != nil {
		return "", fmt.Errorf("failed to write header: %w", err)
	}

	if err := parser.parseValue(config.RootName, jsonData, 0); err != nil {
		return "", fmt.Errorf("failed to parse value: %w", err)
	}

	if err := parser.writer.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush: %w", err)
	}

	// Capture output size for stats
	parser.outputSize = parser.buf.Len()

	return parser.buf.String(), nil
}

// ParseFile is a convenience function to parse a JSON file to N4L
func ParseFile(inputPath, outputPath string, options ...func(*Config)) error {
	config := Config{
		InputFile:  inputPath,
		OutputFile: outputPath,
	}

	// Apply options
	for _, opt := range options {
		opt(&config)
	}

	parser := NewParser(config)
	return parser.Parse()
}

// Option functions for convenience

// WithChapter sets the chapter name
func WithChapter(name string) func(*Config) {
	return func(c *Config) {
		c.ChapterName = name
	}
}

// WithContext sets context tags
func WithContext(tags ...string) func(*Config) {
	return func(c *Config) {
		c.ContextTags = tags
	}
}

// WithTypes enables type annotations
func WithTypes() func(*Config) {
	return func(c *Config) {
		c.IncludeTypes = true
	}
}

// WithPrettyPrint enables pretty printing with comments
func WithPrettyPrint() func(*Config) {
	return func(c *Config) {
		c.PrettyPrint = true
	}
}

// WithMaxDepth sets maximum nesting depth
func WithMaxDepth(depth int) func(*Config) {
	return func(c *Config) {
		c.MaxDepth = depth
	}
}

// WithArrowStyle sets the semantic arrow style to use
func WithArrowStyle(style ArrowStyle) func(*Config) {
	return func(c *Config) {
		c.ArrowStyle = style
	}
}

// WithSemantic enables semantic arrow selection (shorthand for ArrowStyleSemantic)
func WithSemantic() func(*Config) {
	return func(c *Config) {
		c.ArrowStyle = ArrowStyleSemantic
	}
}

// WithBidirectional enables bidirectional arrows with inverse relationships
func WithBidirectional() func(*Config) {
	return func(c *Config) {
		c.ArrowStyle = ArrowStyleBidirectional
	}
}

// WithAutoContext enables automatic context tag detection from JSON structure
func WithAutoContext() func(*Config) {
	return func(c *Config) {
		c.AutoContext = true
	}
}

// WithAliases enables generation of @alias references for complex objects
func WithAliases() func(*Config) {
	return func(c *Config) {
		c.GenerateAliases = true
	}
}

// WithPreserveStructure enables (pt-of) arrows to preserve hierarchical structure
func WithPreserveStructure() func(*Config) {
	return func(c *Config) {
		c.PreserveStructure = true
	}
}

// WithSequenceForList enables _sequence_ mode for arrays when order matters
func WithSequenceForList() func(*Config) {
	return func(c *Config) {
		c.UseSequenceForList = true
	}
}
