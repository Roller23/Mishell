package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

const refString = "ref: refs/heads/"
const headPath = ".git/HEAD"

func getCurrDir() (string, *string) {
	dir, _ := os.Getwd()
	var branchName *string = nil
	if fileExists(headPath) {
		bytes, err := ioutil.ReadFile(headPath)
		if err == nil {
			firstLine := strings.Split(string(bytes), "\n")[0]
			idx := strings.Index(firstLine, refString)
			if idx != -1 {
				branch := firstLine[idx+len(refString):]
				branchName = &branch
			}
		}
	}
	return filepath.Base(dir), branchName
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
		currDir, branch := getCurrDir()
		fmt.Printf("%s %s", color.YellowString(">"), color.HiCyanString(currDir))
		if branch != nil {
			fmt.Printf("%s%s%s", color.BlueString(":("), color.HiRedString(*branch), color.BlueString(")"))
		}
		fmt.Print(" ")
		input, err := reader.ReadString('\n')
		if err != nil {
			abort(err.Error())
		}
		input = strings.TrimSuffix(input, "\n")
		handleInput(input)
	}
}
