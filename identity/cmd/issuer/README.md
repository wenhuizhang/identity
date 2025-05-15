# Issuer entry point

The Issuer CLI tool contains a web mode that serves a web interface for taking the same actions that are available in the command line interface.

Then, you can run the following commands from the root of the project to start the CLI in web mode:

```sh
cd identity
go build -o identity cmd/issuer/main.go
```

Finally, you can run the built CLI with the following command:

```sh
./identity
```
