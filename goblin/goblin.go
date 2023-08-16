package goblin

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/dglitxh/tyrenz/helpers"
)

var logfn *os.File = helpers.CreateLogFile(".goblogs")

type Step struct {
	Name string   `json:"name"`
	Dir  string   `json:"dir"`
	Cmd  string   `json:"cmd"`
	Msg  string   `json:"msg"`
	Args []string `json:"args"`
}

type Process []Step

func (p *Process) NewStep(s Step) Step {
	step := Step{
		Name: s.Name,
		Dir:  s.Dir,
		Msg:  s.Msg,
		Cmd:  s.Cmd,
		Args: s.Args,
	}
	*p = append(*p, step)
	return step
}

func ScanText(p *Process) error {
	var reader *bufio.Reader = bufio.NewReader(os.Stdin)
	var name string
	var msg string
	var cmd string
	var dir string
	var args []string
	var err error
	for {
		name, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		msg, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		dir, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		cmd, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		// args, err = append(args, ) reader.ReadString('\n'); if err != nil {
		// fmt.Println(err)
		// return err
		// }

	}
}

func (s Step) TimeOutExecute() error {
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

func (s Step) Execute() error {
	task := exec.Command(s.Cmd, s.Args...)
	task.Dir = s.Dir
	if err := task.Run(); err != nil {
		return err
	}
	return nil
}

func (p *Process) Run(timeout bool) error {

	for i, v := range *p {
		if !timeout {
			if err := v.Execute(); err != nil {
				helpers.Logger(logfn, fmt.Sprintf("Error: %v at step %d (%s)", err, i, v.Name))
				return err
			}
		} else {
			if err := v.TimeOutExecute(); err != nil {
				helpers.Logger(logfn, fmt.Sprintf("Error: %v at step %d (%s)", err, i, v.Name))
				return err
			}
		}

	}
	return nil
}
