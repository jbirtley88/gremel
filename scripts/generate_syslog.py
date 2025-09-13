#!/usr/bin/python3

import random
import datetime

# Host IPs
hosts = [
    "192.168.0.10",
    "10.0.2.5",
    "172.16.5.22",
    "8.8.4.4",
    "203.0.113.55"
]

# Subsystems and sample messages
subsystems = {
    "systemd": [
        "{service}: Sent signal SIGHUP to main process {pid} ({daemon}) on client request.",
        "{service}: Deactivated successfully.",
        "Finished {service} - {description}.",
        "{service}: Consumed {cpu_time:.3f}s CPU time.",
        "{service}: Scheduled restart job, restart counter is at {counter}.",
        "Stopped {service} - {description}.",
        "Started {service} - {description}.",
    ],
    "sshd": [
        "Accepted password for {user} from {src_ip} port {port} ssh2",
        "Failed password for {user} from {src_ip} port {port} ssh2",
        "Connection closed by {src_ip} port {port} [preauth]",
        "Received disconnect from {src_ip} port {port}: {reason} [preauth]",
        "PAM authentication failure for {user} from {src_ip}",
    ],
    "kernel": [
        "[{timestamp}] {level} {message}",
        "[{timestamp}] {level} CPU{cpu}: Core temperature above threshold, cpu clock throttled",
        "[{timestamp}] {level} eth{eth_num}: Link is Down",
        "[{timestamp}] {level} EXT4-fs (sda1): mounted filesystem with ordered data mode",
        "[{timestamp}] {level} usb {usb_num}-1: USB disconnect, device number {dev_num}"
    ],
    "kubelet": [
        "Flag --container-runtime-endpoint has been deprecated, This parameter should be set via the config file specified by the Kubelet's --config flag.",
        "Flag --pod-infra-container-image has been deprecated, will be removed in a future release.",
        "I{date} {time} {pid} server.go:{line}] \"Kubelet version\" kubeletVersion=\"v{version}\"",
        "I{date} {time} {pid} server.go:{line}] \"Golang settings\" GOGC=\"\" GOMAXPROCS=\"\" GOTRACEBACK=\"\"",
        "I{date} {time} {pid} server.go:{line}] \"Client rotation is on, will bootstrap in background\"",
        "E{date} {time} {pid} log.go:{line}] \"RuntimeConfig from runtime service failed\" err=\"rpc error: code = Unimplemented desc = unknown method RuntimeConfig for service runtime.v1.RuntimeService\""
    ],
    "nginx": [
        "worker process started",
        "worker process exited on signal {signal}",
        "reloading configuration",
        "connection closed while waiting for request",
        "client sent invalid method while reading client request line, client: {src_ip}, server: {server_name}"
    ],
    "myapp": [
        "[{level}] {message}",
        "[{level}] User {user} performed {action}",
        "[{level}] API error: {error}",
        "[{level}] DB connection pool size: {pool_size}",
        "[{level}] Debug info: {debug_info}"
    ]
}

# Log levels
levels = ["DEBUG", "INFO", "NOTICE", "WARNING", "ERROR", "CRITICAL", "ALERT", "EMERGENCY"]
systemd_services = [
    "rsyslog.service", "logrotate.service", "kubelet.service", "ssh.service", "nginx.service"
]
users = ["alice", "bob", "carol", "daniel", "eve", "frank"]
daemons = ["rsyslogd", "sshd", "nginx", "kubelet", "python", "java"]
kernel_messages = [
    "Initializing cgroup subsys cpuset",
    "Memory: {mem}K/{total_mem}K available",
    "EXT4-fs (sda1): mounted filesystem with ordered data mode",
    "usb {usb_num}-1: USB disconnect, device number {dev_num}",
    "eth{eth_num}: Link is Down"
]
app_messages = [
    "Started background job {job_id}",
    "Completed task {task_id} successfully",
    "Cache miss for key {key}",
    "API error: {error}",
    "User {user} logged in"
]
nginx_signals = ["HUP", "TERM", "USR1", "USR2"]
actions = ["login", "logout", "update_profile", "delete_account", "upload_file"]
errors = ["timeout", "connection refused", "invalid data", "permission denied", "not found"]

def random_iso8601():
    # Random date in the past month
    base = datetime.datetime(2025, 9, 1)
    delta = datetime.timedelta(seconds=random.randint(0, 3600*24*12))
    dt = base + delta
    microsec = random.randint(0, 999999)
    dt = dt.replace(microsecond=microsec)
    # Random timezone offset (+01:00 or +00:00)
    tz = random.choice(["+01:00", "+00:00"])
    return dt.strftime("%Y-%m-%dT%H:%M:%S.%f")[:-3] + tz

