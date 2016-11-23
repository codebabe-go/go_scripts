package main

import (
	"time"
	"./common"
	"fmt"
)

// 定时清理桌面
// 定时任务的单位默认为秒(s)
func CleanDesktop(t time.Duration, fn func()) {
	ticker := time.NewTicker(time.Second * t)
	for range ticker.C {
		fn()
	}
}

func removeIntoTrash() {
	dirInfo := common.Do(common.NO_CONDITION, "ls ~/Desktop/")
	if len(dirInfo.Output) != 0 {
		mvResult := common.Do(common.NO_CONDITION, "mv ~/Desktop/* ~/.Trash")
		common.CheckError(mvResult.Err)
	} else {
		fmt.Printf("current time = %v, desktop is clean\n", time.Now())
	}

}

func main() {
	//args := os.Args
	//if len(args) > 1 {
	//	duration := args[1]
	//	CleanDesktop(duration, removeIntoTrash)
	//}

	CleanDesktop(1, func() {
		fmt.Println("tick")
	})
}
