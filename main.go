package main

import (
	"os"
	"redisdual/cmd"
)

func main(){
	cmd.MainStart(os.Args[1:])
}
