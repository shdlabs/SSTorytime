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