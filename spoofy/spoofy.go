package spoofy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Spoof (ext, fn string, paste bool) error {
	scanner := bufio.NewScanner(os.Stdin)
	var doc string
	var empty int
	fmt.Println(`**** type "end..." (on a new line) after pasting text or press enter key 16 times,
	   to exit when in paste mode. ***`)
	fmt.Println("Enter text here: ")
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			empty += 1
			if !paste {
				break
			}
			if empty >= 15 {
				break
			}
		}else {
			if strings.Trim(line, " ")=="end..." {
				break
			}
			empty = 0
		}	
		doc += line+"\n"
	}
	
	if scanner.Err() != nil {
		fmt.Println(scanner.Err(), "scanner error")
	}
	if len(doc) != 0 {
		if err := os.WriteFile(fn+"."+ext, []byte(doc), 0644); err != nil {
		fmt.Println(err)
		return err
	}
	}
	return nil
}