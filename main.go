package main

import (
	cm "atcoder-streak/commit"
	nt "atcoder-streak/notify"
	tm "atcoder-streak/timer"
	"strconv"
	"time"
)

func main() {
	notify := nt.NewNotify()
	client := cm.NewClient(cm.Name, cm.Repository)
	cm.NewGetter(client)
	initTrials := 0
	for {
		err := client.InitStreak()
		if err != nil {
			initTrials += 1
		} else {
			break
		}

		if initTrials == 10 {
			msg := "InitStreak missed 10 times"
			notify.SendNotify(msg)
			time.Sleep(1 * time.Hour)
			initTrials = 0
		}
	}

	{
		msg := "\nCurrent streak: " + strconv.Itoa(client.ReferStreak()) + "days"
		notify.SendNotify(msg)
	}
	timer := tm.NewTimer()
	go timer.FlagTimer()
	go timer.UpdateTimer()

	errCount := 0
	const errLimit = 5
	for {
		select {
		case <-timer.ChFlag:
			client.ResetFlag()
		case <-timer.ChUpdate:
			if client.ReferTimeoutFlag() {
				err := client.InitStreak()
				if err != nil {
					errCount += 1
				} else {
					errCount = 0
				}

				if errCount == errLimit {
					msg := "InitStreak() missed errLimit times"
					notify.SendNotify(msg)
					errCount = 0
				}
			}
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
