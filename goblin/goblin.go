package goblin

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dglitxh/tyrenz/helpers"
)

var logfn *os.File = helpers.CreateLogFile(".goblogs")

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
	for i, v := range *p {
		if err := v.Execute(); err != nil {
			helpers.Logger(logfn, fmt.Sprintf("Error: %v at step %d (%s)" , err, i, v.Name))
			return err
		}
	}
	return nil 
}