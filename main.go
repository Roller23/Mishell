package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func getCurrDir() string {
	dir, _ := os.Getwd()
	return filepath.Base(dir)
}

var builtInCommands = map[string]func(string){
	"cd": func(input string) {
		args := strings.Split(input, " ")
		if len(args) < 2 {
			return
		}
		err := os.Chdir(args[1])
		if err != nil {
			log.Println("Error:", err)
			return
		}
	},
	"exit": func(input string) {
		os.Exit(0)
	},
	"help": func(input string) {
		fmt.Println("WIP")
	},
}

func abort(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func handleInput(input string) {
	args := strings.Split(input, " ")
	if fn, ok := builtInCommands[args[0]]; ok {
		fn(input)
		return
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		log.Println("Error:", err)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Welcome to Mishell 0.1\n\n")

	for {
		fmt.Printf("%s %s ", color.YellowString(">"), color.CyanString(getCurrDir()))
		input, err := reader.ReadString('\n')
		if err != nil {
			abort(err.Error())
		}
		input = strings.TrimSuffix(input, "\n")
		handleInput(input)
	}
}
