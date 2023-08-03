package goblin

import (
	"os/exec"
)

type Process struct {
	name string
	dir string
	cmd string
	msg string 
	args []string
}

func Execute () {
	exec.Command()
}