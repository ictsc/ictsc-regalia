//go:build tools
package main

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/sqldef/sqldef/cmd/psqldef"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	_ "connectrpc.com/connect/cmd/protoc-gen-connect-go"
)
