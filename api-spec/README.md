# API spec

Schema of the external API types that are served by the Identity apps.

## Generate gRPC code

To generate the gRPC code from the API spec, run the following command from the root of the repository:

```bash
cd scripts && ./buf-generate.sh
```

This will generate the Golang code in the `identity/api` directory.

**Note**: The Proto Documentation, the OpenAPI Client and the JSON Schema will be generated in the `api-spec/static` directory.
