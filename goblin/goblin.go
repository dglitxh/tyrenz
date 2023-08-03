package goblin

import (
	"os/exec"
)

type Step struct {
	name string
	dir string
	cmd string
	msg string 
	args []string
}

type Process []Step

func Execute () {
	exec.Command()
}