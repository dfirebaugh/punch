punch:
	go build ./cmd/punch/

run-ast-explorer:
	bash ./scripts/build_wasm.sh
	go run ./tools/ast_explorer/

clean:
	rm punch
