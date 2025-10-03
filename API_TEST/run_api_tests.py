#!/usr/bin/env python3
"""
SSTorytime API Test Suite - Main Test Runner

Comprehensive test runner for all SSTorytime API endpoints.
This script executes all test suites and generates a unified report.

Usage:
    python3 run_api_tests.py [options]

Options:
    --quick         Run only basic tests (faster execution)
    --searchn4l     Run only SearchN4L endpoint tests
    --status        Run only Status/Static endpoint tests
    --verbose       Enable verbose output
    --save-samples  Save sample responses for documentation
    --no-color      Disable colored output
"""

import sys
import argparse
import time
import json
from pathlib import Path
import subprocess
import requests
from typing import Dict, Any, List

# Import our test modules
try:
    from test_searchn4l_api import SearchN4LAPITester
    from test_status_api import StatusAPITester
except ImportError as e:
    print(f"❌ Failed to import test modules: {e}")
    print("   Make sure test_searchn4l_api.py and test_status_api.py are in the same directory")
    sys.exit(1)

# Configuration
BASE_URL = "http://localhost:8080"
COLORS = {
    'RED': '\033[91m',
    'GREEN': '\033[92m', 
    'YELLOW': '\033[93m',
    'BLUE': '\033[94m',
    'MAGENTA': '\033[95m',
    'CYAN': '\033[96m',
    'WHITE': '\033[97m',
    'BOLD': '\033[1m',
    'UNDERLINE': '\033[4m',
    'END': '\033[0m'
}

