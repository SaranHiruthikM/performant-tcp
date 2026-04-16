# Performant TCP Server

## Overview

This project provides a really efficient way to handle incoming network requests. It's built to keep your service responsive and stable, even when you're getting a lot of traffic, by distributing work and preventing overload before it becomes a problem. Essentially, it helps you manage high-volume network interactions smoothly.

## Features

- **Efficient Worker Pool**: Distributes incoming TCP connections across a pool of workers, ensuring requests are processed concurrently without blocking the main listener.
- **Token Bucket Rate Limiting**: Implements a robust token bucket algorithm to control the rate of incoming requests, protecting your service from being overwhelmed.
- **Prometheus Metrics Integration**: Exposes custom metrics for total requests processed and total requests rate-limited, allowing for easy monitoring and observability.
- **Simple TCP Listener**: Sets up a basic TCP server that can accept and respond to client connections.

## Getting Started

First things first, you'll need Go installed on your machine. This project was developed with Go 1.25.8, so make sure you're using a compatible version.

### Installation

1.  **Clone the Repository**:

    ```bash
    git clone https://github.com/SaranHiruthikM/performant-tcp.git
    cd performant-tcp
    ```

2.  **Download Dependencies**:
    ```bash
    go mod download
    ```

### Environment Variables

Before running the server, you'll need to set up some environment variables. These control things like the server port, worker pool size, and rate limiter settings.

| Variable             | Example Value | Description                                                                        |
| :------------------- | :------------ | :--------------------------------------------------------------------------------- |
| `SERVER_PORT`        | `8080`        | The port the TCP server will listen on.                                            |
| `SERVER_WORKERS`     | `5`           | The number of worker goroutines in the pool to process connections.                |
| `SERVER_QUEUE_SIZE`  | `100`         | The size of the job queue for the worker pool.                                     |
| `SERVER_TOKEN_RATE`  | `10`          | The rate at which tokens are added to the rate limiter bucket (tokens per second). |
| `SERVER_TOKEN_LIMIT` | `20`          | The maximum number of tokens the rate limiter bucket can hold.                     |
| `METRICS_PATH`       | `/metrics`    | The HTTP path where Prometheus metrics will be exposed.                            |
| `METRICS_PORT`       | `9090`        | The port where the Prometheus metrics server will listen.                          |

**Example `.env` file content:**

```
SERVER_PORT=8080
SERVER_WORKERS=5
SERVER_QUEUE_SIZE=100
SERVER_TOKEN_RATE=10
SERVER_TOKEN_LIMIT=20
METRICS_PATH=/metrics
METRICS_PORT=9090
```

## Usage

Once you've set your environment variables, you can run the server:

1.  **Run the application**:

    ```bash
    go run cmd/main.go
    ```

    You'll see a log message confirming the server has started, like `Server started at :8080`.

2.  **Send a request to the TCP server**:
    You can use `netcat` or `curl` to interact with the server.

    **Using `netcat`**:

    ```bash
    echo "Hello Server" | nc localhost 8080
    ```

    The server will respond with:

    ```
    HTTP/1.1 200 OK

    Hello
    ```

    **Using `curl`**:

    ```bash
    curl -v localhost:8080
    ```

    The server will respond similarly, indicating a successful connection and a "Hello" message.

3.  **Test the Rate Limiter**:
    If you send requests faster than the `SERVER_TOKEN_RATE` allows, some connections will be rate-limited.

    ```bash
    # Try sending many requests quickly
    for i in $(seq 1 30); do echo "Request $i" | nc -w 1 localhost 8080 & done
    ```

    You'll see some responses like:

    ```
    HTTP/1.1 429 Too Many Requests

    Rate limit exceeded
    ```

4.  **Access Prometheus Metrics**:
    Open your browser or use `curl` to view the exposed metrics:
    ```bash
    curl http://localhost:9090/metrics
    ```
    You'll see output similar to this, showing the counts of processed and rate-limited requests:
    ```
    # HELP total_requests_processed Number of requests successfully processed
    # TYPE total_requests_processed counter
    total_requests_processed 100
    # HELP total_requests_rate_limited Number of requests successfully rate limited
    # TYPE total_requests_rate_limited counter
    total_requests_rate_limited 5
    ```

## Technologies Used

| Technology                                                          | Description                                         |
| :------------------------------------------------------------------ | :-------------------------------------------------- |
| [Go](https://golang.org/)                                           | The primary language for building the server logic. |
| [Prometheus Client Go](https://github.com/prometheus/client_golang) | Go client library for Prometheus metrics.           |

## Contributing

Hey, if you've got ideas on how to make this project even better, I'd love to hear them! Feel free to open an issue to discuss features or bugs, or even better, submit a pull request with your changes. Just make sure your code adheres to standard Go formatting and includes clear commit messages.

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![Prometheus](https://img.shields.io/badge/Prometheus-E6522C?style=for-the-badge&logo=prometheus&logoColor=white)](https://prometheus.io/)
