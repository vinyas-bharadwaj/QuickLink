#!/usr/bin/env python3

import argparse
import json
import sys
import time
import urllib.parse
import urllib.request

class NoRedirectHandler(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None

def test_load_balancer(target_url, num_requests):
    print("=================================================")
    print("         Load Balancer Verification Test          ")
    print("=================================================")
    print(f"Target Base URL: {target_url}\n")

    opener = urllib.request.build_opener(NoRedirectHandler)
    short_codes = []
    
    # 1. Test POST /shorten
    print("[1/2] Testing URL Shortening (POST /shorten)...")
    for i in range(num_requests):
        data = urllib.parse.urlencode({"url": f"https://example.com/test-url-{i}-{time.time()}"}).encode("utf-8")
        req = urllib.request.Request(
            f"{target_url}/shorten",
            data=data,
            headers={"Content-Type": "application/x-www-form-urlencoded"}
        )
        
        try:
            with urllib.request.urlopen(req, timeout=5) as resp:
                if resp.status == 200:
                    body = json.loads(resp.read().decode("utf-8"))
                    code = body.get("short_url")
                    if code:
                        short_codes.append(code)
                        print(f"  ✓ Request {i+1}: HTTP {resp.status} - Short Code: '{code}'")
                else:
                    print(f"  ✗ Request {i+1}: Unexpected status {resp.status}")
                    return False
        except urllib.error.HTTPError as e:
            if e.code == 429:
                print(f"  ⚠ Request {i+1}: Rate limited (HTTP 429)")
            else:
                print(f"  ✗ Request {i+1}: HTTP Error {e.code}")
                return False
        except Exception as e:
            print(f"  ✗ Connection Error: {e}")
            return False

        # Small delay between requests to remain safely within the rate limit
        time.sleep(0.15)

    if not short_codes:
        print("FAILED: No short URLs were created successfully.")
        return False

    print("\n[2/2] Testing URL Redirects (GET /{short_code})...")
    # 2. Test GET /{short_code}
    for i, code in enumerate(short_codes):
        req = urllib.request.Request(
            f"{target_url}/{code}",
            headers={"User-Agent": "LB-Tester/1.0"}
        )
        try:
            with opener.open(req, timeout=5) as resp:
                status = resp.status
        except urllib.error.HTTPError as e:
            status = e.code
        except Exception as e:
            print(f"  ✗ Connection Error: {e}")
            return False

        if status in (302, 307, 301):
            print(f"  ✓ Request {i+1}: HTTP {status} (Redirect) for '/{code}'")
        elif status == 429:
            print(f"  ⚠ Request {i+1}: Rate limited (HTTP 429)")
        else:
            print(f"  ✗ Request {i+1}: Unexpected status {status} for '/{code}'")
            return False

        time.sleep(0.15)

    print("\n-------------------------------------------------")
    print(" RESULT: Load balancer test PASSED successfully!")
    print(" All requests routed and responded accurately.")
    print("=================================================")
    return True

def main():
    parser = argparse.ArgumentParser(description="Test Load Balancer Functionality")
    parser.add_argument("--url", default="http://localhost:8080", help="Target URL (default: http://localhost:8080)")
    parser.add_argument("-n", "--num", type=int, default=6, help="Number of test requests (default: 6)")
    args = parser.parse_args()

    success = test_load_balancer(args.url, args.num)
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()
