#!/usr/bin/python3

import random
import datetime

METHODS = ['GET']*8 + ['POST', 'PUT', 'DELETE']
STATUS_CODES = [200]*8 + [400, 401, 403, 404, 500, 501, 503]
URL_PATHS = [
    '/home', '/about', '/api/data', '/login', '/logout', '/products', '/cart', '/checkout',
    '/search', '/blog', '/posts', '/profile', '/settings', '/api/update', '/api/delete'
]
QUERY_PARAMS = [
    '?q=test', '?user=demo', '?sort=desc', '?id=123', '?page=2', '?ref=google', '?cat=books', '?debug=true'
]
USER_AGENTS = [
    'Mozilla/5.0 (Windows NT 10.0; Win64; x64)',
    'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)',
    'curl/7.68.0',
    'Mozilla/5.0 (X11; Linux x86_64)',
    'PostmanRuntime/7.28.4',
    'Googlebot/2.1 (+http://www.google.com/bot.html)',
]
REFERERS = [
    'https://google.com', 'https://bing.com', 'https://github.com', '-', 'https://facebook.com', 'https://twitter.com'
]

def random_ip():
    return '.'.join(str(random.randint(1, 254)) for _ in range(4))

def random_user():
    # 50% chance of being "-"
    return '-' if random.random() < 0.5 else 'user' + str(random.randint(1, 100))

def random_url():
    base = random.choice(URL_PATHS)
    # 40% chance of having query parameters
    if random.random() < 0.4:
        base += random.choice(QUERY_PARAMS)
    return base

def random_datetime():
    now = datetime.datetime.now()
    # Random offset within last 30 days
    delta = datetime.timedelta(seconds=random.randint(0, 2592000))
    dt = now - delta
    # Timezone offset: -0700 or +0000, etc
    tz_offset = random.choice(['-0700', '+0000', '+0200', '+0530'])
    return dt.strftime('%d/%b/%Y:%H:%M:%S ') + tz_offset

def random_size():
    # Most GETs are small (<5kb), POST/PUT/DELETE can be larger
    return random.randint(10, 1000) if random.random() < 0.7 else random.randint(2000, 3000)

def random_latency():
    # In milliseconds, simulating server timings
    return int(round(random.uniform(0.1, 2.0), 3)*1000)

def generate_log_line():
    ip = random_ip()
    user = random_user()
    ts = random_datetime()
    method = random.choice(METHODS)
    url = random_url()
    proto = random.choice(['HTTP/1.0', 'HTTP/1.1', 'HTTP/2.0'])
    status = random.choice(STATUS_CODES)
    size = random_size()
    referer = random.choice(REFERERS)
    ua = random.choice(USER_AGENTS)
    latency = random_latency()
    # Apache combined log format:
    # %h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-Agent}i\" [latency]
    log = (
        f'{ip} - {user} [{ts}] "{method} {url} {proto}" {status} {size} '
        f'"{referer}" "{ua}" {latency}'
    )
    return log

if __name__ == "__main__":
    for _ in range(1000):
        print(generate_log_line())

