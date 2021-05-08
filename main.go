package main

import (
	ats "atcoder-streak/commit"
)

func main() {
	client := ats.NewClient(ats.Name, ats.Repository)
	client.InitStreak()
	client.ShowStreak()
}
