// Copyright 2024 uhppoted@twyst.co.za. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

/*
Package uhppote/codegen generates the integration test vectors for the project
submodules.
*/
package main

import (
	"fmt"
	"log/slog"
)

const VERSION = "v0.0.0"

func main() {
	fmt.Println()
	fmt.Printf("  ** integration tests: codegen %v\n", VERSION)
	fmt.Println()

	generate()
}

func infof(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-16v %v", tag, format)
	msg := fmt.Sprintf(f, args...)

	slog.Info(msg)
}

func errorf(tag string, format string, args ...any) {
	f := fmt.Sprintf("%-16v %v", tag, format)
	msg := fmt.Sprintf(f, args...)

	slog.Error(msg)
}
