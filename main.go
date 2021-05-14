package main

import (
	ats "atcoder-streak/commit"
	"fmt"
)

func main() {
	client := ats.NewClient(ats.Name, ats.Repository)
	client.InitStreak()
	timer := ats.NewTimer()
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
