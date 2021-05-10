package main

import (
	ats "atcoder-streak/commit"
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
		case <-timer.ChUpdate:
			client.UpdateStreak()
		}
	}

}
