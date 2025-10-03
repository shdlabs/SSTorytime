#!/usr/bin/env python3
"""
SSTorytime API Test Suite - SearchN4L Endpoint Tests

This test suite validates the /searchN4L endpoint functionality
against the actual implementation running on localhost:8080.

Tests are now driven by JSON configuration files to make it easier
for non-Python developers to add new tests.
"""

import requests
import json
import time
from typing import Dict, Any, List, Optional
from dataclasses import dataclass
from pathlib import Path

# Configuration
BASE_URL = "http://localhost:8080"
SEARCHN4L_ENDPOINT = f"{BASE_URL}/searchN4L"
CONFIG_FILE = Path(__file__).parent / "searchn4l_tests.json"

@dataclass
class TestResult:
    """Results of an individual test case"""
    test_name: str
    success: bool
    response_code: int
    response_time: float
    response_data: Optional[Dict[str, Any]]
    error_message: Optional[str] = None

class SearchN4LAPITester:
    """Test suite for the SearchN4L API endpoint - now JSON-driven"""
    
    def __init__(self):
        self.results: List[TestResult] = []
        self.session = requests.Session()
        self.config = self.load_test_config()
        
    def load_test_config(self) -> Dict[str, Any]:
        """Load test configuration from JSON file"""
        try:
            with open(CONFIG_FILE, 'r', encoding='utf-8') as f:
                return json.load(f)
        except FileNotFoundError:
            print(f"❌ Configuration file not found: {CONFIG_FILE}")
            print("   Using fallback configuration")
            return self.get_fallback_config()
        except json.JSONDecodeError as e:
            print(f"❌ Error parsing configuration file: {e}")
            print("   Using fallback configuration")
            return self.get_fallback_config()
    
    def get_fallback_config(self) -> Dict[str, Any]:
        """Provide fallback configuration if JSON file is unavailable"""
        return {
            "test_categories": [
                {
                    "id": "basic_searches",
                    "name": "Basic Search Tests",
                    "tests": [
                        {
                            "id": "basic_word_search",
                            "name": "Basic word search",
                            "parameters": {"name": "test"},
                            "expected_status": 200,
                            "validation": {"response_type": "text/html", "contains_text": True}
                        }
                    ]
                }
            ],
            "quick_test_categories": ["basic_searches"],
            "comprehensive_test_categories": ["basic_searches"]
        }
        
    def run_test(self, test_name: str, params: Dict[str, str], method: str = "POST") -> TestResult:
        """Execute a single test case"""
        print(f"🧪 Running test: {test_name}")
        
        start_time = time.time()
        try:
            if method.upper() == "POST":
                response = self.session.post(SEARCHN4L_ENDPOINT, data=params, timeout=30)
            else:
                response = self.session.get(SEARCHN4L_ENDPOINT, params=params, timeout=30)
                
            response_time = time.time() - start_time
            
            # Try to parse JSON response
            try:
                response_data = response.json() if response.text else None
            except json.JSONDecodeError:
                # Handle non-JSON responses
                response_data = {"raw_response": response.text}
            
            success = response.status_code == 200
            if not success:
                print(f"   ❌ Failed with status {response.status_code}")
            else:
                print(f"   ✅ Success ({response_time:.3f}s)")
                
            result = TestResult(
                test_name=test_name,
                success=success,
                response_code=response.status_code,
                response_time=response_time,
                response_data=response_data,
                error_message=None if success else f"HTTP {response.status_code}"
            )
            
        except Exception as e:
            response_time = time.time() - start_time
            print(f"   ❌ Exception: {str(e)}")
            result = TestResult(
                test_name=test_name,
                success=False,
                response_code=0,
                response_time=response_time,
                response_data=None,
                error_message=str(e)
            )
        
        self.results.append(result)
        return result
    
    def run_json_test(self, test_config: Dict[str, Any]) -> TestResult:
        """Execute a test case from JSON configuration"""
        test_name = test_config.get("name", "Unnamed test")
        print(f"🧪 Running test: {test_name}")
        
        # Extract test parameters
        params = test_config.get("parameters", {})
        method = test_config.get("method", "POST")
        headers = test_config.get("headers", {})
        expected_status = test_config.get("expected_status", 200)
        
        start_time = time.time()
        try:
            # Prepare request
            if method.upper() == "POST":
                response = self.session.post(SEARCHN4L_ENDPOINT, data=params, headers=headers, timeout=30)
            else:
                response = self.session.get(SEARCHN4L_ENDPOINT, params=params, headers=headers, timeout=30)
                
            response_time = time.time() - start_time
            
            # Try to parse JSON response
            try:
                response_data = response.json() if response.text else None
            except json.JSONDecodeError:
                # Handle non-JSON responses
                response_data = {"raw_response": response.text}
            
            # Validate response according to test configuration
            success = self.validate_response(response, response_data, test_config)
            
            if not success:
                print(f"   ❌ Failed with status {response.status_code}")
            else:
                print(f"   ✅ Success ({response_time:.3f}s)")
                
            result = TestResult(
                test_name=test_name,
                success=success,
                response_code=response.status_code,
                response_time=response_time,
                response_data=response_data,
                error_message=None if success else f"HTTP {response.status_code}"
            )
            
        except Exception as e:
            response_time = time.time() - start_time
            print(f"   ❌ Exception: {str(e)}")
            result = TestResult(
                test_name=test_name,
                success=False,
                response_code=0,
                response_time=response_time,
                response_data=None,
                error_message=str(e)
            )
        
        self.results.append(result)
        return result
    
    def validate_response(self, response, response_data: Optional[Dict], test_config: Dict[str, Any]) -> bool:
        """Validate response according to test configuration"""
        validation = test_config.get("validation", {})
        expected_status = test_config.get("expected_status", 200)
        
        # Check status code
        if response.status_code != expected_status:
            # Special handling for tests that should fail
            if validation.get("should_fail", False):
                return response.status_code >= 400
            return False
        
        # Check response type
        expected_type = validation.get("response_type")
        if expected_type:
            content_type = response.headers.get("Content-Type", "")
            if expected_type not in content_type:
                return False
        
        # Check minimum response size
        min_size = validation.get("min_response_size", 0)
        if len(response.text) < min_size:
            return False
        
        # Check if response contains text
        if validation.get("contains_text", False):
            if not response.text or len(response.text.strip()) == 0:
                return False
        
        # Check CORS headers
        if validation.get("has_cors_headers", False):
            cors_headers = ["Access-Control-Allow-Origin", "Access-Control-Allow-Methods"]
            if not any(header in response.headers for header in cors_headers):
                return False
        
        return True
    
    def run_category_tests(self, category_id: str):
        """Run all tests in a specific category"""
        category = next((cat for cat in self.config["test_categories"] if cat["id"] == category_id), None)
        if not category:
            print(f"❌ Category '{category_id}' not found in configuration")
            return
        
        print(f"\n📁 Category: {category.get('name', category_id)}")
        
        for test_config in category.get("tests", []):
            self.run_json_test(test_config)
    
    def test_basic_searches(self):
        """Test Category 1: Basic search functionality"""
        print("\n📁 Category 1: Basic Search Tests")
        
        # Test 1.1: Simple word search
        self.run_test(
            "Basic word search",
            {"name": "test"}
        )
        
        # Test 1.2: Search for concept
        self.run_test(
            "Concept search", 
            {"name": "sperm whale"}
        )
        
        # Test 1.3: Search in specific chapter
        self.run_test(
            "Chapter-specific search",
            {"name": "kubernetes \\\\chapter \"Kubernetes notes\""}
        )
        
        # Test 1.4: Empty search (should default to reminders)
        self.run_test(
            "Empty search (reminders)",
            {"name": ""}
        )
        
        # Test 1.5: Help query
        self.run_test(
            "Help query",
            {"name": "\\\\help"}
        )
        
        # Test 1.6: Direct node pointer access
        self.run_test(
            "Direct node access",
            {"name": "", "nclass": "5", "ncptr": "1957", "chapcontext": ""}
        )
    
    def test_search_modifiers(self):
        """Test Category 2: Search with modifiers"""
        print("\n📁 Category 2: Search Modifier Tests")
        
        # Test 2.1: Limit results
        self.run_test(
            "Limited search results",
            {"name": "test \\\\limit 3"}
        )
        
        # Test 2.2: Range/depth modifier
        self.run_test(
            "Range modifier",
            {"name": "whale \\\\range 5"}
        )
        
        # Test 2.3: Context search
        self.run_test(
            "Context search",
            {"name": "whale \\\\context \"sperm whale\""}
        )
        
        # Test 2.4: Chapter contents
        self.run_test(
            "Chapter contents",
            {"name": "\\\\chapter \"Kubernetes notes\" \\\\contents"}
        )
        
        # Test 2.5: Table of contents
        self.run_test(
            "Table of contents",
            {"name": "\\\\toc"}
        )
        
        # Test 2.6: Statistics query
        self.run_test(
            "Statistics query",
            {"name": "test \\\\stats"}
        )
    
    def test_path_solving(self):
        """Test Category 3: Path solving queries"""
        print("\n📁 Category 3: Path Solving Tests")
        
        # Test 3.1: From-to path solving
        self.run_test(
            "Path from-to",
            {"name": "\\\\from testing \\\\to deployment"}
        )
        
        # Test 3.2: Path with chapter context
        self.run_test(
            "Path with chapter",
            {"name": "\\\\from API \\\\to server \\\\chapter \"Kubernetes notes\""}
        )
        
        # Test 3.3: Path with arrow constraints
        self.run_test(
            "Path with arrows",
            {"name": "\\\\from whale \\\\to testing \\\\arrow \"follows on from\""}
        )
        
        # Test 3.4: Path with depth limit
        self.run_test(
            "Path with depth",
            {"name": "\\\\from test \\\\to whale \\\\depth 3"}
        )
    
    def test_sequential_searches(self):
        """Test Category 4: Sequential and story searches"""
        print("\n📁 Category 4: Sequential Search Tests")
        
        # Test 4.1: Sequence search
        self.run_test(
            "Sequence search",
            {"name": "whale \\\\sequence"}
        )
        
        # Test 4.2: Story search
        self.run_test(
            "Story search", 
            {"name": "test \\\\story"}
        )
        
        # Test 4.3: Sequence with limit
        self.run_test(
            "Limited sequence",
            {"name": "kubernetes \\\\seq \\\\limit 5"}
        )
        
        # Test 4.4: Story with chapter
        self.run_test(
            "Chapter story",
            {"name": "\\\\story \\\\chapter \"Samples from Darwin.dat\""}
        )
    
    def test_browsing_queries(self):
        """Test Category 5: Chapter and context browsing"""
        print("\n📁 Category 5: Browsing Tests")
        
        # Test 5.1: Browse notes
        self.run_test(
            "Browse notes",
            {"name": "\\\\notes"}
        )
        
        # Test 5.2: Browse specific chapter
        self.run_test(
            "Browse chapter",
            {"name": "\\\\browse \\\\chapter \"Kubernetes notes\""}
        )
        
        # Test 5.3: Page-specific query
        self.run_test(
            "Page query",
            {"name": "\\\\page 1 \\\\chapter \"Kubernetes notes\""}
        )
        
        # Test 5.4: Context within chapter
        self.run_test(
            "Chapter context",
            {"name": "\\\\chapter \"Samples from Darwin.dat\" \\\\context \"species\""}
        )
    
    def test_special_queries(self):
        """Test Category 6: Special queries and edge cases"""
        print("\n📁 Category 6: Special Queries")
        
        # Test 6.1: Remind query
        self.run_test(
            "Remind query",
            {"name": "\\\\remind"}
        )
        
        # Test 6.2: About query
        self.run_test(
            "About query",
            {"name": "\\\\about test"}
        )
        
        # Test 6.3: Arrow-specific search
        self.run_test(
            "Arrow search",
            {"name": "\\\\arrow \"follows on from\""}
        )
        
        # Test 6.4: Multiple search terms
        self.run_test(
            "Multiple terms",
            {"name": "test whale kubernetes"}
        )
        
        # Test 6.5: Complex query combination
        self.run_test(
            "Complex query",
            {"name": "whale \\\\chapter \"Samples from MobyDick.dat\" \\\\context \"sperm whale\" \\\\limit 10"}
        )
    
    def test_http_methods(self):
        """Test Category 7: HTTP method variations"""
        print("\n📁 Category 7: HTTP Method Tests")
        
        # Test 7.1: GET method
        self.run_test(
            "GET method search",
            {"name": "test"},
            method="GET"
        )
        
        # Test 7.2: POST method (default)
        self.run_test(
            "POST method search",
            {"name": "test"},
            method="POST"
        )
    
    def test_last_seen_updates(self):
        """Test Category 8: Last seen functionality"""
        print("\n📁 Category 8: Last Seen Tests")
        
        # Test 8.1: Update last seen section
        self.run_test(
            "Update last seen section",
            {"name": "\\\\lastnptr", "chapcontext": "Kubernetes notes"}
        )
        
        # Test 8.2: Update last seen node pointer
        self.run_test(
            "Update last seen node",
            {"name": "\\\\lastnptr", "nclass": "5", "ncptr": "1957", "chapcontext": "test"}
        )
    
    def test_error_conditions(self):
        """Test Category 9: Error handling"""
        print("\n📁 Category 9: Error Handling Tests")
        
        # Test 9.1: Invalid HTTP method
        try:
            response = self.session.put(SEARCHN4L_ENDPOINT, data={"name": "test"})
            result = TestResult(
                test_name="Invalid HTTP method",
                success=response.status_code == 405,  # Method Not Allowed expected
                response_code=response.status_code,
                response_time=0.0,
                response_data=None,
                error_message=None if response.status_code == 405 else f"Expected 405, got {response.status_code}"
            )
            self.results.append(result)
            print(f"🧪 Running test: Invalid HTTP method")
            print(f"   {'✅' if result.success else '❌'} {'Success' if result.success else f'Failed: {result.error_message}'}")
        except Exception as e:
            print(f"🧪 Running test: Invalid HTTP method")
            print(f"   ❌ Exception: {str(e)}")
        
        # Test 9.2: Very long query
        self.run_test(
            "Very long query",
            {"name": "test " * 1000}
        )
        
        # Test 9.3: Special characters
        self.run_test(
            "Special characters",
            {"name": "test@#$%^&*()"}
        )
        
        # Test 9.4: Unicode characters
        self.run_test(
            "Unicode characters",
            {"name": "测试 кириллица العربية"}
        )
    
    def run_quick_tests(self):
        """Execute quick test categories (from JSON config)"""
        print("🚀 Starting SSTorytime SearchN4L API Test Suite")
        print(f"🎯 Target endpoint: {SEARCHN4L_ENDPOINT}")
        print("=" * 60)
        
        quick_categories = self.config.get("quick_test_categories", [])
        for category_id in quick_categories:
            self.run_category_tests(category_id)
        
        # Generate summary report
        self.generate_report()
    
    def run_comprehensive_tests(self):
        """Execute all test categories (from JSON config)"""
        print("🚀 Starting SSTorytime SearchN4L API Test Suite")
        print(f"🎯 Target endpoint: {SEARCHN4L_ENDPOINT}")
        print("=" * 60)
        
        comprehensive_categories = self.config.get("comprehensive_test_categories", [])
        for category_id in comprehensive_categories:
            self.run_category_tests(category_id)
        
        # Generate summary report
        self.generate_report()
    
    def run_all_tests(self):
        """Execute all test categories (alias for comprehensive tests)"""
        self.run_comprehensive_tests()
    
    def generate_report(self):
        """Generate comprehensive test report"""
        print("\n" + "=" * 60)
        print("📊 TEST SUMMARY REPORT")
        print("=" * 60)
        
        total_tests = len(self.results)
        passed_tests = sum(1 for r in self.results if r.success)
        failed_tests = total_tests - passed_tests
        
        print(f"Total Tests: {total_tests}")
        print(f"✅ Passed: {passed_tests}")
        print(f"❌ Failed: {failed_tests}")
        print(f"Success Rate: {(passed_tests/total_tests)*100:.1f}%")
        
        if self.results:
            avg_response_time = sum(r.response_time for r in self.results) / len(self.results)
            print(f"Average Response Time: {avg_response_time:.3f}s")
            
            fastest = min(self.results, key=lambda r: r.response_time)
            slowest = max(self.results, key=lambda r: r.response_time)
            print(f"Fastest: {fastest.test_name} ({fastest.response_time:.3f}s)")
            print(f"Slowest: {slowest.test_name} ({slowest.response_time:.3f}s)")
        
        # Failed tests details
        if failed_tests > 0:
            print("\n❌ FAILED TESTS:")
            for result in self.results:
                if not result.success:
                    print(f"   • {result.test_name}: {result.error_message or f'HTTP {result.response_code}'}")
        
        # Response type analysis
        print("\n📈 RESPONSE TYPE ANALYSIS:")
        response_types = {}
        for result in self.results:
            if result.success and result.response_data:
                resp_type = result.response_data.get("Response", "Unknown")
                response_types[resp_type] = response_types.get(resp_type, 0) + 1
        
        for resp_type, count in response_types.items():
            print(f"   • {resp_type}: {count} responses")
        
        # Save detailed results to file
        self.save_detailed_results()
    
    def save_detailed_results(self):
        """Save detailed test results to JSON file"""
        timestamp = time.strftime("%Y%m%d_%H%M%S")
        filename = f"searchn4l_test_results_{timestamp}.json"
        
        # Prepare data for JSON serialization
        results_data = []
        for result in self.results:
            results_data.append({
                "test_name": result.test_name,
                "success": result.success,
                "response_code": result.response_code,
                "response_time": result.response_time,
                "response_data": result.response_data,
                "error_message": result.error_message
            })
        
        report_data = {
            "test_suite": "SearchN4L API Tests",
            "endpoint": SEARCHN4L_ENDPOINT,
            "timestamp": timestamp,
            "summary": {
                "total_tests": len(self.results),
                "passed": sum(1 for r in self.results if r.success),
                "failed": sum(1 for r in self.results if not r.success),
                "success_rate": (sum(1 for r in self.results if r.success) / len(self.results)) * 100 if self.results else 0
            },
            "results": results_data
        }
        
        try:
            with open(filename, 'w', encoding='utf-8') as f:
                json.dump(report_data, f, indent=2, ensure_ascii=False)
            print(f"\n💾 Detailed results saved to: {filename}")
        except Exception as e:
            print(f"\n⚠️  Failed to save results: {e}")

def main():
    """Main entry point"""
    tester = SearchN4LAPITester()
    
    try:
        # Quick connectivity test
        response = requests.get(f"{BASE_URL}/status", timeout=5)
        if response.status_code != 200:
            print(f"❌ Server not accessible at {BASE_URL}")
            print(f"   Status: {response.status_code}")
            return
        
        print(f"✅ Server connectivity confirmed")
        
        # Run all tests
        tester.run_all_tests()
        
    except requests.exceptions.ConnectionError:
        print(f"❌ Cannot connect to server at {BASE_URL}")
        print("   Make sure the SSTorytime server is running on port 8080")
    except Exception as e:
        print(f"❌ Unexpected error: {e}")

if __name__ == "__main__":
    main()