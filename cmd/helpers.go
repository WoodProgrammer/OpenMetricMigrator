package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
)

func FileHandler(fileName string, lines []string) {

	output := strings.Join(lines, "\n")                                  // here i join and add the new lines
	err := ioutil.WriteFile(fmt.Sprintf(fileName), []byte(output), 0644) // here i write to file with new data
	if err != nil {
		log.Fatalln(err)
	}
}
