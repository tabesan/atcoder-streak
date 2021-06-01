package main

import (
	cm "atcoder-streak/commit"
	nt "atcoder-streak/notify"
	tm "atcoder-streak/timer"
	"strconv"
)

func main() {
	notify := nt.NewNotify()
	client := cm.NewClient(cm.Name, cm.Repository)
	cm.NewGetter(client)
	client.InitStreak()
	{
		msg := "\nCurrent streak: " + strconv.Itoa(client.ReferStreak()) + "days"
		notify.SendNotify(msg)
	}
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
				msg := "Update error"
				notify.SendNotify(msg)
			}
			msg := "Current streak \n" + strconv.Itoa(client.ReferStreak()) + "days"
			notify.SendNotify(msg)
		}
	}

}
