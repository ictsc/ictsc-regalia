//go:build tools
package main

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/sqldef/sqldef/cmd/psqldef"
)
