# Identity CLI Tool

The Identity CLI tool provides identity management capabilities for the AGNTCY Internet of Agents.
It allows you to create and manage cryptographic keys, register as an issuer, generate metadata for identities, issue badges for Agent and MCP Server identities, and verify badges for existing identities.

## Installation

### Build from Source

To build the Identity CLI from source:

```bash
cd identity
go build -o identity cmd/issuer/main.go
```

Then move the binary to a location in your PATH (optional):

```bash
mv identity /usr/local/bin/
```

This will allow you to run the CLI from anywhere in your terminal with the command `identity`.

### Run without Building

You can also run the CLI directly without building and installing it.
This is useful for testing or development purposes:

```bash
cd identity
go run cmd/issuer/main.go
```

## Usage

The Identity CLI follows a hierarchical command structure:

```
identity [command] [subcommand] [flags]
```

### Core Commands

- **vault**: Manage cryptographic vaults and keys
- **issuer**: Register and manage issuer configurations
- **metadata**: Generate and manage metadata for identities
- **badge**: Issue and publish badges for identities
- **verify**: Verify identity badges
- **config**: Display the current configuration context

### Common Workflows

#### Setting Up a New Identity Environment

1. **Create a vault** to store cryptographic keys:

   ```bash
   identity vault create file -f /path/to/keys.json -n "My Vault"
   ```

2. **Register as an issuer**:

   ```bash
   identity issuer register -o "My Organization" -c "client-id" -s "client-secret" -u "https://idp.example.com"
   ```

3. **Generate metadata**:

   ```bash
   identity metadata generate -i "client-id" -s "client-secret" -u "https://idp.example.com"
   ```

4. **Issue a badge**:

    ```bash
    identity badge issue [type] [options]
    ```

    #### Badge Types and Examples

    Badges can be issued for different types of content:

    - **OASF Files** ([OASF Schema](https://schema.oasf.agntcy.org/objects/agent)):
        ```bash
        identity badge issue oasf -f /path/to/oasf_content.json
        ```

    - **A2A Agent Cards** ([A2A Documentation](https://google.github.io/A2A/tutorials/python/3-agent-skills-and-card/#agent-card)):
        ```bash
        identity badge issue a2a -u http://localhost:9091/.well-known/agent.json
        ```

    - **MCP Servers** ([MCP Specification](https://github.com/modelcontextprotocol/servers)):
        ```bash
        identity badge issue mcp -u http://localhost:9090
        ```

    - **Generic Files**:
        ```bash
        identity badge issue file -f /path/to/badge_content.json
        ```

5. **Publish a badge**:

   ```bash
   identity badge publish
   ```

#### Managing Existing Components

- **List existing vaults**:
  ```bash
  identity vault list
  ```

- **Show details of an issuer**:
  ```bash
  identity issuer show -i [issuer-id]
  ```

- **Load a different metadata configuration**:
  ```bash
  identity metadata load -m [metadata-id]
  ```

- **View current configuration**:
  ```bash
  identity config
  ```

- **Verify a badge from a file**:
  ```bash
  identity verify -f /path/to/badge.json
  ```

## Documentation

For more detailed documentation on each command:

```bash
identity [command] --help
```

For a full command overview:

```bash
identity --help
```
