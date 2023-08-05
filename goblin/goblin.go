package goblin

import (
	"context"
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

func (p *Process) NewStep (name, dir, cmd, msg string, args []string) Step {
	step := Step{
		name: name,
		dir: dir,
		msg: msg,
		cmd: msg,
		args: args,
	}
	*p = append(*p, step)
	return step
}

func (s Step) Execute () error{
	task := exec.Command(s.cmd, s.args...)
	task.Dir = s.dir
	if err := task.Run(); err != nil {
		return err
	}
	return nil 
}

func (p *Process) Run (ctx context.Context) error {
	for _, v := range *p {
		v.Execute()
	}
	return nil 
}