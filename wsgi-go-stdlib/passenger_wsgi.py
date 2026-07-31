import atexit
import os
import socket
import subprocess
import threading

import requests

GO_BINARY = os.path.join(os.path.dirname(__file__), "wsgi-go-stdlib")
GO_PORT = int(os.environ.get("GO_BACKEND_PORT", "18080"))
GO_HOST = "127.0.0.1"

_process = None
_lock = threading.Lock()


def _port_open():
    try:
        with socket.create_connection((GO_HOST, GO_PORT), timeout=0.2):
            return True
    except OSError:
        return False


def _ensure_go_running():
    global _process

    if _port_open():
        return

    with _lock:
        if _port_open():
            return

        _process = subprocess.Popen(
            [GO_BINARY],
            env={**os.environ, "PORT": str(GO_PORT)},
            stdout=subprocess.DEVNULL,
            stderr=subprocess.STDOUT,
        )


@atexit.register
def _cleanup():
    global _process
    if _process:
        _process.terminate()


def application(environ, start_response):
    _ensure_go_running()

    path_info = environ.get("PATH_INFO", "")
    url = f"http://{GO_HOST}:{GO_PORT}{path_info}"

    query_string = environ.get("QUERY_STRING")
    if query_string:
        url += f"?{query_string}"

    content_length = int(environ.get("CONTENT_LENGTH") or 0)
    body = environ["wsgi.input"].read(content_length)

    headers = {}
    for key, value in environ.items():
        if key.startswith("HTTP_"):
            headers[key[5:].replace("_", "-")] = value

    if environ.get("CONTENT_TYPE"):
        headers["Content-Type"] = environ["CONTENT_TYPE"]

    response = requests.request(
        method=environ.get("REQUEST_METHOD", "GET"),
        url=url,
        headers=headers,
        data=body,
        allow_redirects=False,
        timeout=15,
    )

    response_headers = [
        (key, value)
        for key, value in response.headers.items()
        if key.lower() not in ("transfer-encoding", "connection", "content-encoding")
    ]

    start_response(f"{response.status_code} {response.reason}", response_headers)
    return [response.content]
