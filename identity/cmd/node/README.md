# Node Backend

The `Node Backend` is a core component of the Identity Platform, responsible for allowing organizations to register, issue and verify agent identities.
It provides a RESTful API for interaction with the Identity services.

## Prerequisites

To run the `Node Backend` locally, you need to have the following installed:

- [Docker](https://docs.docker.com/get-docker/)
- [Golang](https://go.dev/doc/install) 1.24 or later
- [Make](https://www.gnu.org/software/make/)

## Installation

## Starting the Node Backend with Docker

To start the `Node` locally with `Docker`, run the following command:

```bash
make start_node
```

To stop the `Node,` run:

```bash
make stop_node
```

> [!NOTE]
> This will deploy a local persistent `Postgres` database and a local `Node` backend.
> The `Postgres` database will be available at `0.0.0.0:5432`and the `Node` will be available at`0.0.0.0:8080`.

## Development

## Starting the Node Backend with Docker in Development Mode

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
>
> This will deploy a local persistent `Postgres` database and a local `Node` backend.
> The `Postgres` database will be available at `0.0.0.0:5432`and the `Node` will be available at`0.0.0.0:8080`.
