# Test Configuration Guide

This guide explains how to add new tests to the SSTorytime API test suite without needing to know Python. All test definitions are stored in JSON configuration files that are easy to read and modify.

## Configuration Files

- **`searchn4l_tests.json`** - Defines all tests for the SearchN4L API endpoint
- **`status_tests.json`** - Defines all tests for the Status API endpoint and static file serving

## Basic Structure

Each configuration file has the following structure:

```json
{
  "description": "Human-readable description of what this file tests",
  "version": "1.0.0",
  "api_endpoint": "/endpoint-name",
  "test_categories": [
    {
      "id": "unique_category_id",
      "name": "Display Name for Category",
      "description": "What this category tests",
      "icon": "📁",
      "tests": [
        {
          "id": "unique_test_id",
          "name": "Display name for test",
          "description": "What this specific test does",
          "parameters": {
            "name": "search query or parameters"
          },
          "expected_status": 200,
          "validation": {
            "response_type": "text/html",
            "contains_text": true,
            "min_response_size": 100
          }
        }
      ]
    }
  ],
  "quick_test_categories": ["category1", "category2"],
  "comprehensive_test_categories": ["category1", "category2", "category3"]
}
```

## Adding a New Test

### Step 1: Choose the Right File
- For SearchN4L queries → Edit `searchn4l_tests.json`
- For server status, static files, CORS → Edit `status_tests.json`

### Step 2: Find the Right Category
Look at existing categories and find the one that best fits your test:
- `basic_searches` - Simple search terms
- `search_modifiers` - Searches with special modifiers like \\limit, \\range
- `path_solving` - Path finding queries (\\from X \\to Y)
- `error_handling` - Tests that should handle errors gracefully

### Step 3: Add Your Test
Add a new test object to the `tests` array in the appropriate category:

```json
{
  "id": "my_new_test",
  "name": "My New Test",
  "description": "Tests something specific I want to verify",
  "parameters": {
    "name": "your search query here"
  },
  "expected_status": 200,
  "validation": {
    "response_type": "text/html",
    "contains_text": true,
    "min_response_size": 50
  }
}
```

### Step 4: Test Configuration Reference

#### Test Properties
- **`id`** (required) - Unique identifier, use lowercase with underscores
- **`name`** (required) - Human-readable test name shown in output
- **`description`** (required) - Explains what the test validates
- **`parameters`** (required) - Query parameters to send to the API
- **`expected_status`** (required) - HTTP status code expected (usually 200)
- **`validation`** (required) - How to validate the response

#### Parameter Examples

**Simple search:**
```json
"parameters": {
  "name": "whale"
}
```

**Search with modifiers:**
```json
"parameters": {
  "name": "test \\\\limit 5"
}
```

**Path solving:**
```json
"parameters": {
  "name": "\\\\from whale \\\\to testing"
}
```

**Direct node access:**
```json
"parameters": {
  "name": "",
  "nclass": "5",
  "ncptr": "1957",
  "chapcontext": ""
}
```

#### Validation Options

**Basic validation:**
```json
"validation": {
  "response_type": "text/html",
  "contains_text": true,
  "min_response_size": 100
}
```

**JSON response validation:**
```json
"validation": {
  "response_type": "application/json",
  "required_fields": ["server_status", "database_status"],
  "field_validations": {
    "available_topics": {
      "type": "array"
    }
  }
}
```

**Error testing:**
```json
"validation": {
  "should_fail": true,
  "expected_error": "404 Not Found"
}
```

**Performance testing:**
```json
"validation": {
  "max_response_time": 5.0,
  "performance_test": true
}
```

#### HTTP Methods
Most tests use GET, but you can specify other methods:

```json
{
  "id": "post_test",
  "name": "POST request test",
  "method": "POST",
  "parameters": {
    "name": "test"
  },
  "expected_status": 200,
  "validation": {
    "response_type": "text/html"
  }
}
```

#### Custom Headers
For CORS or other header testing:

