package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dfirebaugh/punch/internal/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! Welcome to the Punch🥊 programming language!\n",
		user.Username)
	fmt.Printf("Feel free to type in commands\n")
	r := repl.New(os.Stdin, os.Stdout)
	r.Start()
}
