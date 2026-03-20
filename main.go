package main

import (
	"github.com/rajathjn/deny-by-default-as-a-service/cmd"
)

func main() {
	address := "0.0.0.0:8080"
	cmd.Server(address)
}
