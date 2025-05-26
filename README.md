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
- 🛠️ Use the [Issuer CLI](identity/cmd/issuer/README.md) to manage issuers and credentials
- 🏗️ Deploy the [Node Backend](identity/cmd/node/README.md) to handle identity management
- 🔐 Start issuing and verifying credentials with our [Quick Start](https://docs.agntcy.org/pages/identity-howto.html#quick-start)

## Quick Start

### Prerequisites

To run the `Node Backend` the `Issuer CLI`, and the `Samples` locally, you need to have the following installed:

- [Docker](https://docs.docker.com/get-docker/)
- [Golang](https://go.dev/doc/install) 1.24 or later
- [Make](https://www.gnu.org/software/make/)
- [Ollama](https://ollama.com/download)
- Python 3.12 or later

### Step 1: Install the Issuer CLI

```bash
go install github.com/agntcy/identity/cmd/issuer@latest
```

### Step 2: Clone the Repository

```bash
git clone https://github.com/agntcy/identity.git
```

### Step 3: Start the Node Backend with Docker

Run the following command from the root of the repository:

```bash
make start_node
```

### Step 4: Run the Samples with Docker

1. Run the Llama 3.2 model:

   ```bash
   ollama run llama3.2
   ```

2. Navigate to the `samples` directory and run the following command:

   ```bash
   docker-compose up -d
   ```

### Step 5: Register as an Issuer

### Step 6: Onboard an MCP Server

### Step 7: Verify Credentials

## Development

For more detailed development instructions please refer to the following sections:

- [Node Backend](identity/cmd/node/README.md)
- [Issuer CLI](identity/cmd/issuer/README.md)
- [Samples](samples/README.md)
- [Api Specs](api-specs/README.md)

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

```

```
