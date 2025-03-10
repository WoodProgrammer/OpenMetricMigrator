package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/rs/zerolog/log"
)

func FileHandler(fileName string, lines []string) {

	output := strings.Join(lines, "\n")
	err := ioutil.WriteFile(fmt.Sprintf(fileName), []byte(output), 0644)
	if err != nil {
		log.Err(err).Msg("Error while writing file")
	}
}
