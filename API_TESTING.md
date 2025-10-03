# API Testing

The comprehensive API test suite for SSTorytime has been moved to the `API_TEST/` directory.

## Quick Start

```bash
cd API_TEST
python3 run_api_tests.py --quick
```

## Full Documentation

See [API_TEST/API_TESTS_README.md](API_TEST/API_TESTS_README.md) for complete documentation.

## Files in API_TEST/

- **`run_api_tests.py`** - Main test runner
- **`test_searchn4l_api.py`** - SearchN4L endpoint tests  
- **`test_status_api.py`** - Status endpoint and static file tests
- **`API_TESTS_README.md`** - Complete documentation
- **`unified_api_test_report_*.json`** - Test result reports

## Requirements

- Python 3.6+
- `requests` library (`pip install requests`)
- SSTorytime server running on localhost:8080