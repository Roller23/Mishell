package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/peterh/liner"
)

const refString = "ref: refs/heads/"
const headPath = ".git/HEAD"

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func hasDotGit(dirPath string) (bool, string) {
	hPath := path.Join(dirPath, headPath)
	return fileExists(hPath), hPath
}

func insideRepo(fPath string) (bool, string) {
	changed := true
	oldPath := fPath
	for changed {
		if ok, hPath := hasDotGit(fPath); ok {
			return true, hPath
		}
		fPath = path.Dir(fPath)
		changed = fPath != oldPath
		oldPath = fPath
	}
	return false, ""
}

func getCurrDir() (string, *string) {
	dir, _ := os.Getwd()
	var branchName *string = nil
	if inside, hPath := insideRepo(dir); inside {
		bytes, err := ioutil.ReadFile(hPath)
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
			fmt.Println("Error:", err)
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

func getBranchStr(name string) string {
	return fmt.Sprintf("%s%s%s", color.HiBlueString(":("), color.HiRedString(name), color.HiBlueString(")"))
}

func getDirStr(dir string) string {
	return fmt.Sprintf("%s %s", color.YellowString(">"), color.HiCyanString(dir))
}

func getInputLine() string {
	currDir, branch := getCurrDir()
	line := getDirStr(currDir)
	if branch != nil {
		line += getBranchStr(*branch)
	}
	line += " "
	return line
}

func getInput(line *liner.State) string {
	fmt.Print(getInputLine())
	input, err := line.Prompt("")
	if err != nil {
		abort(err.Error())
	}
	line.AppendHistory(input)
	return input
}

func main() {
	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)
	fmt.Print("Welcome to Mishell 1.0\n\n")

	for {
		input := getInput(line)
		if input == "" {
			continue
		}
		handleInput(input)
	}
}
