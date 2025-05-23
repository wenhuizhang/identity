# Identity

[![Lint](https://github.com/cisco-eti/pyramid/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/marketplace/actions/super-linter)
[![Contributor-Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-fbab2c.svg)](CODE_OF_CONDUCT.md)

<p align="center">
  <a href="https://agntcy.org">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="img/_logo-Agntcy_White@2x.png" width="300">
      <img alt="" src="img/_logo-Agntcy_FullColor@2x.png" width="300">
    </picture>
  </a>
  <br />
  <caption>Welcome to the <b>Identity</b> repository</caption>
</p>

---

**Explore comprehensive guides and best practices for implementing and managing identity for agents.**

## Getting Started

- 🌐 Explore our full [Documentation](https://spec.identity.agntcy.org) to understand our platform's capabilities
- 📚 Dive into our [API Specs](https://spec.identity.agntcy.org/protodocs/agntcy/identity/core/v1alpha1/id.proto) for detailed API documentation
- 📦 Install our sample agents and MCP servers [Samples](samples/README.md)
- 🔐 Start issuing and verifying credentials with our [Quick Start](https://docs.agntcy.org/pages/identity-howto.html#quick-start)

## Development

### Prerequisites

To run the `Node` backend and the `Issuer` client locally, you need to have the following installed:

- [Docker](https://docs.docker.com/get-docker/)
- [Golang](https://go.dev/doc/install) 1.24 or later
- [Make](https://www.gnu.org/software/make/)

### API Specs

The [API Specs](api-spec/README.md) can be found in the `api-spec` directory.

To generate the API specs, run the following command:

```bash
make generate_proto
```

> [!NOTE]
> This will generate the `Protobuf` definitions, the `OpenAPI` specs and the `gRPC` stubs for the `Node` backend and the `Issuer` client.

> [!NOTE]
> The `Proto` definitions are generated in the `api-spec/proto` directory.
> The `Proto Messages and Enums` are generated from the `Go` types from the `core` package.
> The `Protobuf Services` are generated from the `api-spec/proto` directory.

### Starting the Node backend with Docker

To start the `Node` in development mode with `Docker`, run the following command:

```bash
make start_node dev=true
```

To stop the `Node,` run:

```bash
make stop_node
```

> [!NOTE]
> The `dev=true` flag is used to build the docker containers from the source code.
> This is useful for development purposes. If you want to use the pre-built images, you can omit this flag.

> [!NOTE]
> This will deploy a local persistent `Postgres` database and a local `Node` backend.
> The `Postgres` database will be available at `0.0.0.0:5432`and the `Node` will be available at`0.0.0.0:8080`.

### Building and Running the Issuer client locally

## Roadmap

See the [open issues](https://github.com/agntcy/identity/issues) for a list
of proposed features (and known issues).

## Contributing

Contributions are what make the open source community such an amazing place to
learn, inspire, and create. Any contributions you make are **greatly
appreciated**. For detailed contributing guidelines, please see
[CONTRIBUTING.md](CONTRIBUTING.md).

## Copyright Notice

[Copyright Notice and License](LICENSE)

Distributed under Apache 2.0 License. See LICENSE for more information.
Copyright [AGNTCY](https://github.com/agntcy) Contributors.
