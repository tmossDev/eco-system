#!/usr/bin/env python3
"""Local HTTP proxy for Kubernetes ingress hostnames.

This lets a browser reach local ingress hosts without adding each host to
/etc/hosts. Configure the browser or OS HTTP proxy to this process.
"""

from __future__ import annotations

import argparse
import http.client
import json
import select
import socket
import socketserver
import subprocess
import sys
import threading
import time
from http.server import BaseHTTPRequestHandler
from urllib.parse import urlsplit


HOP_BY_HOP_HEADERS = {
    "connection",
    "keep-alive",
    "proxy-authenticate",
    "proxy-authorization",
    "proxy-connection",
    "te",
    "trailer",
    "transfer-encoding",
    "upgrade",
}


def run_kubectl(args: list[str]) -> str:
    try:
        completed = subprocess.run(
            ["kubectl", *args],
            check=True,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
        )
    except FileNotFoundError:
        raise RuntimeError("kubectl is required but was not found on PATH")
    except subprocess.CalledProcessError as exc:
        message = exc.stderr.strip() or exc.stdout.strip()
        raise RuntimeError(f"kubectl failed: {message}")

    return completed.stdout


def ingress_data(namespace: str) -> tuple[set[str], str | None]:
    raw = run_kubectl(["get", "ingress", "-n", namespace, "-o", "json"])
    payload = json.loads(raw)
    hosts: set[str] = set()
    addresses: list[str] = []

    for item in payload.get("items", []):
        for rule in item.get("spec", {}).get("rules", []):
            host = rule.get("host")
            if host:
                hosts.add(host)

        for entry in item.get("status", {}).get("loadBalancer", {}).get("ingress", []):
            address = entry.get("ip") or entry.get("hostname")
            if address:
                addresses.append(address)

    return hosts, addresses[0] if addresses else None


class ThreadingHTTPServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
    allow_reuse_address = True
    daemon_threads = True


