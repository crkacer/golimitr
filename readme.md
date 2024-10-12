# Rate Limiting Middleware in Go

A lightweight, flexible rate-limiting middleware for Go that implements both **Token Bucket** and **Leaky Bucket** algorithms. The middleware can be used in HTTP applications to control the rate of incoming requests, helping prevent abuse or overuse of resources.

## Features

- Supports **Token Bucket** and **Leaky Bucket** rate-limiting strategies.
- Easy-to-use middleware for any Go HTTP server.
- Customizable rate limits based on request strategies.
- Useful for microservices, API gateways, and high-traffic applications.

## Rate Limiting Strategies
- **Token Bucket**: Limits the number of requests allowed in a burst. Requests are refilled based on the specified time interval.

- **Leaky Bucket**: Allows requests at a steady rate over time. Excess requests will "leak" over time and requests beyond the bucket’s capacity are rejected.

## Installation

First, make sure you have Go installed on your system (Go 1.17+ recommended).

You can install the middleware package using `go get`:

```bash
go get github.com/crkacer/golimitr
```

Then import the package into your project:

```go
import "github.com/crkacer/golimitr"
```

## Usage
You can use the middleware in your Go HTTP server by specifying the rate-limiting strategy, request limit, and the rate interval. Here's a quick example of how to apply both Token Bucket and Leaky Bucket strategies.

## Example: Using Token Bucket Rate Limiting

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/crkacer/golimitr"
)

func main() {
    // Apply token bucket rate limiting middleware
    middleware, err := RateLimitingMiddleware("token_bucket", 5, time.Second)
    if err != nil {
        fmt.Println("Error creating middleware:", err)
        return
    }

    // Sample HTTP handler
    http.Handle("/", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Request successful!"))
    })))

    // Start the server
    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", nil)
}
```
## Example: Using Leaky Bucket Rate Limiting

```go
package main

import (
    "fmt"
    "net/http"
    "time"
    "github.com/crkacer/golimitr"
)

func main() {
    // Apply leaky bucket rate limiting middleware
    middleware, err := RateLimitingMiddleware("leaky_bucket", 3, time.Millisecond*500)
    if err != nil {
        fmt.Println("Error creating middleware:", err)
        return
    }

    // Sample HTTP handler
    http.Handle("/", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Request successful!"))
    })))

    // Start the server
    fmt.Println("Starting server on :8080")
    http.ListenAndServe(":8080", nil)
}
```

## Parameters
- **strategy**: Either "token_bucket" or "leaky_bucket". Defines the rate-limiting strategy to use.
- **limit**: The maximum number of requests allowed in the given time period (interval for Token Bucket, capacity for Leaky Bucket).
- **intervalOrRate**: The time interval for Token Bucket or the leak rate for Leaky Bucket.