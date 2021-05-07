package main

import (
	ats "atcoder-streak/commit"
	"fmt"
)

func main() {
	client := ats.NewClient(ats.Name, ats.Repository)
	client.GETRequest()
	fmt.Println(client.Data.GetDate())
}