class SSTorytimeAPITestRunner:
    """Main test runner for all SSTorytime API tests"""
    
    def __init__(self, args):
        self.args = args
        self.use_colors = not args.no_color and sys.stdout.isatty()
        self.start_time = time.time()
        self.test_results = {}
        
    def colorize(self, text: str, color: str) -> str:
        """Apply color to text if colors are enabled"""
        if not self.use_colors:
            return text
        return f"{COLORS.get(color, '')}{text}{COLORS['END']}"
    
    def print_header(self, text: str):
        """Print a formatted header"""
        print("\n" + "=" * 80)
        print(self.colorize(text, 'BOLD'))
        print("=" * 80)
    
    def print_section(self, text: str):
        """Print a formatted section header"""
        print("\n" + self.colorize(f"📋 {text}", 'CYAN'))
        print("-" * 60)
    
    def check_server_availability(self) -> bool:
        """Check if the SSTorytime server is running and accessible"""
        self.print_section("Server Connectivity Check")
        
        try:
            print(f"🔍 Checking server at {BASE_URL}...")
            response = requests.get(f"{BASE_URL}/status", timeout=10)
            
            if response.status_code == 200:
                print(self.colorize("✅ Server is running and accessible", 'GREEN'))
                
                # Get server info
                try:
                    status_data = response.json()
                    topics_count = len(status_data.get('available_topics', []))
                    print(f"   📊 Available topics: {topics_count}")
                    print(f"   💾 Database status: {status_data.get('database_status', 'Unknown')}")
                except:
                    pass
                    
                return True
            else:
                print(self.colorize(f"❌ Server returned status {response.status_code}", 'RED'))
                return False
                
        except requests.exceptions.ConnectionError:
            print(self.colorize("❌ Connection refused - server not running", 'RED'))
            print(f"   Make sure SSTorytime server is running on port 8080")
            print(f"   You can start it with: cd src/server && air")
            return False
        except requests.exceptions.Timeout:
            print(self.colorize("❌ Connection timeout", 'RED'))
            return False
        except Exception as e:
            print(self.colorize(f"❌ Unexpected error: {e}", 'RED'))
            return False
    
    def run_searchn4l_tests(self) -> Dict[str, Any]:
        """Run SearchN4L API tests"""
        self.print_section("SearchN4L API Tests")
        
        try:
            tester = SearchN4LAPITester()
            
            if self.args.quick:
                print("🏃 Running quick SearchN4L tests...")
                tester.run_quick_tests()
            else:
                print("🔍 Running comprehensive SearchN4L tests...")
                tester.run_comprehensive_tests()
            
            # Extract results
            total = len(tester.results)
            passed = sum(1 for r in tester.results if r.success)
            failed = total - passed
            
            results = {
                'total': total,
                'passed': passed,
                'failed': failed,
                'success_rate': (passed / total * 100) if total > 0 else 0,
                'avg_response_time': sum(r.response_time for r in tester.results) / total if total > 0 else 0,
                'details': tester.results
            }
            
            print(f"\n📊 SearchN4L Results: {self.colorize(f'{passed}/{total}', 'GREEN' if failed == 0 else 'YELLOW')} passed")
            
            return results
            
        except Exception as e:
            print(self.colorize(f"❌ SearchN4L tests failed: {e}", 'RED'))
            return {'total': 0, 'passed': 0, 'failed': 1, 'success_rate': 0, 'error': str(e)}
    
    def run_status_tests(self) -> Dict[str, Any]:
        """Run Status and static file tests"""
        self.print_section("Status & Static File Tests")
        
        try:
            tester = StatusAPITester()
            
            if self.args.quick:
                print("🏃 Running quick Status tests...")
                tester.run_quick_tests()
            else:
                print("🔍 Running comprehensive Status tests...")
                tester.run_comprehensive_tests()
            
            # Extract results
            total = len(tester.results)
            passed = sum(1 for r in tester.results if r.success)
            failed = total - passed
            
            results = {
                'total': total,
                'passed': passed,
                'failed': failed,
                'success_rate': (passed / total * 100) if total > 0 else 0,
                'avg_response_time': sum(r.response_time for r in tester.results) / total if total > 0 else 0,
                'details': tester.results
            }
            
            print(f"\n📊 Status Results: {self.colorize(f'{passed}/{total}', 'GREEN' if failed == 0 else 'YELLOW')} passed")
            
            return results
            
        except Exception as e:
            print(self.colorize(f"❌ Status tests failed: {e}", 'RED'))
            return {'total': 0, 'passed': 0, 'failed': 1, 'success_rate': 0, 'error': str(e)}
    
    def save_sample_responses(self):
        """Save sample API responses for documentation"""
        if not self.args.save_samples:
            return
            
        self.print_section("Saving Sample Responses")
        
        samples = {}
        
        try:
            # Get status sample
            print("📄 Capturing status response...")
            response = requests.get(f"{BASE_URL}/status", timeout=10)
            if response.status_code == 200:
                samples['status_endpoint'] = {
                    'url': '/status',
                    'method': 'GET',
                    'response_code': response.status_code,
                    'headers': dict(response.headers),
                    'data': response.json()
                }
            
            # Get search samples
            search_samples = [
                {"name": "test", "description": "Basic search"},
                {"name": "whale \\\\chapter \"Samples from MobyDick.dat\"", "description": "Chapter-specific search"},
                {"name": "\\\\help", "description": "Help query"},
                {"name": "kubernetes \\\\stats", "description": "Statistics query"}
            ]
            
            samples['searchn4l_endpoint'] = []
            
            for sample in search_samples:
                print(f"📄 Capturing SearchN4L response: {sample['description']}...")
                try:
                    response = requests.post(f"{BASE_URL}/searchN4L", data=sample, timeout=15)
                    if response.status_code == 200:
                        samples['searchn4l_endpoint'].append({
                            'description': sample['description'],
                            'url': '/searchN4L',
                            'method': 'POST',
                            'request_data': sample,
                            'response_code': response.status_code,
                            'headers': dict(response.headers),
                            'data': response.json()
                        })
                except Exception as e:
                    print(f"   ⚠️  Failed to capture {sample['description']}: {e}")
            
            # Save samples to file
            timestamp = time.strftime("%Y%m%d_%H%M%S")
            filename = f"api_response_samples_{timestamp}.json"
            
            with open(filename, 'w', encoding='utf-8') as f:
                json.dump(samples, f, indent=2, ensure_ascii=False)
            
            print(self.colorize(f"💾 Sample responses saved to: {filename}", 'GREEN'))
            
        except Exception as e:
            print(self.colorize(f"❌ Failed to save samples: {e}", 'RED'))
    
    def generate_unified_report(self):
        """Generate a unified test report"""
        self.print_header("📊 UNIFIED TEST REPORT")
        
        total_time = time.time() - self.start_time
        
        # Overall statistics
        overall_total = sum(results.get('total', 0) for results in self.test_results.values())
        overall_passed = sum(results.get('passed', 0) for results in self.test_results.values())
        overall_failed = overall_total - overall_passed
        overall_success_rate = (overall_passed / overall_total * 100) if overall_total > 0 else 0
        
        print(f"🕒 Total execution time: {total_time:.2f} seconds")
        print(f"🧪 Total tests: {overall_total}")
        print(f"✅ Passed: {self.colorize(str(overall_passed), 'GREEN')}")
        print(f"❌ Failed: {self.colorize(str(overall_failed), 'RED' if overall_failed > 0 else 'GREEN')}")
        print(f"📈 Success rate: {self.colorize(f'{overall_success_rate:.1f}%', 'GREEN' if overall_success_rate >= 90 else 'YELLOW')}")
        
        # Per-test-suite breakdown
        print("\n📋 Test Suite Breakdown:")
        for suite_name, results in self.test_results.items():
            total = results.get('total', 0)
            passed = results.get('passed', 0)
            failed = results.get('failed', 0)
            success_rate = results.get('success_rate', 0)
            avg_time = results.get('avg_response_time', 0)
            
            status_color = 'GREEN' if failed == 0 else 'YELLOW' if success_rate >= 50 else 'RED'
            
            print(f"   • {suite_name}:")
            print(f"     Tests: {total}, Passed: {passed}, Failed: {failed}")
            print(f"     Success Rate: {self.colorize(f'{success_rate:.1f}%', status_color)}")
            print(f"     Avg Response Time: {avg_time:.3f}s")
            
            if 'error' in results:
                error_msg = results['error']
                print(f"     {self.colorize(f'Error: {error_msg}', 'RED')}")
        
        # Performance analysis
        if overall_total > 0:
            all_response_times = []
            for results in self.test_results.values():
                if 'details' in results:
                    all_response_times.extend([r.response_time for r in results['details']])
            
            if all_response_times:
                print("\n⚡ Performance Analysis:")
                print(f"   Fastest response: {min(all_response_times):.3f}s")
                print(f"   Slowest response: {max(all_response_times):.3f}s")
                print(f"   Average response: {sum(all_response_times)/len(all_response_times):.3f}s")
        
        # Recommendations
        print("\n💡 Recommendations:")
        
        if overall_success_rate >= 95:
            print(self.colorize("   🎉 Excellent! API is working very well.", 'GREEN'))
        elif overall_success_rate >= 80:
            print(self.colorize("   👍 Good! Minor issues may need attention.", 'YELLOW'))
        else:
            print(self.colorize("   ⚠️  Significant issues detected. Review failed tests.", 'RED'))
        
        if overall_failed > 0:
            print(f"   🔍 Review failed tests for potential server issues")
            print(f"   📋 Check server logs for error details")
        
        # Save unified report
        self.save_unified_report(total_time, overall_total, overall_passed, overall_failed)
    
    def save_unified_report(self, total_time: float, overall_total: int, overall_passed: int, overall_failed: int):
        """Save unified report to JSON file"""
        timestamp = time.strftime("%Y%m%d_%H%M%S")
        filename = f"unified_api_test_report_{timestamp}.json"
        
        report_data = {
            "test_suite": "SSTorytime API Comprehensive Test Suite",
            "base_url": BASE_URL,
            "timestamp": timestamp,
            "execution_time": total_time,
            "summary": {
                "total_tests": overall_total,
                "passed": overall_passed,
                "failed": overall_failed,
                "success_rate": (overall_passed / overall_total * 100) if overall_total > 0 else 0
            },
            "test_suites": self.test_results,
            "environment": {
                "python_version": sys.version,
                "platform": sys.platform,
                "args": vars(self.args)
            }
        }
        
        try:
            with open(filename, 'w', encoding='utf-8') as f:
                json.dump(report_data, f, indent=2, ensure_ascii=False, default=str)
            print(f"\n💾 Unified report saved to: {self.colorize(filename, 'CYAN')}")
        except Exception as e:
            print(f"\n⚠️  Failed to save unified report: {e}")
    
    def run_all_tests(self):
        """Execute all configured test suites"""
        self.print_header("🚀 SSTorytime API Test Suite")
        print(f"🎯 Target server: {self.colorize(BASE_URL, 'CYAN')}")
        print(f"⚙️  Mode: {self.colorize('Quick' if self.args.quick else 'Comprehensive', 'YELLOW')}")
        
        # Check server availability first
        if not self.check_server_availability():
            print(self.colorize("\n❌ Cannot proceed with tests - server not available", 'RED'))
            return False
        
        # Save sample responses if requested
        self.save_sample_responses()
        
        # Run test suites based on arguments
        if not self.args.status:  # Run SearchN4L unless explicitly disabled
            self.test_results['SearchN4L API'] = self.run_searchn4l_tests()
        
        if not self.args.searchn4l:  # Run Status unless explicitly disabled
            self.test_results['Status & Static Files'] = self.run_status_tests()
        
        # Generate unified report
        self.generate_unified_report()
        
        # Return overall success
        overall_failed = sum(results.get('failed', 0) for results in self.test_results.values())
        return overall_failed == 0

