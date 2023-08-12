package goblin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/dglitxh/tyrenz/helpers"
)

var logfn *os.File = helpers.CreateLogFile(".goblogs")

type Step struct {
	Name string `json:"name"`
	Dir string	`json:"dir"`
	Cmd string	`json:"cmd"`
	Msg string  `json:"msg"`
	Args []string `json:"args"`
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

func (s Step) TimeOutExecute () error{
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	task := exec.CommandContext(ctx, s.Cmd, s.Args...)
	task.Dir = s.Dir
	if err := task.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			helpers.Logger(logfn, fmt.Sprintf("Error: %v, timeout", err))
			return fmt.Errorf("error: %v, timeout", err)
		}
		return err
	}
	return nil 
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