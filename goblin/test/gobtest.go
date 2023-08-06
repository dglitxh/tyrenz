package gobliin_test

import (
	"fmt"
	"testing"

	"github.com/dglitxh/tyrenz/goblin"
)

func  TestAddStep (t *testing.T) {
	proc := make(goblin.Process, 0)
	step := goblin.Step{
		Name: "make directory",
		Dir: ".",
		Msg: "create a new dir",
		Cmd: "mkdir temp",
		Args: make([]string, 0),

	}
	
	t.Run(step.Name, func(t *testing.T) {
		s := proc.NewStep(step)
		if step.Name != s.Name {
			fmt.Errorf()
		}
	})
}