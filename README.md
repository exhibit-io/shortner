# Redirector Package

## Overview

The `redirector` package provides a simple URL shortening and redirection service using Redis for storage. This package allows you to generate short URLs that redirect to longer, original URLs. It also tracks the number of visits to each shortened URL.

## Features

- **Create Short URLs:** Generate a short URL that maps to a longer URL.
- **Redirect Short URLs:** Redirect requests from the short URL to the original long URL.
- **Visit Tracking:** Count the number of visits for each short URL.
- **List All Short URLs:** Retrieve a list of all short URLs and their corresponding original URLs.

## Dependencies

- **Go Redis Client:** For interacting with Redis.
- **HTTPRouter:** For handling HTTP requests and routing.
- **Custom Config Package:** A configuration package for managing Redis and service settings.

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/exhibit-io/redirector.git
    cd redirector
    ```

2. Install dependencies:

    ```sh
    go get -u github.com/go-redis/redis/v8
    go get -u github.com/julienschmidt/httprouter
    ```

3. Build the package:

    ```sh
    go build -o redirector .
    ```

## Usage

### Initialization

Before using the `redirector` package, you need to initialize it with the required configuration:

```go
import (
    "github.com/exhibit-io/redirector/config"
    "github.com/exhibit-io/redirector"
)

func main() {
    config := config.LoadConfig() // Load your configuration
    redirector.Init(config)
}
```

### Create a Short URL

To create a new short URL, send a POST request to the `/create` endpoint with the JSON body containing the `url` and optional `expiresIn` fields:

```json
{
    "url": "https://example.com",
    "expiresIn": 3600
}
```

The service will respond with a JSON object containing the generated short URL.

### Redirect a Short URL

To redirect to the original URL, access the short URL in your browser. The service will automatically redirect you to the original URL.

### List All Short URLs

To get a list of all short URLs and their corresponding original URLs, send a GET request to the `/urls` endpoint. The service will respond with a JSON object containing all URL mappings.

## Example

```go
package main

import (
    "log"
    "net/http"

    "github.com/exhibit-io/redirector"
    "github.com/julienschmidt/httprouter"
)

func main() {
    config := config.LoadConfig()
    redirector.Init(config)

    router := httprouter.New()
    router.POST("/create", redirector.CreateRedirectURL)
    router.GET("/urls", redirector.GetAllRedirectURLs)
    router.GET("/:url", redirector.HandleURLRedirection)

    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## Configuration

The package expects a configuration structure for Redis and the redirector service. The configuration should provide the Redis address, password, and URI for the redirector service.

Example configuration structure:

```go
type Config struct {
    Redis struct {
        Addr     string
        Password string
    }
    Redirector struct {
        URI string
    }
}
```

## License

This package is licensed under the MIT License.