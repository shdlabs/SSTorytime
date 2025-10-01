#!/bin/bash

#******************************************************************
# SSTorytime API Integration Tests
# Tests the JSON API endpoints with curl
#******************************************************************

BASE_URL="http://localhost:8080"
API_ENDPOINT="$BASE_URL/searchN4L"
TIMEOUT=10
VERBOSE=false

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

#******************************************************************
# Helper Functions
#******************************************************************

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[PASS]${NC} $1"
    ((TESTS_PASSED++))
}

log_error() {
    echo -e "${RED}[FAIL]${NC} $1"
    ((TESTS_FAILED++))
}

log_warning() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# Test if server is running
check_server() {
    log_info "Checking if server is running on $BASE_URL..."
    if curl -s --max-time 5 "$BASE_URL" > /dev/null; then
        log_success "Server is responding"
        return 0
    else
        log_error "Server is not responding. Please start the server first."
        exit 1
    fi
}

# Generic API test function
test_api() {
    local test_name="$1"
    local query="$2"
    local expected_response_type="$3"
    local expected_content_check="$4"
    
    ((TESTS_RUN++))
    log_info "Testing: $test_name"
    
    # Make the API call
    local response=$(timeout $TIMEOUT curl -s -X POST -F "name=$query" "$API_ENDPOINT" 2>/dev/null)
    
    if [ -z "$response" ]; then
        log_error "$test_name: No response received (timeout or connection error)"
        return 1
    fi
    
    # Check if response is valid JSON
    if ! echo "$response" | jq . > /dev/null 2>&1; then
        log_error "$test_name: Response is not valid JSON"
        if [ "$VERBOSE" = true ]; then
            echo "Response: $response"
        fi
        return 1
    fi
    
    # Check response type
    local actual_response_type=$(echo "$response" | jq -r '.Response // "null"')
    if [ "$actual_response_type" != "$expected_response_type" ]; then
        log_error "$test_name: Expected response type '$expected_response_type', got '$actual_response_type'"
        return 1
    fi
    
    # Check content
    if [ -n "$expected_content_check" ]; then
        local content_result=$(echo "$response" | jq -r "$expected_content_check" 2>/dev/null)
        if [ "$content_result" = "null" ] || [ "$content_result" = "false" ] || [ -z "$content_result" ]; then
            log_error "$test_name: Content check failed: $expected_content_check"
            if [ "$VERBOSE" = true ]; then
                echo "Content: $(echo "$response" | jq '.Content')"
            fi
            return 1
        fi
    fi
    
    log_success "$test_name: Response type=$actual_response_type, content validated"
    return 0
}

# Test with detailed output
test_api_detailed() {
    local test_name="$1"
    local query="$2"
    
    log_info "=== Detailed Test: $test_name ==="
    log_info "Query: $query"
    
    local response=$(timeout $TIMEOUT curl -s -X POST -F "name=$query" "$API_ENDPOINT" 2>/dev/null)
    
    if [ -n "$response" ]; then
        echo "Response Type: $(echo "$response" | jq -r '.Response // "ERROR"')"
        echo "Content Type: $(echo "$response" | jq -r '.Content | type')"
        echo "Content Keys: $(echo "$response" | jq -r '.Content | if type == "object" then keys else "N/A (not object)" end')"
        if echo "$response" | jq -e '.Content | has("length")' > /dev/null 2>&1; then
            echo "Content Length: $(echo "$response" | jq -r '.Content | length')"
        fi
        echo ""
    else
        log_error "No response received"
    fi
}

#******************************************************************
# Test Cases
#******************************************************************

run_basic_tests() {
    log_info "=== Running Basic API Tests ==="
    
    # Test 1: Simple word search (should return Orbits)
    test_api "Simple word search" "moon" "Orbits" '.Content | length > 0'
    
    # Test 2: Notes view (should return PageMap)
    test_api "Notes view" "\\notes moon" "PageMap" '.Content.Title != null and .Content.Notes != null'
    
    # Test 3: Table of contents (should return TOC)
    test_api "Table of contents" "\\chapters \\limit 5" "TOC" '.Content | length > 0 and .[0].Chapter != null'
    
    # Test 4: Empty query (should return default)
    test_api "Empty query" "" "PageMap" '.Content.Title != null'
    
    # Test 5: Help command
    test_api "Help command" "\\help" "PageMap" '.Content.Title != null'
    
    # Test 6: Stats command  
    test_api "Stats command" "\\stats" "STAT" '.Content != null'
}

run_advanced_tests() {
    log_info "=== Running Advanced API Tests ==="
    
    # Test complex queries
    test_api "Word with context" "moon \\context astronomy" "Orbits" '.Content | length > 0'
    
    test_api "Chapter-specific search" "moon \\chapter \"Moon tidal effect\"" "Orbits" '.Content | length > 0'
    
    test_api "Limited results" "moon \\limit 2" "Orbits" '.Content | length <= 2'
    
    test_api "Notes with limit" "\\notes moon \\limit 3" "PageMap" '.Content.Notes | length <= 3'
}

run_edge_cases() {
    log_info "=== Running Edge Case Tests ==="
    
    # Test non-existent terms
    test_api "Non-existent term" "xyznonsenseword123" "Orbits" true
    
    # Test special characters
    test_api "Special characters" "moon & sun" "Orbits" true
    
    # Test very long query
    test_api "Long query" "$(printf 'moon %.0s' {1..50})" "Orbits" true
}

