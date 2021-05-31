package main

import (
	cm "atcoder-streak/commit"
	tm "atcoder-streak/timer"
	"fmt"
)

func main() {
	client := cm.NewClient(cm.Name, cm.Repository)
	cm.NewGetter(client)
	client.InitStreak()
	client.ShowStreak()
	timer := tm.NewTimer()
	go timer.FlagTimer()
	go timer.UpdateTimer()

	for {
		select {
		case <-timer.ChFlag:
			client.ResetFlag()
		case <-timer.ChUpdate:
			err := client.UpdateStreak()
			if err != nil {
				fmt.Println(err)
			}
		}
	}

}
