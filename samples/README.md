# Identity Samples

You can find sample code and examples in the `samples` directory.
These samples are designed to help you understand how to use the Identity platform effectively.
The samples are composed of a `Currency Exchange A2A Agent` leveraging a `Currency Exchange MCP Server` and an `OASF Agent Definition`.

You can create and verify the following type of badges:

- A2A Agent Badge based on the [`Currency Exchange A2A Agent`](agent/a2a)
- OASF Agent Badge based on a [OASF Agent Definition Example](agent/oasf)
- MCP Server Badge based on the [`Currency Exchange MCP Server`](mcp)

> [!NOTE]
> These samples are based on the following [A2A Agent Example](https://github.com/google/A2A/blob/main/samples/python/agents/langgraph).
> The `OASF Agent Definitions` could be found in the [Agent Directory](https://hub.agntcy.org/explore).

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
   docker compose up -d
   ```

## Testing the Samples

### A2A Agent

To test the A2A Agent sample, navigate to the `samples/agent/a2a/currency_exchange`
directory and run the following command:

```bash
python test_client.py
```

### MCP Server

To test the MCP Server sample, navigate to the `samples/mcp/currency_exchange`
directory and run the following command:

```bash
python test_client.py
```
