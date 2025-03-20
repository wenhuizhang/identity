# Assets

This directory contains static web assets that are served by the CLI tool in web mode.
The assets come from the `ui/dist` directory in the root of the project, and are copied here by the `go generate` command in the `cmd/issuer/main.go` file.

The `go generate` command copies the contents of the `ui/dist` directory to the `pkg/assets/web` directory.

In order to run the CLI, you must first make sure that the `ui/dist` directory is up to date by running the `yarn build` command in the `ui` directory, and then run the `go generate` command in the `cmd/issuer/main.go` file.
