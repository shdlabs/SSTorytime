#!/usr/bin/env python3
"""
SSTorytime API Test Suite - Status Endpoint Tests

This test suite validates the /status endpoint and static file serving
functionality against the actual implementation running on localhost:8080.

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
STATUS_ENDPOINT = f"{BASE_URL}/status"
CONFIG_FILE = Path(__file__).parent / "status_tests.json"

@dataclass
class TestResult:
    """Results of an individual test case"""
    test_name: str
    success: bool
    response_code: int
    response_time: float
    response_data: Optional[Dict[str, Any]]
    headers: Optional[Dict[str, str]] = None
    error_message: Optional[str] = None

class StatusAPITester:
    """Test suite for the Status API endpoint and static files - now JSON-driven"""
    
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
                    "id": "status_endpoint",
                    "name": "Status Endpoint Tests",
                    "tests": [
                        {
                            "id": "basic_status_request",
                            "name": "Basic status request",
                            "url": "/status",
                            "method": "GET",
                            "expected_status": 200,
                            "validation": {"response_type": "application/json"}
                        }
                    ]
                }
            ],
            "quick_test_categories": ["status_endpoint"],
            "comprehensive_test_categories": ["status_endpoint"]
        }
        
    def run_json_test(self, test_config: Dict[str, Any]) -> TestResult:
        """Execute a test case from JSON configuration"""
        test_name = test_config.get("name", "Unnamed test")
        print(f"🧪 Running test: {test_name}")
        
        # Extract test parameters
        url = test_config.get("url", "/status")
        method = test_config.get("method", "GET")
        headers = test_config.get("headers", {})
        params = test_config.get("parameters", {})
        expected_status = test_config.get("expected_status", 200)
        
        # Build full URL
        full_url = f"{BASE_URL}{url}" if url.startswith("/") else url
        if not url.startswith("http") and not url.startswith("/"):
            full_url = f"{BASE_URL}/{url}"
        
        start_time = time.time()
        try:
            # Prepare request
            if method.upper() == "POST":
                response = self.session.post(full_url, data=params, headers=headers, timeout=30)
            elif method.upper() == "OPTIONS":
                response = self.session.options(full_url, headers=headers, timeout=30)
            else:
                response = self.session.get(full_url, params=params, headers=headers, timeout=30)
                
            response_time = time.time() - start_time
            
            # Try to parse JSON response
            try:
                response_data = response.json() if response.text and 'application/json' in response.headers.get('Content-Type', '') else None
            except json.JSONDecodeError:
                response_data = None
            
            # Store response headers
            response_headers = dict(response.headers)
            
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
                headers=response_headers,
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
                headers=None,
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
        
        # Check required fields for JSON responses
        required_fields = validation.get("required_fields", [])
        if required_fields and response_data:
            missing_fields = [field for field in required_fields if field not in response_data]
            if missing_fields:
                print(f"   ⚠️  Missing fields: {missing_fields}")
                return False
            print(f"   ✅ All expected fields present")
        
        # Validate field types
        field_validations = validation.get("field_validations", {})
        if field_validations and response_data:
            for field, field_config in field_validations.items():
                if field in response_data:
                    expected_type = field_config.get("type")
                    if expected_type == "array" and not isinstance(response_data[field], list):
                        return False
                    elif expected_type == "string" and not isinstance(response_data[field], str):
                        return False
                    elif expected_type == "number" and not isinstance(response_data[field], (int, float)):
                        return False
        
        # Check for specific validation requirements
        if validation.get("response_type") == "application/json" and response_data:
            if "available_topics" in response_data and isinstance(response_data["available_topics"], list):
                print(f"   ✅ Available topics: {len(response_data['available_topics'])} chapters")
        
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
        print("🚀 Starting SSTorytime Status API Test Suite")
        print(f"🎯 Target server: {BASE_URL}")
        print("=" * 60)
        
        quick_categories = self.config.get("quick_test_categories", [])
        for category_id in quick_categories:
            self.run_category_tests(category_id)
        
        # Generate summary report
        self.generate_report()
    
    def run_comprehensive_tests(self):
        """Execute all test categories (from JSON config)"""
        print("🚀 Starting SSTorytime Status API Test Suite")
        print(f"🎯 Target server: {BASE_URL}")
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
        
        # HTTP status code analysis
        print("\n📈 HTTP STATUS CODE ANALYSIS:")
        status_codes = {}
        for result in self.results:
            code = result.response_code
            status_codes[code] = status_codes.get(code, 0) + 1
        
        for code, count in sorted(status_codes.items()):
            status_text = {
                200: "OK",
                404: "Not Found", 
                405: "Method Not Allowed",
                500: "Server Error"
            }.get(code, "Unknown")
            print(f"   • {code} ({status_text}): {count} responses")
        
        # Failed tests details
        if failed_tests > 0:
            print("\n❌ FAILED TESTS:")
            for result in self.results:
                if not result.success:
                    print(f"   • {result.test_name}: {result.error_message or f'HTTP {result.response_code}'}")
        
        # Save detailed results to file
        self.save_detailed_results()
    
    def save_detailed_results(self):
        """Save detailed test results to JSON file"""
        timestamp = time.strftime("%Y%m%d_%H%M%S")
        filename = f"status_test_results_{timestamp}.json"
        
        # Prepare data for JSON serialization
        results_data = []
        for result in self.results:
            results_data.append({
                "test_name": result.test_name,
                "success": result.success,
                "response_code": result.response_code,
                "response_time": result.response_time,
                "response_data": result.response_data,
                "headers": result.headers,
                "error_message": result.error_message
            })
        
        report_data = {
            "test_suite": "Status API Tests",
            "base_url": BASE_URL,
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
    tester = StatusAPITester()
    
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