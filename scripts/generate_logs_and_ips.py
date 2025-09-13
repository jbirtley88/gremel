#!/usr/bin/env python3
"""
Generate two datasets:

1) ipaddresses.csv
   - Columns: ip,datacenter
   - 25% of IPs in 'datacenter1'
   - Remaining IPs distributed among 'datacenter2', 'datacenter3', 'datacenter4'

2) weblogs.log (Apache Combined Log Format + trailing latency field)
   - ~30% of requests are 'GET /api/foo'
   - Final field is the request latency in milliseconds (integer), with no label or 'ms' suffix
   - For GET /api/foo:
       * If the client IP belongs to datacenter1, latency > 2000ms
       * If the client IP belongs to other datacenters, latency < 1000ms
   - Other requests have realistic distributions and include a latency value as the final field too

Usage:
  python generate_logs_and_ips.py --ips 1000 --lines 10000 --seed 42

Outputs:
  - ipaddresses.csv
  - weblogs.log
"""
import argparse
import csv
import random
from datetime import datetime, timedelta, timezone

DATACENTERS = ["datacenter1", "datacenter2", "datacenter3", "datacenter4"]

# Browsers and clients for realism
USER_AGENTS = [
    "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:128.0) Gecko/20100101 Firefox/128.0",
    "curl/8.7.1",
    "python-requests/2.31.0",
    "Go-http-client/2.0",
]

REFERERS = [
    "-",
    "https://www.google.com/",
    "https://www.bing.com/search?q=example",
    "https://news.ycombinator.com/",
    "https://twitter.com/",
    "https://example.com/",
    "https://m.example.com/",
    "https://docs.example.com/guide",
]

OTHER_PATHS = [
    "/", "/index.html", "/about", "/contact",
    "/api/bar", "/api/baz", "/api/v1/items", "/api/v1/items/123",
    "/login", "/logout", "/signup",
    "/products/123", "/products/456", "/search?q=widgets",
    "/static/app.js", "/static/style.css", "/static/logo.png",
]

METHODS_OTHER = ["GET", "POST", "PUT", "DELETE"]

STATUSES_COMMON = [200, 200, 200, 200, 200, 201, 204, 301, 302, 304, 400, 401, 403, 404, 500, 502, 503]

