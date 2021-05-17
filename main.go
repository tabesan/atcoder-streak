package main

import (
	cmt "atcoder-streak/commit"
	tm "atcoder-streak/timer"
	"context"
	"fmt"
	"time"
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
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			client.UpdateStreak(ctx)
			select {
			case <-ctx.Done():
				fmt.Println("timeout")
				client.Timeouted()
			}
		}
	}

}
