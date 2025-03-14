# Issuer entry point

The Issuer CLI tool contains a web mode that serves a web interface for taking the same actions that are available in the command line interface.

The web assets come from the `ui/dist` directory in the root of the project, and are copied here by the `go generate` command in the `cmd/issuer/main.go` file.

The `go generate` command copies the contents of the `ui/dist` directory to the `pkg/assets/web` directory.

In order to run the CLI, you must first make sure that the `ui/dist` directory is up to date by running the `yarn build` command in the `ui` directory, and then run the `go generate` command in the `cmd/issuer/main.go` file.

Then, you can run the following commands from the root of the project to start the CLI in web mode:

```sh
cd pyramid 
go generate cmd/issuer/main.go
go run cmd/issuer/main.go web [port]
```

After copying the files to the `pkg/assets/web` directory, with the go generate command, you are then able to build the CLI, with the following command:

```sh
go build -o pyramid cmd/issuer/main.go
```

Finally, you can run the built CLI with the following command:

```sh
./pyramid
```