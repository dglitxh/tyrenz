package helpers

import (
	"log"
	"os"
)


func Logger (txt... interface{}) {
	f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	var output string
	for _, v := range txt {
		output += " "+v.(string)
	}
	log.Println(output)
}