class IngressProxyHandler(BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"
    hosts: set[str] = set()
    target_host = "127.0.0.1"
    http_port = 80
    https_port = 443
    state_lock = threading.Lock()

    @classmethod
    def update_state(cls, hosts: set[str], target_host: str) -> None:
        with cls.state_lock:
            cls.hosts = hosts
            cls.target_host = target_host

    @classmethod
    def state(cls) -> tuple[set[str], str, int, int]:
        with cls.state_lock:
            return set(cls.hosts), cls.target_host, cls.http_port, cls.https_port

    def do_CONNECT(self) -> None:
        hosts, target_host, _, https_port = self.state()
        host, port = split_host_port(self.path, https_port)
        if host not in hosts:
            self.send_error(502, f"Unknown ingress host: {host}")
            return

        try:
            upstream = socket.create_connection(
                (target_host, port if port != 443 else https_port),
                timeout=10,
            )
        except OSError as exc:
            self.send_error(502, f"Could not connect to ingress: {exc}")
            return

        self.send_response(200, "Connection Established")
        self.end_headers()
        tunnel(self.connection, upstream)

    def do_GET(self) -> None:
        self.forward()

    def do_HEAD(self) -> None:
        self.forward()

    def do_POST(self) -> None:
        self.forward()

    def do_PUT(self) -> None:
        self.forward()

    def do_PATCH(self) -> None:
        self.forward()

    def do_DELETE(self) -> None:
        self.forward()

    def do_OPTIONS(self) -> None:
        self.forward()

    def forward(self) -> None:
        hosts, target_host, http_port, https_port = self.state()
        request = urlsplit(self.path)
        if request.scheme and request.netloc:
            original_host, original_port = split_host_port(
                request.netloc,
                443 if request.scheme == "https" else 80,
            )
            path = request.path or "/"
            if request.query:
                path = f"{path}?{request.query}"
        else:
            host_header = self.headers.get("Host", "")
            original_host, original_port = split_host_port(host_header, 80)
            path = self.path

        if original_host not in hosts:
            self.send_error(502, f"Unknown ingress host: {original_host}")
            return

        port = https_port if original_port == 443 else http_port
        body = self.read_body()
        headers = {
            key: value
            for key, value in self.headers.items()
            if key.lower() not in HOP_BY_HOP_HEADERS and key.lower() != "host"
        }
        headers["Host"] = original_host

        connection = http.client.HTTPConnection(target_host, port, timeout=30)
        try:
            connection.request(self.command, path, body=body, headers=headers)
            response = connection.getresponse()
            self.send_response(response.status, response.reason)
            for key, value in response.getheaders():
                if key.lower() not in HOP_BY_HOP_HEADERS:
                    self.send_header(key, value)
            self.end_headers()
            if self.command != "HEAD":
                while True:
                    chunk = response.read(64 * 1024)
                    if not chunk:
                        break
                    self.wfile.write(chunk)
        except OSError as exc:
            self.send_error(502, f"Could not proxy request: {exc}")
        finally:
            connection.close()

    def read_body(self) -> bytes | None:
        length = self.headers.get("Content-Length")
        if not length:
            return None
        return self.rfile.read(int(length))

    def log_message(self, fmt: str, *args: object) -> None:
        sys.stderr.write(f"{self.address_string()} - {fmt % args}\n")


def split_host_port(value: str, default_port: int) -> tuple[str, int]:
    if value.startswith("["):
        host, _, port = value[1:].partition("]:")
        return host, int(port) if port else default_port

    host, sep, port = value.partition(":")
    return host, int(port) if sep and port else default_port


def tunnel(client: socket.socket, upstream: socket.socket) -> None:
    sockets = [client, upstream]
    try:
        while True:
            readable, _, errored = select.select(sockets, [], sockets, 60)
            if errored or not readable:
                break
            for current in readable:
                other = upstream if current is client else client
                data = current.recv(64 * 1024)
                if not data:
                    return
                other.sendall(data)
    finally:
        upstream.close()


def refresh_ingress_state(
    namespace: str,
    extra_hosts: list[str],
    target_host_override: str | None,
    interval: int,
) -> None:
    previous_hosts = set(IngressProxyHandler.hosts)
    previous_target = IngressProxyHandler.target_host

    while True:
        time.sleep(interval)
        try:
            hosts, discovered_target = ingress_data(namespace)
            hosts.update(extra_hosts)
            target_host = target_host_override or discovered_target
            if not hosts or not target_host:
                continue
        except (RuntimeError, json.JSONDecodeError) as exc:
            print(f"Could not refresh ingress hosts: {exc}", file=sys.stderr, flush=True)
            continue

        IngressProxyHandler.update_state(hosts, target_host)
        if hosts != previous_hosts or target_host != previous_target:
            print(
                f"Refreshed {len(hosts)} ingress host(s) to {target_host}",
                flush=True,
            )
            previous_hosts = hosts
            previous_target = target_host


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Proxy local browser traffic to Kubernetes ingress hosts.",
    )
    parser.add_argument("-n", "--namespace", default="eco-test")
    parser.add_argument("--listen-host", default="127.0.0.1")
    parser.add_argument("--listen-port", type=int, default=18080)
    parser.add_argument("--target-host", help="Override discovered ingress IP/hostname")
    parser.add_argument("--http-port", type=int, default=80)
    parser.add_argument("--https-port", type=int, default=443)
    parser.add_argument(
        "--refresh-interval",
        type=int,
        default=15,
        help="Seconds between ingress refreshes. Use 0 to disable refresh.",
    )
    parser.add_argument(
        "--host",
        action="append",
        dest="extra_hosts",
        default=[],
        help="Additional ingress hostname to allow. Can be passed more than once.",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_args()
    try:
        hosts, discovered_target = ingress_data(args.namespace)
    except (RuntimeError, json.JSONDecodeError) as exc:
        raise SystemExit(str(exc))

    hosts.update(args.extra_hosts)
    target_host = args.target_host or discovered_target

    if not hosts:
        raise SystemExit(f"No ingress hosts found in namespace {args.namespace}")
    if not target_host:
        raise SystemExit(
            "No ingress address found. Pass --target-host, or port-forward Traefik "
            "and use --target-host 127.0.0.1 --http-port 8080."
        )

    IngressProxyHandler.update_state(hosts, target_host)
    IngressProxyHandler.http_port = args.http_port
    IngressProxyHandler.https_port = args.https_port

    server = ThreadingHTTPServer(
        (args.listen_host, args.listen_port),
        IngressProxyHandler,
    )
    print(f"Proxy listening on http://{args.listen_host}:{args.listen_port}", flush=True)
    print(f"Forwarding {len(hosts)} ingress host(s) to {target_host}", flush=True)
    for host in sorted(hosts):
        print(f"  {host}", flush=True)
    if args.refresh_interval > 0:
        threading.Thread(
            target=refresh_ingress_state,
            args=(
                args.namespace,
                args.extra_hosts,
                args.target_host,
                args.refresh_interval,
            ),
            daemon=True,
        ).start()
        print(
            f"Refreshing ingress hosts every {args.refresh_interval} seconds",
            flush=True,
        )
    server.serve_forever()


if __name__ == "__main__":
    main()