run_performance_tests() {
    log_info "=== Running Performance Tests ==="
    
    local start_time=$(date +%s.%N)
    
    # Run multiple quick tests
    for i in {1..5}; do
        timeout 5 curl -s -X POST -F "name=moon" "$API_ENDPOINT" > /dev/null
        if [ $? -ne 0 ]; then
            log_warning "Performance test $i: Request timed out or failed"
        fi
    done
    
    local end_time=$(date +%s.%N)
    local duration=$(echo "$end_time - $start_time" | bc)
    
    log_info "Performance test completed: 5 requests in ${duration}s"
}

run_json_structure_tests() {
    log_info "=== Running JSON Structure Validation ==="
    
    # Test each response type has expected structure
    local queries=("moon:Orbits" "\\notes moon:PageMap" "\\chapters \\limit 3:TOC" "\\stats:STAT")
    
    for query_type in "${queries[@]}"; do
        IFS=':' read -r query expected_type <<< "$query_type"
        
        log_info "Validating JSON structure for: $query"
        local response=$(timeout $TIMEOUT curl -s -X POST -F "name=$query" "$API_ENDPOINT" 2>/dev/null)
        
        if [ -n "$response" ]; then
            # Validate required fields
            local has_response=$(echo "$response" | jq -e 'has("Response")' 2>/dev/null)
            local has_content=$(echo "$response" | jq -e 'has("Content")' 2>/dev/null)
            local has_time=$(echo "$response" | jq -e 'has("Time")' 2>/dev/null)
            local has_intent=$(echo "$response" | jq -e 'has("Intent")' 2>/dev/null)
            local has_ambient=$(echo "$response" | jq -e 'has("Ambient")' 2>/dev/null)
            
            if [ "$has_response" = "true" ] && [ "$has_content" = "true" ] && [ "$has_time" = "true" ]; then
                log_success "JSON structure valid for $expected_type"
            else
                log_error "JSON structure invalid for $expected_type"
                echo "  Response field: $has_response"
                echo "  Content field: $has_content" 
                echo "  Time field: $has_time"
            fi
        else
            log_error "No response for $query"
        fi
    done
}

#******************************************************************
# Main Execution
#******************************************************************

show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help      Show this help message"
    echo "  -v, --verbose   Enable verbose output"
    echo "  -s, --server    Server URL (default: $BASE_URL)"
    echo "  -t, --timeout   Request timeout in seconds (default: $TIMEOUT)"
    echo "  --basic         Run only basic tests"
    echo "  --advanced      Run only advanced tests"
    echo "  --edge          Run only edge case tests"
    echo "  --performance   Run only performance tests"
    echo "  --json          Run only JSON structure tests"
    echo "  --detailed      Show detailed output for key tests"
    echo ""
    echo "Examples:"
    echo "  $0                    # Run all tests"
    echo "  $0 --basic --verbose  # Run basic tests with verbose output"
    echo "  $0 --detailed         # Show detailed test information"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_usage
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -s|--server)
            BASE_URL="$2"
            API_ENDPOINT="$BASE_URL/searchN4L"
            shift 2
            ;;
        -t|--timeout)
            TIMEOUT="$2"
            shift 2
            ;;
        --basic)
            RUN_BASIC=true
            shift
            ;;
        --advanced)
            RUN_ADVANCED=true
            shift
            ;;
        --edge)
            RUN_EDGE=true
            shift
            ;;
        --performance)
            RUN_PERFORMANCE=true
            shift
            ;;
        --json)
            RUN_JSON=true
            shift
            ;;
        --detailed)
            RUN_DETAILED=true
            shift
            ;;
        *)
            log_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Main test execution
main() {
    echo "================================================"
    echo "       SSTorytime API Integration Tests"
    echo "================================================"
    echo "Server: $BASE_URL"
    echo "Timeout: ${TIMEOUT}s"
    echo "Verbose: $VERBOSE"
    echo "================================================"
    echo ""
    
    # Check if server is running
    check_server
    echo ""
    
    # Run tests based on flags
    if [ "$RUN_DETAILED" = true ]; then
        test_api_detailed "Basic Search" "moon"
        test_api_detailed "Notes View" "\\notes moon"
        test_api_detailed "Table of Contents" "\\chapters \\limit 3"
        test_api_detailed "Statistics" "\\stats"
    elif [ "$RUN_BASIC" = true ]; then
        run_basic_tests
    elif [ "$RUN_ADVANCED" = true ]; then
        run_advanced_tests
    elif [ "$RUN_EDGE" = true ]; then
        run_edge_cases
    elif [ "$RUN_PERFORMANCE" = true ]; then
        run_performance_tests
    elif [ "$RUN_JSON" = true ]; then
        run_json_structure_tests
    else
        # Run all tests
        run_basic_tests
        echo ""
        run_advanced_tests
        echo ""
        run_edge_cases
        echo ""
        run_json_structure_tests
        echo ""
        run_performance_tests
    fi
    
    # Summary
    echo ""
    echo "================================================"
    echo "              Test Summary"
    echo "================================================"
    echo "Tests Run:    $TESTS_RUN"
    echo "Tests Passed: $TESTS_PASSED"
    echo "Tests Failed: $TESTS_FAILED"
    
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}All tests passed!${NC}"
        exit 0
    else
        echo -e "${RED}Some tests failed.${NC}"
        exit 1
    fi
}

# Check dependencies
if ! command -v jq &> /dev/null; then
    log_error "jq is required but not installed. Please install jq first."
    exit 1
fi

if ! command -v curl &> /dev/null; then
    log_error "curl is required but not installed. Please install curl first."
    exit 1
fi

# Run main function
main