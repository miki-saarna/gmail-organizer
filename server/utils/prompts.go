package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type ConfirmationMsg string

func (c *ConfirmationMsg) AskForConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s \n\n[y/n]: ", *c)

		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		input = strings.ToLower(strings.TrimSpace(input))

		if input == "y" {
			return true
		} else if input == "n" {
			return false
		} else {
			fmt.Printf("Invalid input: %s. Please enter only \"y\" or \"n\".\n", input)
		}
	}
}