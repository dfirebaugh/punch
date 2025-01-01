package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/dfirebaugh/punch/internal/lexer"
	"github.com/dfirebaugh/punch/internal/parser"
)

//go:embed static/*
var staticFiles embed.FS

func parseHandler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if rec := recover(); rec != nil {
			errMessage := fmt.Sprintf("An error occurred: %v", rec)
			http.Error(w, errMessage, http.StatusInternalServerError)
		}
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody struct {
		Source string `json:"source"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	l := lexer.New("example", requestBody.Source)
	p := parser.New(l)

	program := p.ParseProgram("ast_explorer")

	astJSON, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		http.Error(w, "Failed to generate AST JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(astJSON)
}

func main() {
	publicFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic(fmt.Errorf("failed to get subdirectory: %w", err))
	}
	fileServer := http.FileServer(http.FS(publicFS))

	http.Handle("/", fileServer)

	http.HandleFunc("/parse", parseHandler)

	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
