package helpers

import (
	"log"
	"os"
	"strconv"
)

func CreateLogFile (name string) *os.File {
	f, err := os.OpenFile(name+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func Logger (f *os.File, txt... interface{}) {	
	defer f.Close()
	log.SetOutput(f)
	var output string
	for _, v := range txt {
		if _, ok := v.(string); ok {
			output += " "+v.(string)
		}else if _, ok := v.(int); ok {
			output += " "+strconv.Itoa(v.(int))
		}else if _, ok := v.(bool); ok {
			output += " "+strconv.FormatBool(v.(bool))
		}
	}
	log.Println(output)
}

