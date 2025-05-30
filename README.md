# Identity

[![Lint](https://github.com/agntcy/identity/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/marketplace/actions/super-linter)
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
- 🛠️ Use the [Issuer CLI](cmd/issuer/README.md) to manage issuers and credentials
- 🏗️ Deploy the [Node Backend](cmd/node/README.md) to handle identity management
- 🔐 Start issuing and verifying credentials with our [Quick Start](https://docs.agntcy.org/pages/identity-howto.html#quick-start)

## Quick Start

### Prerequisites

To run the `Node Backend` the `Issuer CLI`, and the `Samples` locally, you need to have the following installed:

- [Docker](https://docs.docker.com/get-docker/)
- [Ollama](https://ollama.com/download)
- [Okta CLI](https://cli.okta.com/manual/#installation)

### Step 1: Install the Issuer CLI

Download the `Issuer CLI` binary corresponding to your platform from the [latest releases](https://github.com/agntcy/identity/releases).

> [!NOTE]
> On some platforms you might need to add execution permissions and/or approve the binary in `System Security Settings`.
>
> For easier use, consider moving the binary to your `$PATH` or to the `/usr/local/bin` folder.

If you have `Golang` set up locally, you could also use the `go install command`:

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
./deployments/scripts/identity/launch_node.sh
```

Or use `make` if available locally:

```bash
make start_node
```

### Step 4: Run the Samples with Docker

1. Run the Llama 3.2 model:

   ```bash
   ollama run llama3.2
   ```

2. Navigate to the `samples` directory and run the following command
   to deploy the `Currency Exchange A2A Agent` leveraging the `Currency Exchange MCP Server`:

   ```bash
   docker-compose up -d
   ```

3. [Optional] Test the samples using the provided [test clients](./samples/README.md#testing-the-samples).

### Step 5: Create a local Vault and generate keys

1. Create a local vault to store generated cryptographic keys:

   ```bash
   identity vault connect file -f ~/.identity/vault.json -v "My Vault"
   ```

2. Generate a new key pair and store it in the vault:

   ```bash
   identity vault key generate
   ```

### Step 6: Register as an Issuer

For this quick start we will use Okta as an IdP to create an application for the Issuer:

1. Run the following command from the root repository to create a new Okta application:

   ```bash
   . ./demo/scripts/create_okta_app
   ```

2. In the interactive prompt, choose the following options:

   `> 4: Service (Machine-to-Machine)`, `> 5: Other`

3. Register the Issuer using the `Issuer CLI` and the environment variables from the previous step:

   ```bash
   identity issuer register -o "My Organization" \
       -c "$OKTA_OAUTH2_CLIENT_ID" -s "$OKTA_OAUTH2_CLIENT_SECRET" -u "$OKTA_OAUTH2_ISSUER"
   ```

> [!NOTE]
> You can now access the `Issuer's Well-Known Public Key` at [`http://localhost:4000/v1alpha1/issuer/{common_name}/.well-known/jwks.json`](http://localhost:4000/v1alpha1/issuer/{common_name}/.well-known/jwks.json),
> where `{common_name}` is the common name you provided during registration.

### Step 7: Generate metadata for an MCP Server

Create a second application for the MCP Server metadata using the Okta, similar to the previous step:

1. Run the following command from the root repository to create a new Okta application:

   ```bash
   . ./demo/scripts/create_okta_app
   ```

2. In the interactive prompt, choose the following options:

   `> 4: Service (Machine-to-Machine)`, `> 5: Other`

3. Generate metadata for the MCP Server using the `Issuer CLI` and the environment variables from the previous step:

   ```bash
   identity metadata generate -c "$OKTA_OAUTH2_CLIENT_ID" \
       -s "$OKTA_OAUTH2_CLIENT_SECRET" -u "$OKTA_OAUTH2_ISSUER"
   ```

> [!NOTE]
> When successful, this command will print the metadata ID, which you will need in the next step to view published badges that are linked to this metadata.

### Step 8: Issue and Publish a Badge for the MCP Server

1. Issue a badge for the MCP Server:

   ```bash
   identity badge issue mcp -u http://localhost:9090 -n "My MCP Server"
   ```

2. Publish the badge:

   ```bash
   identity badge publish
   ```

> [!NOTE]
> You can now access the `VCs as a Well-Known` at [`http://localhost:4000/v1alpha1/vc/{metadata_id}/.well-known/vcs.json`](http://localhost:4000/v1alpha1/vc/{client_id}/.well-known/vcs.json),
> where `{metadata_id}` is the metadata ID you generated in the previous step.

### (Optional) Step 9: Verify a Published Badge

You can use the `Issuer CLI` to verify a published badge any published badge, not just those that you issued yourself.
This allows others to verify the Agent and MCP badges you publish.

1. Download the badge that you created in the previous step:

   ```bash
   # Download the published badges linked to the metadata, replace {metadata_id} with the actual metadata ID
   curl -o vcs.json http://localhost:4000/v1alpha1/vc/{metadata_id}/.well-known/vcs.json
   ```

2. Verify the badges using the `Issuer CLI`:

   ```bash
   identity verify -f vcs.json
   ```

## Development

For more detailed development instructions please refer to the following sections:

- [Node Backend](cmd/node/README.md)
- [Issuer CLI](cmd/issuer/README.md)
- [Samples](samples/README.md)
- [Api Spec](api-spec/README.md)

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
