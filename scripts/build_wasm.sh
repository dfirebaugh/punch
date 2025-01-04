#!/bin/bash

cp $(go env GOROOT)/misc/wasm/wasm_exec.js ./tools/ast_explorer/static/
GOOS=js GOARCH=wasm go build -o ./tools/ast_explorer/static/main.wasm ./tools/punchgen/ 
