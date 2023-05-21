package todo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)
var fn string = "todo.json"
var reader *bufio.Reader = bufio.NewReader(os.Stdin)
type Todo struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Completed  bool `json:"completed"`
	Created time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

type TodoList []Todo

func (tl *TodoList) AddTodo (completed bool) error {

	fmt.Println("Please enter new todo title: ")
	title, err1 := reader.ReadString('\n'); if err1 != nil {
		fmt.Println(err1)
		return err1
	}
	fmt.Println("Please enter todo description: ")
	desc, err2 := reader.ReadString('\n'); if err2 != nil {
		fmt.Println(err2)
		return err2
	}
	desc = strings.TrimSuffix(desc, "\n")
	title = strings.TrimSuffix(title, "\n")
	item := Todo{
		Id: len(*tl)+1,
		Title: title,
		Completed: completed,
		Description: desc,
		Created: time.Now(),
		Modified: time.Time{},
	}
	*tl = append(*tl, item)
	tl.SaveTodo()
	fmt.Println("Task added succesfully.")
	return nil
}

func (tl *TodoList) ReadTodo () error {

	f, err := os.ReadFile(fn); if err != nil {
		return err
	}
	if err2 := json.Unmarshal(f, tl); err2 != nil {
		return err2
	}
	return nil
}

func (tl *TodoList) ListTodo () error {
	for _, v := range *tl {
		val, err := json.MarshalIndent(v, "  ", "  "); if err != nil {
			return err
		}
		fmt.Println(string(val))
	}
  return nil
}

func (tl *TodoList) GetTodo (id string) error {
	for _, v := range *tl {
		i, err := strconv.Atoi(id); if err != nil {
			return err
		}
		if v.Id == i {
			val, err := json.MarshalIndent(v, "  ", ""); if err != nil {
			return err
		}
		 fmt.Println(string(val))
		}
	}
  return nil
}

func (tl *TodoList) DeleteTodo (id string) error {
	list := make([]Todo, 0)
	for _, v := range *tl {
		i, err := strconv.Atoi(id); if err != nil {
			return err
		}
		if v.Id != i {
			list = append(list, v)
		}
	}
    *tl = list
	tl.SaveTodo()
	fmt.Println("task deleted succesfully")
	return nil
}

func (tl *TodoList) ToggleComplete (id string) error {
	list := make([]Todo, 0)
	for _, v := range *tl {
		i, err := strconv.Atoi(id); if err != nil {
			return err
		}
		if v.Id == i {
			v.Completed = !v.Completed
		}
		list = append(list, v)	
}
	*tl = list
	tl.SaveTodo()
	fmt.Println("task completion succesfully toggled")
	return nil
}

func (tl *TodoList) SaveTodo() error {
	it, err := json.Marshal(tl); if err != nil {
		return err
	}
	os.WriteFile(fn, []byte(it), 0644)
	return nil
}

func (tl *TodoList) EditTodo(id string) error {
	list := make([]Todo, 0)

	fmt.Println("Please enter new title: ")
	title, err1 := reader.ReadString('\n'); if err1 != nil {
		fmt.Println(err1)
		return err1
	}
	fmt.Println("Please enter new description: ")
	desc, err2 := reader.ReadString('\n'); if err2 != nil {
		fmt.Println(err2)
		return err2
	}
	for _, v := range *tl {
		i, err := strconv.Atoi(id); if err != nil {
			return err
		}
		if v.Id == i {
			if len(desc) > 0{
				v.Description = strings.TrimSuffix(desc, "\n")
			}
			if len(title) > 0{
				v.Title = strings.TrimSuffix(title, "\n")
			}
		}
		list = append(list, v)	
}
    *tl = list
	tl.SaveTodo()
	fmt.Println("Task updated succesfully")
	return nil
}