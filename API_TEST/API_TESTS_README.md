# SSTorytime API Test Suite

Comprehensive API testing framework for the SSTorytime server running on localhost:8080.

## 📁 Files

- **`run_api_tests.py`** - Main test runner with unified reporting
- **`test_searchn4l_api.py`** - Comprehensive tests for `/searchN4L` endpoint
- **`test_status_api.py`** - Tests for `/status` endpoint and static files

## 🚀 Quick Start

Make sure the SSTorytime server is running:
```bash
cd src/server
air  # or go run http_server.go
```

Run all tests:
```bash
python3 run_api_tests.py
```

## 📋 Test Categories

### SearchN4L API Tests (`/searchN4L`)
1. **Basic Searches** - Simple word/concept searches
2. **Search Modifiers** - `\\limit`, `\\range`, `\\context`, `\\chapter`, `\\stats`
3. **Path Solving** - `\\from` and `\\to` queries
4. **Sequential Searches** - `\\sequence`, `\\story`
5. **Browsing Queries** - `\\notes`, `\\browse`, `\\page`
6. **Special Queries** - `\\help`, `\\remind`, `\\about`
7. **HTTP Methods** - GET/POST variations
8. **Last Seen Updates** - `\\lastnptr` functionality
9. **Error Handling** - Invalid requests, edge cases

### Status & Static Tests (`/status`, static files)
1. **Status Endpoint** - JSON response validation
2. **Static Files** - HTML, JS, CSS serving
3. **CORS Headers** - Cross-origin request support
4. **Error Endpoints** - 404 handling
5. **Server Capabilities** - Content-Type, headers
6. **Performance** - Response times, load handling
7. **Data Validation** - Response structure verification

## 🎯 Usage Options

```bash
# Run all tests (comprehensive)
python3 run_api_tests.py

# Quick test run (basic tests only)
python3 run_api_tests.py --quick

# Test only SearchN4L endpoint
python3 run_api_tests.py --searchn4l

# Test only Status/Static endpoints  
python3 run_api_tests.py --status

# Save sample responses for documentation
python3 run_api_tests.py --save-samples

# Disable colored output
python3 run_api_tests.py --no-color

# Get help
python3 run_api_tests.py --help
```

## 📊 Expected API Responses

### Status Endpoint (`/status`)
```json
{
  "server_status": "OK",
  "database_status": "OK", 
  "available_topics": ["chapter1", "chapter2", ...],
  "timestamp": "2025-10-03T14:36:19.335237717+03:00"
}
```

### SearchN4L Endpoint (`/searchN4L`)
```json
{
  "Response": "Orbits|PathSolve|Stories|ChapterContexts|Stats|PageMap",
  "Content": [...],
  "Time": "Fri:Hr14:Qu3-Min40_45",
  "Intent": "search_query",
  "Ambient": "N_Autumn, S_Spring, Afternoon, ..."
}
```

## 🔍 Search Query Examples

| Query Type | Example | Description |
|------------|---------|-------------|
| Basic Search | `name=test` | Simple word search |
| Chapter Search | `name=whale \\chapter "Samples from MobyDick.dat"` | Search within chapter |
| Path Solving | `name=\\from testing \\to deployment` | Find paths between concepts |
| Statistics | `name=test \\stats` | Get search statistics |
| Help | `name=\\help` | Get help information |
| Sequence | `name=whale \\sequence` | Sequential/story search |
| Limited Results | `name=test \\limit 5` | Limit result count |

## 📈 Report Generation

Tests automatically generate detailed reports:

- **Console Output** - Real-time test results with colors
- **JSON Reports** - Detailed results saved to timestamped files
- **Unified Report** - Combined results from all test suites
- **Sample Responses** - API response examples (with `--save-samples`)

## 🛠 Requirements

- Python 3.6+
- `requests` library (`pip install requests`)
- SSTorytime server running on localhost:8080

## 🔧 Server Commands

The tests work with the current server implementation. Make sure you have:

1. **Database populated** with some N4L content
2. **Server running** on port 8080:
   ```bash
   cd /home/alex/SSTorytime/src/server
   air  # for development with auto-reload
   # or
   go run http_server.go  # for simple run
   ```

## 📋 Test Results Interpretation

### Success Indicators
- ✅ **Green**: Test passed successfully
- 📊 **Success Rate > 90%**: API working well
- ⚡ **Response Time < 1s**: Good performance

### Warning Indicators  
- ⚠️ **Yellow**: Some tests failed but system functional
- 📊 **Success Rate 80-90%**: Minor issues present
- ⚡ **Response Time > 2s**: Performance concerns

### Error Indicators
- ❌ **Red**: Test failed
- 📊 **Success Rate < 80%**: Significant issues
- 🔍 **Connection Errors**: Server not running/accessible

## 🧪 Adding New Tests

To add new test cases:

1. **SearchN4L tests**: Edit `test_searchn4l_api.py`, add to appropriate category
2. **Status tests**: Edit `test_status_api.py`, add new test method
3. **New endpoints**: Create new test file, import in `run_api_tests.py`

Example new test:
```python
def test_new_feature(self):
    """Test new API feature"""
    self.run_test(
        "New feature test",
        {"name": "new_query_type"}
    )
```

## 🎯 Based on Actual Implementation

These tests are designed to work with the **actual SSTorytime server code** without modifications. They test:

- Real endpoint paths (`/searchN4L`, `/status`)
- Actual request parameters (`name`, `nclass`, `ncptr`, `chapcontext`)
- Expected response formats (JSON with `Response`, `Content`, `Time`, etc.)
- Current search syntax (`\\chapter`, `\\limit`, `\\from`, `\\to`, etc.)
- CORS headers and HTTP methods supported
- Static file serving from embedded filesystem

The tests **fit the code** rather than requiring code changes to fit the tests.