package gobliin_test

import (
	"fmt"
	"testing"

	"github.com/dglitxh/tyrenz/goblin"
)

var proc goblin.Process

func  TestAddStep (t *testing.T) {
	
	step1 := goblin.Step{
		Name: "make directory",
		Dir: ".",
		Msg: "create a new dir",
		Cmd: "mkdir", 
		Args: []string{"temp"},

	}
	step2 := goblin.Step{
		Name: "add file",
		Dir: "./temp",
		Msg: "add new file to dir",
		Cmd: "touch", 
		Args: []string{"temp.txt", ".gitignore"},
	}
	step3 := goblin.Step{
		Name: "git initialization",
		Dir: ".",
		Msg: "initialize git",
		Cmd: "git", 
		Args: []string{"init"},
	}
	
	
	t.Run(step1.Name, func(t *testing.T) {
		s := proc.NewStep(step1)
		if step1.Name != s.Name {
			t.Errorf("expected %v to be %v", step1.Name, s.Name)
		}
	})
		t.Run(step2.Name, func(t *testing.T) {
		s := proc.NewStep(step2)
		if step2.Name != s.Name {
			t.Errorf("expected %v to be %v", step2.Name, s.Name)
		}
	})
	t.Run(step3.Name, func(t *testing.T) {
		s := proc.NewStep(step3)
		if step3.Name != s.Name {
			t.Errorf("expected %v to be %v", step3.Name, s.Name)
		}
	})
}

func TestRun(t *testing.T) {
	fmt.Println(proc)
	t.Run("Execution", func(t *testing.T) {
	if err := proc.Run(); err != nil {
		t.Errorf("%v occurred while executing tasks", err)
	}
	})
}