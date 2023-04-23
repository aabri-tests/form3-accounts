# Amine Abri - [Form3 Accounts Client Library]
This is a Go client library for interacting with the Form3 Accounts Fake API, developed by **Amine Abri**. 
I have two years of experience with Go, so I am familiar with the language, but I'm still learning and refining my understanding of best practices.

## Features
- Easy-to-use and intuitive API
- Built-in error handling and retries
- Customizable retry policy
- Comprehensive API coverage with clear methods and data structures
- Support for JSON serialization
- Logging capabilities

## Requirements

- Go 1.20 or later

## Installation

You can install **form3-accounts** using `go get`:

```sh
go get -u github.com/aabri-assignments/form3-accounts/v1
```

## Getting Started
To get started with **form3-accounts**, you will need to import the library and configure it with the appropriate options:

```go
package main

import (
  "fmt"
  "time"

  "github.com/aabri-assignments/form3-accounts/v1"
  "github.com/aabri-assignments/form3-accounts/v1/accounts"
  "github.com/aabri-assignments/form3-accounts/v1/accounts/models"
  "github.com/aabri-assignments/form3-accounts/v1/pkg/logging"
  "github.com/google/uuid"
)

func main() {
  opts := accounts.Options{
    BaseURL:      "http://localhost:8080",
    Duration:     3 * time.Minute,
    Retries:      3,
    InitialDelay: 300 * time.Millisecond,
    Multiplier:   2,
    Factor:       0.1,
    LogLevel:     logging.LevelInfo,
  }
  // Create a new Form3 client with the base URL of the API.
  form3Client, err := form3.NewForm3(opts)
  if err != nil {
    panic(err)
  }

  // Start making requests to the Example API
}
```

Now you can start making requests to the Form3 Accounts Fake API. For more detailed usage and [examples](./examples), please check the source code, tests, and the examples folder in the GitHub repository.

## Examples
The examples folder contains sample programs that demonstrate different use cases for the Form3 Accounts Client Library. You can find examples for creating, fetching and deleting accounts.

- [Create Account Example](./examples/create_account)
- [Fetch Account Example](./examples/fetch_account)
- [Delete Account Example](./examples/delete_account)


To run an example, navigate to the example's folder and run the `main.go` file:

```shell
cd examples/create_account
go run main.go
```

Make sure to replace any placeholder values (such as API endpoints or accounts uuid) with your own information before running the examples.

## End-to-End Tests

This library includes end-to-end tests that verify the functionality of the Form3 Accounts Fake API. To run the tests, you can use the provided Docker Compose file. Simply navigate to the root of the repository and run the following command:

```shell
docker-compose up --build
```
This command will build the necessary Docker images and run the end-to-end tests automatically. The tests will use a local instance of the Form3 Accounts Fake API, so make sure that the API is running before running the tests.
