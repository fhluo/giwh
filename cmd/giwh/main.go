package main

import (
	"github.com/fhluo/giwh/cmd"
	"log"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd.Execute()
}
