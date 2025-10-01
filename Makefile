
# Modern Makefile for SSTorytime project
# Builds all executables using the new cmd/ directory structure

# Define all executables
TOOLS = text2N4L pathsolve notes graph_report
N4L_TOOLS = N4L
SEARCH_TOOLS = searchN4L removeN4L
EXAMPLES = API_EXAMPLE_1 API_EXAMPLE_2 API_EXAMPLE_3 API_EXAMPLE_4
SERVER = http_server
DEMO_POCS = dotest_entirecone dotest_getnodes search_coarse_grain search_coarse_grain2 search_coarse_grain_api

ALL_TARGETS = $(TOOLS) $(N4L_TOOLS) $(SEARCH_TOOLS) $(EXAMPLES) $(SERVER) $(DEMO_POCS)

# Default target
all: $(ALL_TARGETS)

# Tools (in cmd/tools/)
$(TOOLS):
	@echo "Building $@..."
	cd cmd/tools && go build -o ../../$@ $@.go

# N4L tools (in cmd/n4l/)
$(N4L_TOOLS):
	@echo "Building $@..."
	cd cmd/n4l && go build -o ../../$@ $@.go

# Search tools (in cmd/search/)
$(SEARCH_TOOLS):
	@echo "Building $@..."
	cd cmd/search && go build -o ../../$@ $@.go

# Examples (in cmd/examples/)
$(EXAMPLES):
	@echo "Building $@..."
	cd cmd/examples && go build -o ../../$@ $@.go

# Server (in cmd/server/)
$(SERVER):
	@echo "Building $@..."
	cd cmd/server && go build -o ../../$@ .

# Demo POCs - DSL and DB connectivity tests (in cmd/demo_pocs/)
$(DEMO_POCS):
	@echo "Building $@..."
	cd cmd/demo_pocs && go build -o ../../$@ $@.go

# Development targets
dev: $(SERVER)
	@echo "Starting development server..."
	./$(SERVER)

# Test DSL parser (N4L language tests)
test-dsl: N4L
	@echo "Running DSL parser tests..."
	cd tests && make

# Test DB connectivity and DSL processing
test-db: $(DEMO_POCS)
	@echo "Testing database connectivity..."
	@echo "Running coarse grain search tests..."
	./search_coarse_grain
	@echo "Running node retrieval tests..."
	./dotest_getnodes

# Process example N4L files into database
load-examples: N4L
	@echo "Loading N4L example files into database..."
	cd examples && make

# Run all tests
test: test-dsl test-db
	@echo "All tests completed"

# Full test cycle: build, load examples, run tests
test-full: all load-examples test
	@echo "Full test cycle completed"

# Clean target
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(ALL_TARGETS)
	rm -f *~ cmd/*~ cmd/*/~ 
	cd tests && make clean
	cd examples && make clean

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Show help
help:
	@echo "SSTorytime Build System"
	@echo ""
	@echo "Available targets:"
	@echo "  all          - Build all executables"
	@echo "  dev          - Build and start development server" 
	@echo "  test-dsl     - Run DSL parser tests"
	@echo "  test-db      - Test database connectivity"
	@echo "  load-examples- Load N4L files into database"
	@echo "  test         - Run DSL and DB tests"
	@echo "  test-full    - Build, load examples, and test everything"
	@echo "  clean        - Remove build artifacts"
	@echo "  deps         - Install/update dependencies"
	@echo "  help         - Show this help"
	@echo ""
	@echo "Executables:"
	@echo "  Core tools:  $(TOOLS)"
	@echo "  N4L tools:   $(N4L_TOOLS)" 
	@echo "  Search:      $(SEARCH_TOOLS)"
	@echo "  Examples:    $(EXAMPLES)"
	@echo "  Server:      $(SERVER)"
	@echo "  DSL/DB tests:$(DEMO_POCS)"

.PHONY: all dev test-dsl test-db load-examples test test-full clean deps help $(ALL_TARGETS)

