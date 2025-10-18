package json2n4l

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestParseString tests the parser with various JSON inputs
func TestParseString(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		rootName string
		wantErr  bool
		contains []string
	}{
		{
			name:     "simple object",
			json:     `{"name": "test", "value": 42}`,
			rootName: "testobj",
			wantErr:  false,
			contains: []string{"testobj", `" (contain) name`, `" (contain) test`, `" (contain) value`, `" (contain) 42`},
		},
		{
			name:     "nested object",
			json:     `{"user": {"name": "Alice"}}`,
			rootName: "root",
			wantErr:  false,
			contains: []string{"root", `" (contain) user`, "user", `" (contain) name`, `" (contain) Alice`},
		},
		{
			name:     "array",
			json:     `{"items": [1, 2, 3]}`,
			rootName: "data",
			wantErr:  false,
			contains: []string{"items", `" (contain) "items[0]"`, `" (contain) "items[1]"`, `" (contain) "items[2]"`},
		},
		{
			name:     "boolean and null",
			json:     `{"active": true, "data": null}`,
			rootName: "test",
			wantErr:  false,
			contains: []string{"active", `" (contain) true`, "data", `" (contain) null`},
		},
		{
			name:     "invalid json",
			json:     `{invalid}`,
			rootName: "test",
			wantErr:  true,
		},
		{
			name:     "empty object",
			json:     `{}`,
			rootName: "empty",
			wantErr:  false,
			contains: []string{"empty"},
		},
		{
			name:     "empty array",
			json:     `{"arr": []}`,
			rootName: "root",
			wantErr:  false,
			contains: []string{"root", `" (contain) arr`, "arr"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				RootName:    tt.rootName,
				ChapterName: "Test Chapter",
			}
			parser := NewParser(config)

			result, err := parser.ParseString(tt.json)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for _, want := range tt.contains {
					if !strings.Contains(result, want) {
						t.Errorf("ParseString() result missing expected content: %q", want)
						t.Logf("Full output:\n%s", result)
					}
				}
			}
		})
	}
}

// TestEscapeN4L tests the N4L escaping rules
func TestEscapeN4L(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"simple", "simple"},
		{"with space", `"with space"`},
		{"with\\backslash", `with\\backslash`},
		{"with{braces}", `with\{braces\}`},
		{"exclaim!", `"exclaim!"`},
		{"array[0]", `"array[0]"`},
		{"normal_key", "normal_key"},
	}

	parser := &Parser{}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := parser.escapeN4L(tt.input)
			if got != tt.want {
				t.Errorf("escapeN4L(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestWithOptions(t *testing.T) {
	config := Config{}

	// Test option functions
	WithChapter("Test Chapter")(&config)
	if config.ChapterName != "Test Chapter" {
		t.Errorf("WithChapter() failed, got %q", config.ChapterName)
	}

	WithContext("tag1", "tag2")(&config)
	if len(config.ContextTags) != 2 {
		t.Errorf("WithContext() failed, got %d tags", len(config.ContextTags))
	}

	WithTypes()(&config)
	if !config.IncludeTypes {
		t.Error("WithTypes() failed")
	}

	WithPrettyPrint()(&config)
	if !config.PrettyPrint {
		t.Error("WithPrettyPrint() failed")
	}

	WithMaxDepth(5)(&config)
	if config.MaxDepth != 5 {
		t.Errorf("WithMaxDepth() failed, got %d", config.MaxDepth)
	}
}

// TestGetStats tests the statistics gathering
func TestGetStats(t *testing.T) {
	config := Config{
		RootName:    "test",
		ChapterName: "Test",
	}
	parser := NewParser(config)

	json := `{"a": 1, "b": {"c": 2}}`
	_, err := parser.ParseString(json)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	stats := parser.GetStats()

	if nodeCount, ok := stats["nodeCount"].(int); !ok || nodeCount == 0 {
		t.Error("GetStats() should return non-zero nodeCount")
	}

	if outputSize, ok := stats["outputSize"].(int); !ok || outputSize == 0 {
		t.Error("GetStats() should return non-zero outputSize")
	}
}

// TestParseFile tests file-based parsing with buffered I/O
func TestParseFile(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "test.json")
	outputFile := filepath.Join(tmpDir, "test.n4l")

	testJSON := `{
		"name": "Test",
		"count": 42,
		"items": ["a", "b", "c"]
	}`

	if err := os.WriteFile(inputFile, []byte(testJSON), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := Config{
		InputFile:   inputFile,
		OutputFile:  outputFile,
		RootName:    "testdata",
		ChapterName: "File Test",
		PrettyPrint: true,
	}

	parser := NewParser(config)
	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Verify output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Parse() did not create output file")
	}

	// Read and verify content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	output := string(content)
	expectedStrings := []string{
		"- File Test",
		"testdata",
		`" (contain) name`,
		`" (contain) Test`,
		`" (contain) count`,
		`" (contain) 42`,
		`" (contain) items`,
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("Output missing expected string: %q", expected)
		}
	}

	stats := parser.GetStats()
	if stats["nodeCount"].(int) == 0 {
		t.Error("Parse() should generate nodes")
	}
}

