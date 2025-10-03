#!/usr/bin/env python3
"""
SSTorytime API Test Suite - Status & Additional Endpoints

This test suite validates the /status endpoint and any other available
endpoints against the actual implementation running on localhost:8080.

Test Categories:
1. Status endpoint functionality
2. Static file serving
3. CORS headers validation  
4. Error handling for non-existent endpoints
5. Server response format validation
"""

import requests
import json
import time
from typing import Dict, Any, List, Optional
from dataclasses import dataclass
from datetime import datetime

# Configuration
BASE_URL = "http://localhost:8080"
STATUS_ENDPOINT = f"{BASE_URL}/status"

@dataclass
class TestResult:
    """Results of an individual test case"""
    test_name: str
    success: bool
    response_code: int
    response_time: float
    response_data: Optional[Dict[str, Any]]
    headers: Optional[Dict[str, str]]
    error_message: Optional[str] = None

class StatusAPITester:
    """Test suite for the Status API and other endpoints"""
    
    def __init__(self):
        self.results: List[TestResult] = []
        self.session = requests.Session()
        
    def run_test(self, test_name: str, url: str, method: str = "GET", 
                 data: Dict[str, str] = None, headers: Dict[str, str] = None) -> TestResult:
        """Execute a single test case"""
        print(f"🧪 Running test: {test_name}")
        
        start_time = time.time()
        try:
            if method.upper() == "GET":
                response = self.session.get(url, params=data, headers=headers, timeout=30)
            elif method.upper() == "POST":
                response = self.session.post(url, data=data, headers=headers, timeout=30)
            elif method.upper() == "OPTIONS":
                response = self.session.options(url, headers=headers, timeout=30)
            else:
                response = self.session.request(method.upper(), url, data=data, headers=headers, timeout=30)
                
            response_time = time.time() - start_time
            
            # Try to parse JSON response
            try:
                response_data = response.json() if response.text else None
            except json.JSONDecodeError:
                # Handle non-JSON responses
                response_data = {"raw_response": response.text[:500]}  # Limit size for HTML/CSS
            
            success = response.status_code in [200, 304]  # 304 for cached resources
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
                headers=dict(response.headers),
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
    
    def test_status_endpoint(self):
        """Test Category 1: Status endpoint functionality"""
        print("\n📁 Category 1: Status Endpoint Tests")
        
        # Test 1.1: Basic status request
        result = self.run_test(
            "Basic status request",
            STATUS_ENDPOINT
        )
        
        # Validate status response structure
        if result.success and result.response_data:
            expected_fields = ["server_status", "database_status", "available_topics", "timestamp"]
            missing_fields = [field for field in expected_fields if field not in result.response_data]
            
            if missing_fields:
                print(f"   ⚠️  Missing fields: {missing_fields}")
            else:
                print(f"   ✅ All expected fields present")
                
            # Validate data types
            if isinstance(result.response_data.get("available_topics"), list):
                print(f"   ✅ Available topics: {len(result.response_data['available_topics'])} chapters")
            else:
                print(f"   ❌ available_topics should be a list")
        
        # Test 1.2: Status with different HTTP methods
        self.run_test(
            "Status POST request",
            STATUS_ENDPOINT,
            method="POST"
        )
        
        self.run_test(
            "Status OPTIONS request", 
            STATUS_ENDPOINT,
            method="OPTIONS"
        )
    
    def test_static_files(self):
        """Test Category 2: Static file serving"""
        print("\n📁 Category 2: Static File Tests")
        
        # Test 2.1: Main HTML page
        self.run_test(
            "Main HTML page",
            BASE_URL
        )
        
        # Test 2.2: Main JavaScript file
        self.run_test(
            "Main JavaScript file",
            f"{BASE_URL}/main.js"
        )
        
        # Test 2.3: CSS file
        self.run_test(
            "CSS stylesheet",
            f"{BASE_URL}/style.css"
        )
        
        # Test 2.4: Non-existent static file
        self.run_test(
            "Non-existent file (404 test)",
            f"{BASE_URL}/nonexistent.js"
        )
    
    def test_cors_headers(self):
        """Test Category 3: CORS headers validation"""
        print("\n📁 Category 3: CORS Header Tests")
        
        # Test 3.1: CORS preflight request
        cors_headers = {
            "Origin": "http://localhost:3000",
            "Access-Control-Request-Method": "POST",
            "Access-Control-Request-Headers": "Content-Type"
        }
        
        result = self.run_test(
            "CORS preflight request",
            f"{BASE_URL}/searchN4L",
            method="OPTIONS",
            headers=cors_headers
        )
        
        # Validate CORS headers in response
        if result.success and result.headers:
            cors_headers_found = []
            expected_cors_headers = [
                "Access-Control-Allow-Origin",
                "Access-Control-Allow-Methods", 
                "Access-Control-Allow-Headers"
            ]
            
            for header in expected_cors_headers:
                if header in result.headers:
                    cors_headers_found.append(header)
                    print(f"   ✅ {header}: {result.headers[header]}")
                else:
                    print(f"   ❌ Missing: {header}")
        
        # Test 3.2: Cross-origin request simulation
        origin_headers = {"Origin": "http://example.com"}
        self.run_test(
            "Cross-origin request",
            STATUS_ENDPOINT,
            headers=origin_headers
        )
    
    def test_error_endpoints(self):
        """Test Category 4: Error handling for non-existent endpoints"""
        print("\n📁 Category 4: Error Endpoint Tests")
        
        # Test 4.1: Non-existent API endpoint
        self.run_test(
            "Non-existent API endpoint",
            f"{BASE_URL}/nonexistent"
        )
        
        # Test 4.2: Malformed endpoint
        self.run_test(
            "Malformed endpoint",
            f"{BASE_URL}//double//slash"
        )
        
        # Test 4.3: Very long URL
        long_path = "/very/long/path/" * 100
        self.run_test(
            "Very long URL",
            f"{BASE_URL}{long_path}"
        )
    
    def test_server_capabilities(self):
        """Test Category 5: Server capabilities and response format"""
        print("\n📁 Category 5: Server Capability Tests")
        
        # Test 5.1: Content-Type headers
        result = self.run_test(
            "JSON Content-Type validation",
            STATUS_ENDPOINT
        )
        
        if result.success and result.headers:
            content_type = result.headers.get("Content-Type", "")
            if "application/json" in content_type:
                print(f"   ✅ Correct JSON Content-Type: {content_type}")
            else:
                print(f"   ⚠️  Unexpected Content-Type: {content_type}")
        
        # Test 5.2: Server identification
        if result.success and result.headers:
            server_header = result.headers.get("Server", "Unknown")
            print(f"   📄 Server: {server_header}")
            
        # Test 5.3: Response encoding
        if result.success and result.headers:
            encoding = result.headers.get("Content-Encoding", "None")
            print(f"   📄 Encoding: {encoding}")
    
    def test_performance_characteristics(self):
        """Test Category 6: Performance and load characteristics"""
        print("\n📁 Category 6: Performance Tests")
        
        # Test 6.1: Multiple rapid requests
        print(f"   🔄 Testing rapid sequential requests...")
        rapid_times = []
        
        for i in range(5):
            result = self.run_test(
                f"Rapid request #{i+1}",
                STATUS_ENDPOINT
            )
            rapid_times.append(result.response_time)
        
        if rapid_times:
            avg_time = sum(rapid_times) / len(rapid_times)
            print(f"   📊 Average response time: {avg_time:.3f}s")
            print(f"   📊 Response time range: {min(rapid_times):.3f}s - {max(rapid_times):.3f}s")
        
        # Test 6.2: Large response handling (search with many results)
        self.run_test(
            "Large response test",
            f"{BASE_URL}/searchN4L",
            method="POST",
            data={"name": "test \\\\limit 50"}
        )
    
    def test_data_validation(self):
        """Test Category 7: Response data validation"""
        print("\n📁 Category 7: Data Validation Tests")
        
        # Get status response for validation
        result = self.run_test(
            "Status data validation",
            STATUS_ENDPOINT
        )
        
        if result.success and result.response_data:
            data = result.response_data
            
            # Test 7.1: Server status values
            server_status = data.get("server_status")
            if server_status == "OK":
                print(f"   ✅ Server status: {server_status}")
            else:
                print(f"   ⚠️  Unexpected server status: {server_status}")
            
            # Test 7.2: Database status values
            db_status = data.get("database_status")
            if db_status == "OK":
                print(f"   ✅ Database status: {db_status}")
            else:
                print(f"   ⚠️  Unexpected database status: {db_status}")
            
            # Test 7.3: Timestamp format validation
            timestamp_str = data.get("timestamp")
            if timestamp_str:
                try:
                    # Try to parse the timestamp
                    timestamp = datetime.fromisoformat(timestamp_str.replace('Z', '+00:00'))
                    print(f"   ✅ Valid timestamp: {timestamp}")
                except ValueError as e:
                    print(f"   ❌ Invalid timestamp format: {e}")
            
            # Test 7.4: Available topics validation
            topics = data.get("available_topics", [])
            print(f"   📄 Available topics count: {len(topics)}")
            
            # Check for expected topics based on our curl test
            expected_topics = ["Kubernetes notes", "SSTorytime help and search examples"]
            found_topics = [topic for topic in expected_topics if topic in topics]
            print(f"   ✅ Found expected topics: {len(found_topics)}/{len(expected_topics)}")
            
            if found_topics != expected_topics:
                missing = [topic for topic in expected_topics if topic not in topics]
                if missing:
                    print(f"   ⚠️  Missing expected topics: {missing}")
    
    def run_all_tests(self):
        """Execute all test categories"""
        print("🚀 Starting SSTorytime Status & Additional API Test Suite")
        print(f"🎯 Target server: {BASE_URL}")
        print("=" * 60)
        
        # Execute all test categories
        self.test_status_endpoint()
        self.test_static_files()
        self.test_cors_headers()
        self.test_error_endpoints()
        self.test_server_capabilities()
        self.test_performance_characteristics()
        self.test_data_validation()
        
        # Generate summary report
        self.generate_report()
    
    def generate_report(self):
        """Generate comprehensive test report"""
        print("\n" + "=" * 60)
        print("📊 STATUS API TEST SUMMARY REPORT")
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
            code_name = {
                200: "OK",
                404: "Not Found", 
                405: "Method Not Allowed",
                500: "Internal Server Error",
                0: "Connection Error"
            }.get(code, "Unknown")
            print(f"   • {code} ({code_name}): {count} responses")
        
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
        filename = f"status_api_test_results_{timestamp}.json"
        
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
            "test_suite": "Status & Additional API Tests",
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
        response = requests.get(STATUS_ENDPOINT, timeout=5)
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