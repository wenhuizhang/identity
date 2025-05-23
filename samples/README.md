# Identity Samples

[![Lint](https://github.com/cisco-eti/pyramid/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/marketplace/actions/super-linter)
[![Contributor-Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-fbab2c.svg)](CODE_OF_CONDUCT.md)

<p align="center">
  <a href="https://agntcy.org">
    <picture>
      <source media="(prefers-color-scheme: dark)" srcset="../img/_logo-Agntcy_White@2x.png" width="300">
      <img alt="" src="../img/_logo-Agntcy_FullColor@2x.png" width="300">
    </picture>
  </a>
  <br />
  <caption>Welcome to the <b>Identity</b> repository</caption>
</p>

---

**Explore comprehensive guides and best practices for implementing and managing identity for agents.**

## Samples

You can find sample code and examples in the `samples` directory.
These samples are designed to help you understand how to use the Identity platform effectively.

You can create and verify the following type of badges:

- A2A Agent Badge based on the [A2A Agent Example](./samples/agent/a2a), based on the [A2A Agent Example](https://github.com/google/A2A/blob/main/samples/python/agents/langgraph)
- OASF Agent Badge based on the [OASF Agent Definition Example](./samples/agent/oasf), based on the [OASF Agent Example](https://hub.agntcy.org/)
- MCP Server Badge based on the [MCP Server Example](./samples/mcp), based on the [MCP Server Example](https://github.com/google/A2A/blob/main/samples/python/agents/langgraph)

## Prerequisites

To run the samples, you need to have the following prerequisites installed:

- [Docker](https://docs.docker.com/engine/install/)
- Ollama
- Python 3.12 or later

To install Ollama:

- First [Install Ollama](https://ollama.com/download)
- Then run the following command to install and run the Llama 3.2 model:

```bash
ollama run llama3.2
```

## Running the Samples

To run the samples, follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/agntcy/identity.git
```

2. Navigate to the `samples` directory and run the following command
   to start the Docker containers:

```bash
cd samples && docker-compose up -d
```

## Testing the Samples

### A2A Agent

To test the A2A Agent sample, run the following command:

```bash
cd samples/agent/a2a && python test_client.py
```

### MCP Server

To test the MCP Server sample, run the following command:

```bash
cd samples/mcp && python test_client.py
```
