package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

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