```json
{
  "id": "cors_test",
  "name": "CORS test",
  "headers": {
    "Origin": "http://localhost:3000",
    "Access-Control-Request-Method": "POST"
  },
  "validation": {
    "has_cors_headers": true
  }
}
```

## Adding a New Category

If you need to test something that doesn't fit existing categories:

```json
{
  "id": "my_new_category",
  "name": "My New Category", 
  "description": "Tests for my specific feature",
  "icon": "📁",
  "tests": [
    {
      "id": "first_test",
      "name": "First Test in New Category",
      "description": "Tests the first thing",
      "parameters": {
        "name": "test query"
      },
      "expected_status": 200,
      "validation": {
        "response_type": "text/html",
        "contains_text": true,
        "min_response_size": 50
      }
    }
  ]
}
```

Don't forget to add your new category to the test execution lists at the bottom:

```json
"quick_test_categories": ["basic_searches", "my_new_category"],
"comprehensive_test_categories": ["basic_searches", "search_modifiers", "my_new_category"]
```

## Common Test Patterns

### Testing SearchN4L Commands

**Help command:**
```json
"parameters": {
  "name": "\\\\help"
}
```

**Table of contents:**
```json
"parameters": {
  "name": "\\\\toc"
}
```

**Chapter contents:**
```json
"parameters": {
  "name": "\\\\chapter \"ChapterName\" \\\\contents"
}
```

**Statistics:**
```json
"parameters": {
  "name": "search_term \\\\stats"
}
```

**Browsing:**
```json
"parameters": {
  "name": "search_term \\\\next"
}
```

### Testing Error Conditions

**Malformed queries:**
```json
{
  "id": "malformed_query",
  "name": "Malformed query test",
  "description": "Test server handles malformed syntax gracefully",
  "parameters": {
    "name": "\\\\from test \\\\to"
  },
  "expected_status": 200,
  "validation": {
    "response_type": "text/html",
    "contains_text": true,
    "min_response_size": 30
  }
}
```

**Special characters:**
```json
{
  "id": "special_chars",
  "name": "Special characters test",
  "parameters": {
    "name": "test @#$%^&*()_+-="
  },
  "expected_status": 200,
  "validation": {
    "response_type": "text/html",
    "contains_text": true,
    "min_response_size": 30
  }
}
```

## Running Your Tests

After adding tests to the JSON files, run them with:

```bash
# Test everything
python3 run_api_tests.py

# Test only quick categories  
python3 run_api_tests.py --quick

# Test only SearchN4L
python3 run_api_tests.py --searchn4l

# Test only Status endpoint
python3 run_api_tests.py --status
```

## Tips for Writing Good Tests

1. **Use descriptive names** - Make it clear what the test validates
2. **Start simple** - Add basic tests before complex edge cases
3. **Test edge cases** - Empty inputs, very long inputs, special characters
4. **Test error conditions** - Make sure the server handles errors gracefully
5. **Check response sizes** - Use appropriate `min_response_size` values
6. **Group related tests** - Put similar tests in the same category

## Validation Guidelines

- **`min_response_size`** should be:
  - 100+ for full search results
  - 50+ for help/info responses  
  - 30+ for error messages or empty results
- **`response_type`** should match what the endpoint actually returns:
  - `text/html` for SearchN4L results
  - `application/json` for Status endpoint
  - `text/css` for CSS files
  - `application/javascript` for JS files

## Common Mistakes to Avoid

1. **Duplicate IDs** - Each test must have a unique `id` within its category
2. **Missing commas** - JSON requires commas between array/object elements
3. **Wrong quotes** - Use double quotes `"` for JSON strings, not single quotes `'`
4. **Backslash escaping** - Use `\\\\` in JSON to represent `\\` in the actual query
5. **Missing validation** - Every test needs a `validation` section

## Need Help?

If you're not sure how to structure a test:
1. Look at existing similar tests in the JSON files
2. Copy an existing test and modify it for your needs
3. Run the tests to see if they work as expected
4. Check the test output for any validation errors

The JSON configuration approach makes it easy to add comprehensive tests without touching any Python code!