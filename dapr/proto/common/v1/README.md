# Common

This folder is intended for the Common protos amongst the packages in the `dapr/proto` folder.

## Proto client generation

Pre-requisites:
1. Install protoc version: [v4.25.4](https://github.com/protocolbuffers/protobuf/releases/tag/v4.25.4)

2. Install protoc-gen-go and protoc-gen-go-grpc

```bash
make init-proto
```

*If* protoc is already installed:

3. Generate gRPC proto clients from the root of the project

```bash
make gen-proto
```

4. See the auto-generated files in `pkg/proto`
