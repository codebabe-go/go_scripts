package main

/**
 * git 自动提交脚本
 * 期间一旦遇到需要手动开那个解决错误, 就会丢出异常终止运行
 */

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
	"errors"
	"bytes"
	"bufio"
)

const NO_CONDITION string = "no condition, cannot match"
const DEFAULT_COMMENT string = "code.babe push code"
const UNMATCHED_CONDITION string = "cannot match condition"
const CONTINUE string = "continue"
const EMPTY string = ""

var ERROR_MAP map[string][]string = make(map[string][]string)

type Result struct {
	output string
	err error
}

// error map information init
func init()  {
	ERROR_MAP["checkout"] = []string{}
	ERROR_MAP["add"] = []string{}
	ERROR_MAP["commit"] = []string{}
	ERROR_MAP["pull"] = []string{}
	ERROR_MAP["push"] = []string{}
}

func NewResult(output string, errMsg string) *Result {
	result := new(Result)
	result.output = output
	if len(errMsg) != 0 {
		result.err = errors.New(errMsg)
	}
	return result
}

func GitPush(comment, checkout string) error {
	branch := checkoutBranch(checkout)
	fmt.Println(branch)
	checkoutInfo := do(NO_CONDITION, "git checkout " + branch)
	toContinue := checkError(checkoutInfo.err)
	goContinue(toContinue, func() {
		// commit
		fmt.Println("commit")
		os.Exit(0)
	})
	checkoutSuccess(checkoutInfo.output)
	gitAdd()
	gitCommit(comment)
	gitPull()
	gitPush(branch)
	return errors.New("")
}

func gitAdd() error {
	result := do(NO_CONDITION, "git add .")
	toContinue := checkError(result.err)
	goContinue(toContinue, func() {
		// commit
		fmt.Println("add")
		os.Exit(0)
	})
	return errors.New("")
}

func gitCommit(comment string) error {
	result := do(NO_CONDITION, fmt.Sprintf(`git commit -m "%s"`, comment))
	toContinue := checkError(result.err)
	goContinue(toContinue, func() {
		// commit
		fmt.Println("commit")
		os.Exit(0)
	})
	return errors.New("")
}

func gitPull() error {
	result := do(NO_CONDITION, "git pull --rebase")
	toContinue := checkError(result.err)
	goContinue(toContinue, func() {
		fmt.Println("pull")
		os.Exit(0)
	})
	return errors.New("")
}

func gitPush(branch string) error {
	result := do(NO_CONDITION, fmt.Sprintf("git push origin %s", branch))
	toContinue := checkError(result.err)
	goContinue(toContinue, func() {
		fmt.Println("push")
		os.Exit(0)
	})
	return errors.New("")
}

func checkoutBranch(checkout string) string {
	if len(checkout) == 0 {
		branchInfo := do(NO_CONDITION, "git branch")
		checkError(branchInfo.err)
		branches := strings.Split(branchInfo.output, "\n")
		var currentBranch string
		for _, branch := range branches {
			if strings.Contains(branch, "*") {
				index := strings.Index(branch, "*")
				// 中间隔了有两个占位符
				currentBranch = branch[index + 2:]
				break
			}
		}
		return currentBranch
	} else {
		return checkout
	}
}

func checkoutSuccess(info string) error {
	if strings.Contains(info, "error") {
		fmt.Println("checkout error info: ", info)
		os.Exit(-1)
	}
	return nil
}

func do(errorCondition, commandLine string) *Result {
	fmt.Printf(`command = "%s" will be executing\n`, commandLine)
	var cmd = exec.Command("/bin/sh", "-c", commandLine)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		//fmt.Println(stderr.String())
		return NewResult("", stderr.String())
	}
	result := out.String()
	//fmt.Println("execution successful", result)
	if strings.Contains(result, errorCondition) {
		return NewResult("", UNMATCHED_CONDITION)
	}

	return NewResult(result, "")
}

// check error
func checkError(err error) bool {
	if err != nil {
		fmt.Printf("error occurred, error = \n%s\n", err.Error())
		fmt.Println(`input "continue", system will resolve it`)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		comparator := strings.TrimRight(input, "\n")
		if comparator == CONTINUE {
			return true
		}
		os.Exit(-1)
	}
	return false
}

func goContinue(isContinued bool, fn func())  {
	if isContinued {
		fmt.Println("it will be continue, high power will release, watch out!")
		fn()
	}
}

func main() {
	//result := do("ls", NO_CONDITION, "-l")
	//fmt.Println(result.output)

	args := os.Args[1:]
	if len(args) == 2 {
		checkError(GitPush(args[0], args[1]))
	} else if len(args) == 1 {
		checkError(GitPush(args[0], EMPTY))
	} else {
		checkError(GitPush(DEFAULT_COMMENT, EMPTY))
	}
}