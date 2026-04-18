# Performant TCP Server

## Overview
Ever need to build a TCP service that can handle a bunch of connections without falling over? This project gives you a solid foundation for that. It's a Go-based TCP server designed to keep things humming smoothly by managing incoming requests efficiently and preventing overload. You get a reliable way to accept connections and process them without your service getting bogged down.

## Features
-   **Concurrent Request Processing**: Handles incoming TCP connections using a worker pool, allowing the server to process multiple requests at the same time efficiently.
-   **Configurable Rate Limiting**: Integrates a token bucket algorithm to control the rate of incoming requests, protecting the server from being overwhelmed by traffic spikes.
-   **Prometheus Metrics Integration**: Exposes key operational metrics (like processed requests and rate-limited requests) via a Prometheus endpoint, making it easy to monitor the server's performance.
-   **Graceful Shutdown**: Ensures all active connections are handled and resources are cleaned up properly when the server is stopped, preventing data loss or abrupt service interruptions.

## Getting Started

### Installation
To get this server up and running on your local machine, follow these steps:

1.  **Clone the Repository**:
    ```bash
    git clone https://github.com/SaranHiruthikM/performant-tcp.git
    cd performant-tcp
    ```

2.  **Install Dependencies**:
    Go modules will handle the dependencies automatically when you build or run the project.
    ```bash
    go mod tidy
    ```

3.  **Build the Executable (Optional)**:
    If you want a standalone executable:
    ```bash
    go build -o server cmd/main.go
    ```

### Environment Variables
The server's behavior is configured using environment variables. You'll need to set these before running the application.

| Variable             | Example Value | Description                                          |
| :------------------- | :------------ | :--------------------------------------------------- |
| `SERVER_PORT`        | `3000`        | The port on which the TCP server will listen.        |
| `SERVER_WORKERS`     | `5`           | The number of worker goroutines to process requests. |
| `SERVER_QUEUE_SIZE`  | `100`         | The size of the job queue for incoming connections.  |
| `SERVER_TOKEN_RATE`  | `10`          | The rate (tokens per second) at which new tokens are added to the rate limiter. |
| `SERVER_TOKEN_LIMIT` | `20`          | The maximum number of tokens the rate limiter bucket can hold. |
| `METRICS_PATH`       | `/metrics`    | The HTTP path where Prometheus metrics are exposed.  |
| `METRICS_PORT`       | `9090`        | The port for the Prometheus metrics server.          |

Example `.env` file content:
```
SERVER_PORT=8080
SERVER_WORKERS=10
SERVER_QUEUE_SIZE=200
SERVER_TOKEN_RATE=50
SERVER_TOKEN_LIMIT=100
METRICS_PATH=/metrics
METRICS_PORT=9090
```

## Usage
Once you've set up your environment variables, running the server is straightforward.

1.  **Run the Server**:
    If you built an executable:
    ```bash
    ./server
    ```
    Or, run directly from source:
    ```bash
    go run cmd/main.go
    ```
    You should see output indicating the server has started, for example: `server running on :8080`.

2.  **Test the TCP Server**:
    You can connect to the server using `netcat` or `telnet`. Open another terminal and try:
    ```bash
    nc localhost 8080
    ```
    Type something and press Enter, then press Ctrl+D (or Ctrl+C to close `netcat`). The server will respond with `HTTP/1.1 200 OK\r\n\r\nHello\n`.

    To observe the rate limiting in action, send many requests quickly. You might see `HTTP/1.1 429 Too Many Requests\r\n\r\nRate limit exceeded\n` if you go over the configured limit.

3.  **Monitor Metrics**:
    While the server is running, open your web browser and navigate to `http://localhost:9090/metrics` (or whatever `METRICS_PORT` and `METRICS_PATH` you configured). You'll see the Prometheus metrics, including `total_requests_processed` and `total_requests_rate_limited`.

## Technologies Used

| Technology | Description                                        | Link                                                                      |
| :--------- | :------------------------------------------------- | :------------------------------------------------------------------------ |
| Go         | The primary language for building the server.      | [https://golang.org/](https://golang.org/)                                |
| Prometheus | Used for collecting and exposing application metrics. | [https://prometheus.io/](https://prometheus.io/)                           |

## Contributing
We'd love for you to contribute to this project! If you have suggestions for improvements, feature requests, or bug reports, please open an issue on the GitHub repository.

If you'd like to contribute code, here's a general guideline:
1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Make your changes and ensure your code follows Go best practices.
4.  Write clear, concise commit messages.
5.  Push your branch and open a pull request.