def parse_arguments():
    """Parse command line arguments"""
    parser = argparse.ArgumentParser(
        description="SSTorytime API Test Suite",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  python3 run_api_tests.py                    # Run all tests
  python3 run_api_tests.py --quick            # Run basic tests only
  python3 run_api_tests.py --searchn4l        # Test only SearchN4L endpoint
  python3 run_api_tests.py --status           # Test only Status endpoint
  python3 run_api_tests.py --save-samples     # Save sample responses
        """
    )
    
    parser.add_argument('--quick', action='store_true',
                       help='Run only basic tests (faster execution)')
    parser.add_argument('--searchn4l', action='store_true',
                       help='Run only SearchN4L endpoint tests')
    parser.add_argument('--status', action='store_true', 
                       help='Run only Status/Static endpoint tests')
    parser.add_argument('--verbose', action='store_true',
                       help='Enable verbose output')
    parser.add_argument('--save-samples', action='store_true',
                       help='Save sample responses for documentation')
    parser.add_argument('--no-color', action='store_true',
                       help='Disable colored output')
    
    return parser.parse_args()

def main():
    """Main entry point"""
    args = parse_arguments()
    
    # Validate argument combinations
    if args.searchn4l and args.status:
        print("❌ Cannot specify both --searchn4l and --status")
        print("   Use one or neither to run specific test suites")
        sys.exit(1)
    
    # Create and run test suite
    runner = SSTorytimeAPITestRunner(args)
    
    try:
        success = runner.run_all_tests()
        
        # Exit with appropriate code
        if success:
            print("\n🎉 All tests completed successfully!")
            sys.exit(0)
        else:
            print("\n⚠️  Some tests failed. Review the report above.")
            sys.exit(1)
            
    except KeyboardInterrupt:
        print("\n\n⚠️  Tests interrupted by user")
        sys.exit(130)
    except Exception as e:
        print(f"\n❌ Unexpected error during test execution: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()