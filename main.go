package main

import (
	cmt "atcoder-streak/commit"
	tm "atcoder-streak/timer"
	"fmt"
)

func main() {
	client := cmt.NewClient(cmt.Name, cmt.Repository)
	client.InitStreak()
	timer := tm.NewTimer()
	go timer.FlagTimer()
	go timer.UpdateTimer()

	for {
		select {
		case <-timer.ChFlag:
			client.FlagReset()
			fmt.Println("flag")
		case <-timer.ChUpdate:
			client.UpdateStreak()
			fmt.Println("update")
		}
	}

}