// TestPrettyPrint verifies pretty print comments
func TestPrettyPrint(t *testing.T) {
	config := Config{
		RootName:    "test",
		ChapterName: "Test",
		PrettyPrint: true,
	}
	parser := NewParser(config)

	json := `{"obj": {"key": "value"}, "arr": [1, 2]}`
	result, err := parser.ParseString(json)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	// Check for pretty print comments
	if !strings.Contains(result, "# Converted from JSON to N4L") {
		t.Error("Pretty print should include conversion comment")
	}
	if !strings.Contains(result, "# JSON Object:") {
		t.Error("Pretty print should include object comments")
	}
	if !strings.Contains(result, "# JSON Array:") {
		t.Error("Pretty print should include array comments")
	}
}

// TestMaxDepth verifies depth limiting
func TestMaxDepth(t *testing.T) {
	config := Config{
		RootName:    "test",
		ChapterName: "Test",
		MaxDepth:    2,
		PrettyPrint: true,
	}
	parser := NewParser(config)

	// Deeply nested JSON
	json := `{"a": {"b": {"c": {"d": "too deep"}}}}`
	result, err := parser.ParseString(json)
	if err != nil {
		t.Fatalf("ParseString() error = %v", err)
	}

	// Should stop at depth 2
	if strings.Contains(result, "too deep") {
		t.Error("MaxDepth should prevent parsing deeply nested values")
	}
	if !strings.Contains(result, "# Max depth reached") {
		t.Error("MaxDepth should add a comment when limit is reached")
	}
}

//
// BENCHMARKS
//

// BenchmarkParseString_Small benchmarks parsing small JSON (< 1KB)
func BenchmarkParseString_Small(b *testing.B) {
	config := Config{
		RootName:    "bench",
		ChapterName: "Benchmark",
	}

	json := `{"name": "test", "value": 42, "active": true}`

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parser := NewParser(config)
		_, err := parser.ParseString(json)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseString_Medium benchmarks parsing medium JSON (~10KB)
func BenchmarkParseString_Medium(b *testing.B) {
	config := Config{
		RootName:    "bench",
		ChapterName: "Benchmark",
	}

	// Generate medium-sized JSON
	json := `{
		"users": [
			{"id": 1, "name": "Alice", "email": "alice@example.com", "active": true},
			{"id": 2, "name": "Bob", "email": "bob@example.com", "active": false},
			{"id": 3, "name": "Charlie", "email": "charlie@example.com", "active": true},
			{"id": 4, "name": "Diana", "email": "diana@example.com", "active": true},
			{"id": 5, "name": "Eve", "email": "eve@example.com", "active": false}
		],
		"metadata": {
			"version": "1.0",
			"timestamp": "2024-01-01T00:00:00Z",
			"count": 5,
			"tags": ["test", "benchmark", "performance"]
		}
	}`

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parser := NewParser(config)
		_, err := parser.ParseString(json)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseString_Large benchmarks parsing large JSON (~100KB)
func BenchmarkParseString_Large(b *testing.B) {
	config := Config{
		RootName:    "bench",
		ChapterName: "Benchmark",
	}

	// Generate large JSON with many repeated structures
	var builder strings.Builder
	builder.WriteString(`{"items": [`)
	for i := 0; i < 100; i++ {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(`{
			"id": ` + strings.Repeat("1", 10) + `,
			"name": "item_` + strings.Repeat("x", 50) + `",
			"description": "` + strings.Repeat("Lorem ipsum dolor sit amet ", 10) + `",
			"tags": ["tag1", "tag2", "tag3", "tag4", "tag5"],
			"metadata": {"key1": "value1", "key2": "value2", "key3": "value3"}
		}`)
	}
	builder.WriteString(`]}`)
	json := builder.String()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parser := NewParser(config)
		_, err := parser.ParseString(json)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseString_DeepNesting benchmarks deeply nested structures
func BenchmarkParseString_DeepNesting(b *testing.B) {
	config := Config{
		RootName:    "bench",
		ChapterName: "Benchmark",
	}

	// Generate deeply nested JSON (10 levels deep)
	json := `{"a": {"b": {"c": {"d": {"e": {"f": {"g": {"h": {"i": {"j": "deep"}}}}}}}}}}`

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parser := NewParser(config)
		_, err := parser.ParseString(json)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseString_WideArray benchmarks arrays with many elements
func BenchmarkParseString_WideArray(b *testing.B) {
	config := Config{
		RootName:    "bench",
		ChapterName: "Benchmark",
	}

	// Generate array with 1000 elements
	var builder strings.Builder
	builder.WriteString(`{"numbers": [`)
	for i := 0; i < 1000; i++ {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(`{"value": `)
		builder.WriteString(strings.Repeat("1", 5))
		builder.WriteString(`}`)
	}
	builder.WriteString(`]}`)
	json := builder.String()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		parser := NewParser(config)
		_, err := parser.ParseString(json)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkParseFile benchmarks file-based parsing with buffered I/O
func BenchmarkParseFile(b *testing.B) {
	tmpDir := b.TempDir()
	inputFile := filepath.Join(tmpDir, "bench.json")

	testJSON := `{
		"data": {
			"users": [
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"}
			],
			"metadata": {"version": "1.0"}
		}
	}`

	if err := os.WriteFile(inputFile, []byte(testJSON), 0o644); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		outputFile := filepath.Join(tmpDir, "bench_"+strings.Repeat("x", 10)+".n4l")
		config := Config{
			InputFile:  inputFile,
			OutputFile: outputFile,
			RootName:   "bench",
		}

		parser := NewParser(config)
		if err := parser.Parse(); err != nil {
			b.Fatal(err)
		}
	}
}