# Month abbreviations for CLF timestamp
MONTH_ABBR = ["Jan", "Feb", "Mar", "Apr", "May", "Jun",
              "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]


def rand_ipv4():
    """
    Generate a plausible IPv4 in a mix of private and public ranges.
    Avoid .0 and .255 in each octet to reduce broadcast/network edge cases.
    """
    choice = random.random()
    if choice < 0.4:
        # 10.0.0.0/8
        return f"10.{random.randint(1, 254)}.{random.randint(1, 254)}.{random.randint(1, 254)}"
    elif choice < 0.7:
        # 172.16.0.0/12
        return f"172.{random.randint(16, 31)}.{random.randint(1, 254)}.{random.randint(1, 254)}"
    elif choice < 0.9:
        # 192.168.0.0/16
        return f"192.168.{random.randint(1, 254)}.{random.randint(1, 254)}"
    else:
        # A few public blocks
        first_octet = random.choice([8, 20, 23, 52, 64, 96, 100, 104, 128, 129, 130, 131, 132, 151, 155, 172, 185, 203])
        return f"{first_octet}.{random.randint(1, 254)}.{random.randint(1, 254)}.{random.randint(1, 254)}"


def generate_unique_ips(n):
    ips = set()
    while len(ips) < n:
        ips.add(rand_ipv4())
    return list(ips)


def assign_datacenters(ips):
    """
    Assign IPs to datacenters such that:
      - 25% -> datacenter1
      - 75% -> split across datacenter2/3/4 (uniform)
    Returns: dict ip -> datacenter, and a dict of dc -> list of ips
    """
    n = len(ips)
    n_dc1 = int(round(n * 0.25))
    remaining = n - n_dc1
    per_other = remaining // 3
    counts = {
        "datacenter1": n_dc1,
        "datacenter2": per_other,
        "datacenter3": per_other,
        "datacenter4": remaining - 2 * per_other,
    }

    random.shuffle(ips)
    assignment = {}
    dc_ip_lists = {dc: [] for dc in DATACENTERS}
    start = 0
    for dc in DATACENTERS:
        end = start + counts[dc]
        for ip in ips[start:end]:
            assignment[ip] = dc
            dc_ip_lists[dc].append(ip)
        start = end
    return assignment, dc_ip_lists


def clf_timestamp(dt):
    """
    Format datetime into Apache Combined Log time: [10/Oct/2000:13:55:36 -0700]
    """
    day = dt.day
    mon = MONTH_ABBR[dt.month - 1]
    year = dt.year
    timestr = dt.strftime("%H:%M:%S")
    offset = dt.strftime("%z")
    return f"[{day:02d}/{mon}/{year}:{timestr} {offset}]"


def format_combined_log(ip, dt, method, path, status, bytes_sent, referer, ua, http_version="HTTP/1.1", duration_ms=None):
    """
    Apache Combined Log Format with an extra trailing latency field (milliseconds):
    %h %l %u %t "%r" %>s %b "%{Referer}i" "%{User-agent}i" <latency_ms>
    - ident (%l) and authuser (%u) are '-'
    - %b is '-' if zero bytes were sent
    - <latency_ms> is an integer millisecond value with no label or unit suffix
    """
    ts = clf_timestamp(dt)
    request = f"{method} {path} {http_version}"
    size_field = "-" if bytes_sent == 0 else str(bytes_sent)
    base = f'{ip} - - {ts} "{request}" {status} {size_field} "{referer}" "{ua}"'
    if duration_ms is not None:
        base = f"{base} {int(duration_ms)}"
    return base


def random_datetime(start_dt, days_span=7):
    """
    Random datetime within [start_dt, start_dt + days_span)
    Timezones: randomly pick from UTC, +0100, -0500
    """
    delta_seconds = random.randint(0, days_span * 24 * 3600 - 1)
    naive = start_dt + timedelta(seconds=delta_seconds)
    tz = random.choice([timezone.utc, timezone(timedelta(hours=1)), timezone(timedelta(hours=-5))])
    aware = naive.replace(tzinfo=tz)
    return aware


def pick_bytes_sent(status):
    # Typical sizes; 304 and 204 often return 0
    if status in (204, 304):
        return 0
    choices = [
        random.randint(200, 1500),
        random.randint(1501, 5000),
        random.randint(5001, 20000),
        random.randint(20001, 150000),
    ]
    weights = [0.6, 0.25, 0.12, 0.03]
    return random.choices(choices, weights=weights, k=1)[0]


def build_foo_requests(count, dc_ip_lists, start_dt):
    """
    Build exactly `count` GET /api/foo requests.
    For datacenter1 IPs -> latency > 2000ms, others -> latency < 1000ms.
    """
    # Allocate proportionally to IP counts per DC
    total_ips = sum(len(v) for v in dc_ip_lists.values())
    dc_allocation = {}
    allocated = 0
    for i, dc in enumerate(DATACENTERS):
        if i < len(DATACENTERS) - 1:
            share = int(round(count * (len(dc_ip_lists[dc]) / total_ips))) if total_ips else 0
            dc_allocation[dc] = share
            allocated += share
        else:
            dc_allocation[dc] = max(0, count - allocated)

    lines = []
    for dc in DATACENTERS:
        if not dc_ip_lists[dc]:
            continue
        for _ in range(dc_allocation[dc]):
            ip = random.choice(dc_ip_lists[dc])
            dt = random_datetime(start_dt)
            method = "GET"
            path = "/api/foo"
            status = random.choices([200, 201, 500, 502, 503], weights=[70, 5, 10, 7, 8], k=1)[0]
            bytes_sent = pick_bytes_sent(status)
            referer = random.choice(REFERERS)
            ua = random.choice(USER_AGENTS)
            # Latency rule:
            duration_ms = random.randint(2001, 5000) if dc == "datacenter1" else random.randint(50, 999)
            line = format_combined_log(ip, dt, method, path, status, bytes_sent, referer, ua, duration_ms=duration_ms)
            lines.append(line)

    # Adjust if rounding misallocated
    if len(lines) > count:
        lines = random.sample(lines, count)
    elif len(lines) < count:
        deficit = count - len(lines)
        for _ in range(deficit):
            dcs = [dc for dc in DATACENTERS if dc_ip_lists[dc]]
            weights = [len(dc_ip_lists[dc]) for dc in dcs]
            dc = random.choices(dcs, weights=weights, k=1)[0]
            ip = random.choice(dc_ip_lists[dc])
            dt = random_datetime(start_dt)
            status = random.choice([200, 201, 500, 502, 503])
            bytes_sent = pick_bytes_sent(status)
            referer = random.choice(REFERERS)
            ua = random.choice(USER_AGENTS)
            duration_ms = random.randint(2001, 5000) if dc == "datacenter1" else random.randint(50, 999)
            line = format_combined_log(ip, dt, "GET", "/api/foo", status, bytes_sent, referer, ua, duration_ms=duration_ms)
            lines.append(line)
    return lines


def build_other_requests(count, all_ips, start_dt):
    """
    Build the remaining non-foo requests with general distributions.
    """
    lines = []
    for _ in range(count):
        ip = random.choice(all_ips)
        dt = random_datetime(start_dt)
        method = random.choices(METHODS_OTHER, weights=[0.75, 0.2, 0.03, 0.02], k=1)[0]
        path = random.choice(OTHER_PATHS)
        status = random.choice(STATUSES_COMMON)
        bytes_sent = pick_bytes_sent(status)
        referer = random.choice(REFERERS)
        ua = random.choice(USER_AGENTS)
        # Generic latency for other endpoints: 30–1500ms
        duration_ms = random.randint(30, 1500)
        line = format_combined_log(ip, dt, method, path, status, bytes_sent, referer, ua, duration_ms=duration_ms)
        lines.append(line)
    return lines


def write_ip_csv(filename, assignment_ordered):
    with open(filename, "w", newline="") as f:
        writer = csv.writer(f)
        writer.writerow(["ip", "datacenter"])
        for ip, dc in assignment_ordered:
            writer.writerow([ip, dc])


def write_weblogs(filename, lines):
    with open(filename, "w") as f:
        for line in lines:
            f.write(line + "\n")


def main():
    parser = argparse.ArgumentParser(description="Generate ipaddresses.csv and weblogs.log per requested pathology (Apache combined + trailing latency ms).")
    parser.add_argument("--ips", type=int, default=1000, help="Number of unique IP addresses to generate (default: 1000)")
    parser.add_argument("--lines", type=int, default=10000, help="Number of log lines to generate (default: 10000)")
    parser.add_argument("--seed", type=int, default=42, help="Random seed for reproducibility")
    parser.add_argument("--ip-file", default="ipaddresses.csv", help="Output CSV file for IP to datacenter mapping")
    parser.add_argument("--log-file", default="weblogs.log", help="Output Apache combined-format log file with trailing latency field")
    args = parser.parse_args()

    random.seed(args.seed)

    # Generate IPs and assignments
    ips = generate_unique_ips(args.ips)
    assignment_map, dc_ip_lists = assign_datacenters(ips)

    # Write IP CSV (stable order for reproducibility)
    assignment_ordered = sorted(assignment_map.items(), key=lambda kv: kv[0])
    write_ip_csv(args.ip_file, assignment_ordered)

    # Build web logs
    total_lines = args.lines
    foo_count = int(round(total_lines * 0.30))  # ~30%
    other_count = total_lines - foo_count

    # Start range for dates: within last 10 days from now
    now = datetime.now(timezone.utc)
    start_dt = now - timedelta(days=10)

    foo_lines = build_foo_requests(foo_count, dc_ip_lists, start_dt)
    other_lines = build_other_requests(other_count, ips, start_dt)

    all_lines = foo_lines + other_lines
    random.shuffle(all_lines)

    write_weblogs(args.log_file, all_lines)

    # Print a small summary
    dc_counts = {dc: 0 for dc in DATACENTERS}
    for ip in ips:
        dc_counts[assignment_map[ip]] += 1

    print(f"Wrote {args.ip_file} with {len(ips)} IPs.")
    print(f"Datacenter distribution (IPs): {dc_counts}")
    print(f"Wrote {args.log_file} with {len(all_lines)} lines.")
    print(f"Targeted 'GET /api/foo' share: ~30% -> {foo_count} lines; others: {other_count}.")
    print("Note: weblogs.log is Apache Combined plus an extra trailing latency-ms integer field.")

if __name__ == "__main__":
    main()
