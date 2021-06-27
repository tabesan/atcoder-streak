package main

import (
	cm "atcoder-streak/commit"
	nt "atcoder-streak/notify"
	tm "atcoder-streak/timer"
	"fmt"
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
		streak, latest, update, err := client.DownloadData()
		if client.IsUpdated(update) {
			break
		}

		if err != nil {
			time.Sleep(time.Minute * 30)
			continue
		}
		client.SetStreak(streak)
		client.SetLatest(latest)
		err = client.InitStreak()
		if err != nil {
			initTrials += 1
			continue
		}

		if initTrials == 10 {
			msg := "InitStreak missed 10 times"
			notify.SendNotify(msg)
			time.Sleep(1 * time.Hour)
			initTrials = 0
		}

		if latest != client.ReferLatestCommit() {
			fmt.Println("Latest", latest, "Ref", client.ReferLatestCommit())
			newStreak = true
			err := client.UploadData()
			if err != nil {
				for {
					time.Sleep(time.Minute * 30)
					err = client.UploadData()
					if err == nil {
						break
					}
				}
			}
		}
		break
	}

	if newStreak {
		msg := "Current streak: " + strconv.Itoa(client.ReferStreak()) + " days"
		notify.SendNotify(msg)
		err := client.UploadData()
		if err != nil {
			for {
				time.Sleep(30 * time.Minute)
				err = client.UploadData()
			}
		}
	}

	timer := tm.NewTimer()
	go timer.Timer()

	/*for {
		select {
		case <-timer.ChFlag:
			if client.ConvJST(time.Now()).Hour() == 0 && !client.ReferResetFlag() {
				msg := "call ResetFlag()"
				notify.SendNotify(msg)
				client.ResetFlag()
				client.UploadData()
			}
		case <-timer.ChUpdate:
			if client.ConvJST(time.Now()).Hour() != 0 {
				client.SetResetFlag(false)
			}
			if !client.ReferUpdateFlag() {
				if client.ReferUpdateFlag() {
					msg := "\nCurrent streak " + strconv.Itoa(client.ReferStreak()) + "days"
					notify.SendNotify(msg)
					err := client.UpdateStreak()
					if err != nil {
						msg := "Update error"
						notify.SendNotify(msg)
					}
					client.UploadData()
				}
			}
		}
	}*/

	for {
		select {
		case <-timer.ChUpdate:
			latest := client.ReferLatestCommit()
			_, _, update, err := client.DownloadData()
			if err != nil || client.IsUpdated(update) {
				continue
			}

			err = client.InitStreak()
			if err != nil {
				for {
					time.Sleep(30 * time.Minute)
					err = client.InitStreak()
					if err == nil {
						break
					}
				}
			}

			if latest == client.ReferLatestCommit() {
				continue
			}

			msg := "Current Streak " + strconv.Itoa(client.ReferStreak()) + " days"
			notify.SendNotify(msg)
			err = client.UploadData()
			if err != nil {
				for {
					time.Sleep(30 * time.Minute)
					err = client.UploadData()
				}
			}
		default:
		}
	}
}
