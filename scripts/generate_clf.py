#!/usr/bin/env python3

import random
import time
from datetime import datetime, timedelta

# Constants
NUM_LINES = 1000
METHODS = ['GET'] * 80 + ['POST'] * 7 + ['PUT'] * 7 + ['DELETE'] * 6
STATUSES = [200] * 80 + [400, 401, 403, 404, 500, 501, 503]
URL_PATHS = [
    '/index.html', '/home', '/api/v1/resource', '/login', '/logout', '/dashboard',
    '/settings', '/user/profile', '/search', '/products', '/cart', '/checkout'
]
QUERY_PARAMS = [
    'user=123', 'id=456', 'q=test', 'sort=asc', 'page=2', 'lang=en', 'token=abcdef',
    'mode=full', 'type=basic', 'category=books', 'filter=active', 'ref=google'
]

def random_ip():
    return ".".join(str(random.randint(1, 254)) for _ in range(4))

def random_user():
    return random.choice(['-', 'alice', 'bob', 'carol', 'dan', '-'])

def random_datetime():
    # Generate a random time within the last 30 days
    now = datetime.now()
    delta = timedelta(seconds=random.randint(0, 30 * 24 * 60 * 60))
    dt = now - delta
    # Format: "02/Jan/2006:15:04:05 -0700"
    tz_offset = random.choice(['+0000', '-0700', '+0530', '-0400', '+0100'])
    return dt.strftime('%d/%b/%Y:%H:%M:%S ') + tz_offset

def random_url():
    base = random.choice(URL_PATHS)
    if random.random() < 0.4:
        # 40% chance to have query params
        params = "&".join(random.sample(QUERY_PARAMS, random.randint(1, 3)))
        return f"{base}?{params}"
    else:
        return base

def random_status():
    if random.random() < 0.8:
        return 200
    else:
        return random.choice([400, 401, 403, 404, 500, 501, 503])

def random_method():
    return random.choice(METHODS)

def random_size():
    # Simulate payload size: GETs smaller, POST/PUT larger
    method = random_method()
    if method == 'GET':
        return random.randint(200, 5000)
    elif method in ['POST', 'PUT']:
        return random.randint(500, 50000)
    else:
        return random.randint(100, 10000)

def random_latency():
    # Latency in ms
    return round(random.uniform(10, 900), 2)

def generate_log_line():
    ip = random_ip()
    user = random_user()
    timestamp = random_datetime()
    method = random_method()
    url = random_url()
    protocol = "HTTP/1.1"
    status = random_status()
    size = random_size()
    latency = random_latency()
    # Common Log Format: %h %l %u %t \"%r\" %>s %b
    # We'll add latency at the end
    return f'{ip} - {user} [{timestamp}] "{method} {url} {protocol}" {status} {size} {latency}'

def main():
    for _ in range(NUM_LINES):
        print(generate_log_line())

if __name__ == "__main__":
    main()

