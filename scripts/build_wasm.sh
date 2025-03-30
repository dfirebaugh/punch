#!/bin/bash

cp $(go env GOROOT)/lib/wasm/wasm_exec.js ./tools/ast_explorer/static/
GOOS=js GOARCH=wasm go build -o ./tools/ast_explorer/static/main.wasm ./tools/punchgen/ 
