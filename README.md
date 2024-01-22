# Gatekeeper

Gatekeeper is a dynamic request handling and routing application designed to manage and enforce various access control rules and rate limits based on request attributes. It supports IP-based restrictions, rate limiting, and can be configured to suit various use cases.

## Overview

Gatekeeper operates using a concept known as "gates." These gates serve as customizable checkpoints within the application's request handling process. Each gate is designed to evaluate specific criteria or conditions associated with incoming requests. Depending on the outcome of these evaluations, Gatekeeper can take various actions such as allowing or denying access, imposing rate limits, or applying custom middleware.

These gates are highly dynamic and can be configured to suit your specific needs. They enable you to define and enforce access control policies, making it possible to secure your application effectively. Whether it's restricting access based on IP addresses, implementing rate limits, or adding JWT authentication middleware, Gatekeeper's gates offer the flexibility and control you need to manage your application's requests efficiently.

By configuring and customizing these gates, you can fine-tune how Gatekeeper handles incoming requests, ensuring the security and performance of your application.

## Features

- IP Restriction: Restrict access to your application based on client IP addresses.
- Rate Limiting: Limit the number of requests from a particular IP or for a specific endpoint within a defined time window.
- JWT Authentication Middleware: Secure your endpoints with JWT-based authentication.
- Flexible Configuration: Easily configure access rules through a YAML file.

## Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.x or later)

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/atanurdemir/gatekeeper.git
   cd gatekeeper
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

### Configuration

1. **Setting up the configuration file:**

Modify the `gates.yaml` file in the src/config directory to define your access control rules. Here's an example configuration:

```yaml
gates:
  - path: "/"
    method: "POST"
    rules:
      - name: "ip_restriction"
        config:
          allowed_ips: ["127.0.0.1", "192.168.1.1"]
      - name: "ip_rate"
        config:
          limit: 5
          duration: 30000
```

### Running the Application

1. **Running in development mode:**
   To run the application in development mode with hot-reloading, use the following command:

   ```bash
   make dev
   ```

2. **Building the application:**

To build the application, use the following command:

```bash
make build
```

The binary will be generated in the dist directory as gatekeeper.
To run the application after building, use:

```bash
make run
```

## Contributing

If you wish to contribute to this project, please feel free to fork the repository and submit a pull request.

## License

This library is licensed under MIT License. Please see the [License](LICENSE) file for more details.