def random_pid():
    return random.randint(100, 2000000)

def random_port():
    return random.randint(1024, 65535)

def random_service():
    return random.choice(systemd_services)

def random_daemon():
    return random.choice(daemons)

def random_counter():
    return random.randint(1, 500000)

def random_cpu_time():
    return round(random.uniform(0.001, 20.0), 3)

def random_line():
    return random.randint(10, 1000)

def random_version():
    return f"{random.randint(1,2)}.{random.randint(0,99)}.{random.randint(0,99)}"

def random_eth_num():
    return random.randint(0, 3)

def random_usb_num():
    return random.randint(1, 7)

def random_dev_num():
    return random.randint(1, 50)

def random_mem():
    return random.randint(1024, 32768)

def random_total_mem():
    return random.randint(32768, 65536)

def random_job_id():
    return f"job{random.randint(1000,9999)}"

def random_task_id():
    return f"task{random.randint(1000,9999)}"

def random_key():
    return f"key_{random.randint(100,999)}"

def random_pool_size():
    return random.randint(1, 100)

def random_debug_info():
    return f"var={random.randint(1,99)}; state={random.choice(['ok','fail'])}"

def random_server_name():
    return random.choice(["localhost", "api.internal", "www.example.com"])

def random_error():
    return random.choice(errors)

def make_log_line():
    timestamp = random_iso8601()
    host = random.choice(hosts)
    subsystem = random.choice(list(subsystems.keys()))
    pid = random_pid()
    level = random.choice(levels)
    user = random.choice(users)
    src_ip = random.choice(hosts)
    port = random_port()
    job_id = random_job_id()
    task_id = random_task_id()
    key = random_key()
    pool_size = random_pool_size()
    debug_info = random_debug_info()
    service = random_service()
    daemon = random_daemon()
    counter = random_counter()
    cpu_time = random_cpu_time()
    line = random_line()
    version = random_version()
    eth_num = random_eth_num()
    usb_num = random_usb_num()
    dev_num = random_dev_num()
    server_name = random_server_name()
    error = random_error()
    action = random.choice(actions)
    mem = random_mem()
    total_mem = random_total_mem()
    signal = random.choice(nginx_signals)
    description = random.choice([
        "Rotate log files.",
        "The Kubernetes Node Agent.",
        "Secure Shell Daemon.",
        "Web Server.",
        "System Logger.",
        "User management."
    ])
    msg_template = random.choice(subsystems[subsystem])
    # Kernel messages may use custom formatting
    if subsystem == "kernel":
        message = msg_template.format(
            timestamp=timestamp,
            level=level,
            cpu=eth_num,
            eth_num=eth_num,
            usb_num=usb_num,
            dev_num=dev_num,
            mem=mem,
            total_mem=total_mem,
            message=random.choice(kernel_messages)
        )
        return f"{timestamp} {host} kernel: {message}"
    elif subsystem == "systemd":
        message = msg_template.format(
            service=service,
            pid=pid,
            daemon=daemon,
            counter=counter,
            cpu_time=cpu_time,
            description=description
        )
        return f"{timestamp} {host} systemd[1]: {message}"
    elif subsystem == "sshd":
        message = msg_template.format(
            user=user,
            src_ip=src_ip,
            port=port,
            reason=error
        )
        return f"{timestamp} {host} sshd[{pid}]: {message}"
    elif subsystem == "kubelet":
        message = msg_template.format(
            date=timestamp[5:10].replace("-", ""),
            time=timestamp[11:19],
            pid=pid,
            line=line,
            version=version
        )
        return f"{timestamp} {host} kubelet[{pid}]: {message}"
    elif subsystem == "nginx":
        message = msg_template.format(
            signal=signal,
            src_ip=src_ip,
            server_name=server_name
        )
        return f"{timestamp} {host} nginx[{pid}]: {message}"
    elif subsystem == "myapp":
        message = msg_template.format(
            level=level,
            message=random.choice(app_messages).format(
                job_id=job_id,
                task_id=task_id,
                key=key,
                error=error,
                user=user
            ),
            user=user,
            action=action,
            error=error,
            pool_size=pool_size,
            debug_info=debug_info
        )
        return f"{timestamp} {host} myapp[{pid}]: {message}"
    else:
        return f"{timestamp} {host} {subsystem}[{pid}]: {level}: unknown message"

def main():
    with open("syslog_sample.log", "w") as f:
        for _ in range(10000):
            f.write(make_log_line() + "\n")

if __name__ == "__main__":
    main()

