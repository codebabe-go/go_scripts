package common

import (
	"fmt"
	"strings"
	"os"
	"os/exec"
	"runtime"
	"bytes"
	"errors"
)

const OS_X string = "darwin"
const WINDOWS string = "windows"

const NO_CONDITION string = "no condition, cannot match"
const UNMATCHED_CONDITION string = "cannot match condition"

type Result struct {
	Output string
	Err error
}

func NewResult(output string, errMsg string) *Result {
	result := new(Result)
	result.Output = output
	if len(errMsg) != 0 {
		result.Err = errors.New(errMsg)
	}
	return result
}

func Do(errorCondition, commandLine string) *Result {
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