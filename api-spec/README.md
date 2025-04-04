# API spec

Schema of the external API types that are served by the Identity apps.

## Generate gRPC code

To generate the gRPC code from the API spec, run the following command from the root of the repository:

```bash
cd scripts && ./buf-generate.sh
```

This will generate the Golang code in the `identity/internal/pkg/generated` directory and the Python code in the `sdk/python/internal/generated` directory.

## OpenAPI Client and Proto Documentation

The API spec is also used to generate the OpenAPI client and the Proto documentation.
To generate the OpenAPI client and the Proto documentation, run the following command from the root of the repository:

```bash
make start_docs
```

You can then access the docs portal at `http://localhost:3000`.
