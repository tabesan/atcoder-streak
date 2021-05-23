package main

import (
	cmt "atcoder-streak/commit"
	tm "atcoder-streak/timer"
	"context"
	"time"
)

func main() {
	client := cmt.NewClient(cmt.Name, cmt.Repository)
	cmt.NewGetter(client)
	client.Getter.InitStreak()
	client.ShowStreak()
	timer := tm.NewTimer()
	go timer.FlagTimer()
	go timer.UpdateTimer()

	for {
		select {
		case <-timer.ChFlag:
			client.ResetFlag()
		case <-timer.ChUpdate:
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
			defer cancel()
			client.UpdateStreak(ctx)
			select {
			case <-ctx.Done():
				client.Timeouted()
			}
		}
	}

}
