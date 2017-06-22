package main

/**
 * git 一键提交工具
 * 期间一旦遇到需要手动解决的错误, 就会丢出异常终止程序
 */

import (
	"os"
	"os/exec"
	"strings"
	"fmt"
	"errors"
	"bytes"
	"bufio"
	"runtime"
)

// TODO: 1.完善各个continue func的功能 2.给出详细的下一步操作提示 3.给出参数完善日志输出(是否要输出到文件中) 4.兼容Windows

const NO_CONDITION string = "no condition, cannot match"
const DEFAULT_COMMENT string = "commit"
const UNMATCHED_CONDITION string = "cannot match condition"
const CONTINUE string = "continue"
const EMPTY string = ""

const CANNOT_PUSH_BRANCH string = "release, master"

const OS_X string = "darwin"
const WINDOWS string = "windows"

var ERROR_MAP map[string][]string = make(map[string][]string)

// error map information init
func init()  {
	ERROR_MAP["checkout"] = []string{}
	ERROR_MAP["add"] = []string{}
	ERROR_MAP["commit"] = []string{}
	ERROR_MAP["pull"] = []string{}
	ERROR_MAP["push"] = []string{}
}

type Result struct {
	output string
	err error
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
	if strings.Contains(CANNOT_PUSH_BRANCH, branch) {
		fmt.Printf("you want to push to [[%s]], confirm to push y/n \n", branch)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if strings.Compare(strings.ToLower(strings.TrimRight(input, "\n")), "y") != 0 {
			fmt.Println("git push abort!")
			os.Exit(0)
		}
	}
	fmt.Printf("branch = %s will be checked out\n", branch)
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
	return nil
}

func gitAdd() error {
	result := do(NO_CONDITION, "git add .")
	toContinue := checkError(result.err)
	goContinue(toContinue, func() {
		// commit
		fmt.Println("add")
		os.Exit(0)
	})
	return nil
}

func gitCommit(comment string) error {
	result := do(NO_CONDITION, fmt.Sprintf(`git commit -m "%s"`, comment))
	toContinue := checkError(result.err)
	goContinue(toContinue, func() {
		// commit
		fmt.Println("commit")
		os.Exit(0)
	})
	return nil
}

func gitPull() error {
	result := do(NO_CONDITION, "git pull --rebase")
	toContinue := checkError(result.err)
	goContinue(toContinue, func() {
		fmt.Println("pull")
		os.Exit(0)
	})
	return nil
}

func gitPush(branch string) error {
	result := do(NO_CONDITION, fmt.Sprintf("git push origin %s", branch))
	toContinue := checkError(result.err)
	goContinue(toContinue, func() {
		fmt.Println("push")
		os.Exit(0)
	})
	return nil
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
		fmt.Printf("checkout error info: \n%s", info)
		os.Exit(-1)
	}
	return nil
}

func do(errorCondition, commandLine string) *Result {
	fmt.Printf("command = '%s' will be executing\n", commandLine)
	var cmd *exec.Cmd
	var osInfo = runtime.GOOS
	if osInfo == WINDOWS {
		cmd = exec.Command("cmd.exe", "/c", commandLine)
	} else if osInfo == OS_X {
		cmd = exec.Command("/bin/sh", "-c", commandLine)
	} else {
		fmt.Printf("cannot match your os = %s\n", osInfo)
		os.Exit(-1)
	}
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

// 检查当前的git环境, 以error来标识是否正常
func CheckGitEnvironment() {
	result := do(NO_CONDITION, "git status")
	if result.err != nil {
		if runtime.GOOS != WINDOWS {
			path := strings.TrimRight(do(NO_CONDITION, "pwd").output, "\n")
			fmt.Printf("current path = %s which is not git repository, you shoul git init first\n", path)
		} else {
			fmt.Println("current path is not git repository, you shoul git init first")
		}
		os.Exit(-1)
	}
}

func main() {
	//result := do("ls", NO_CONDITION, "-l")
	//fmt.Println(result.output)

	CheckGitEnvironment()
	fmt.Printf("current os is %s\n", runtime.GOOS)
	args := os.Args[1:]
	if len(args) == 2 {
		checkError(GitPush(args[0], args[1]))
	} else if len(args) == 1 {
		checkError(GitPush(args[0], EMPTY))
	} else {
		checkError(GitPush(DEFAULT_COMMENT, EMPTY))
	}
	fmt.Println("push success!")
}