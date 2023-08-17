package goblin

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	var argstr string
	var err error
	for {
		fmt.Println("Enter the name of your command: ")
		name, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Enter information your command: ")
		msg, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Enter your working directory: ")
		dir, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("enter the first word in command eg: git, python, go")
		cmd, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		fmt.Println("Enter command args, each separated by a comma.")
		argstr, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return err
		}
		args = strings.Split(argstr, ",")
		step := Step{
			Name: strings.TrimSuffix(name, "\n"),
			Msg:  strings.TrimSuffix(msg, "\n"),
			Dir:  strings.TrimSuffix(dir, "\n"),
			Cmd:  strings.TrimSuffix(cmd, "\n"),
			Args: args,
		}
		p.NewStep(step)
		fmt.Println("")
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
