package spoofy

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Create (ext string, paste bool) error {
	scanner := bufio.NewScanner(os.Stdin)
	var doc string
	var emt int
	fmt.Println("****Press enter key 3 times to complete text.***")
	fmt.Println("Enter text here: ")
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			emt += 1
			if !paste {
				break
			}
			if emt >= 5 || strings.Trim(line, " ") == "end..." {
				break
			}
		}else {
			emt = 0
		}	
		fmt.Println(strings.Trim(line, " "))
		doc += line+"\n"
	}
	
	if scanner.Err() != nil {
		fmt.Println(scanner.Err(), "scanner error")
	}
	if len(doc) != 0 {
		if err := os.WriteFile("file."+ext, []byte(doc), 0644); err != nil {
		fmt.Println(err)
		return err
	}
	}
	return nil
}