package goblin

import (
	"fmt"
	"os/exec"
)

type Step struct {
	Name string
	Dir string
	Cmd string
	Msg string 
	Args []string
}

type Process []Step

func (p *Process) NewStep (s Step) Step {
	step := Step{
		Name: s.Name,
		Dir: s.Dir,
		Msg: s.Msg,
		Cmd: s.Cmd,
		Args: s.Args,
	}
	*p = append(*p, step)
	return step
}

func (s Step) Execute () error{
	task := exec.Command(s.Cmd, s.Args...)
	task.Dir = s.Dir
	if err := task.Run(); err != nil {
		return err
	}
	return nil 
}

func (p *Process) Run () error {
	for _, v := range *p {
		if err := v.Execute(); err != nil {
			fmt.Printf("%v at %s", err, v.Name)
			return err
		}
	}
	return nil 
}