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
	newStreak := false
	for {
		streak, latest, updateFlag, resetFlag, err := client.DownloadData()
		if err != nil {
			time.Sleep(time.Minute * 30)
			continue
		}

		client.SetStreak(streak)
		client.SetLatest(latest)
		client.SetUpdateFlag(updateFlag)
		client.SetResetFlag(resetFlag)
		err = client.InitStreak()
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

		if latest != client.ReferLatestCommit() {
			newStreak = true
			err := client.UploadData()
			if err != nil {
				time.Sleep(time.Minute * 30)
				continue
			}
		}
	}

	if newStreak {
		msg := "Current streak: " + strconv.Itoa(client.ReferStreak()) + "days"
		notify.SendNotify(msg)
	}

	timer := tm.NewTimer()
	go timer.FlagTimer()
	go timer.UpdateTimer()

	for {
		select {
		case <-timer.ChFlag:
			msg := "call ResetFlag()"
			notify.SendNotify(msg)
			client.ResetFlag()
			client.UploadData()
		case <-timer.ChUpdate:
			if !client.ReferUpdateFlag() {
				err := client.UpdateStreak()
				if err != nil {
					msg := "Update error"
					notify.SendNotify(msg)
				}
				if client.ReferUpdateFlag() {
					msg := "\nCurrent streak " + strconv.Itoa(client.ReferStreak()) + "days"
					notify.SendNotify(msg)
					client.UploadData()
				}
			}
		}
	}

}
