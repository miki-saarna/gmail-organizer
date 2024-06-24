package utils

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	clearScreen = "\033[H\033[2J"
)

type Options []string

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

func (o *Options) SelectOption() (string, error) {
	options := *o
	selected := 0

	if err := keyboard.Open(); err != nil {
		return "", fmt.Errorf("could not access keyboard interactions: %v", err.Error())
	}
	defer keyboard.Close()

	fmt.Print(clearScreen)
	
	for {
		fmt.Println("Use the arrow keys to select an option and press Enter:")
		for i, option := range options {
			if i == selected {
				fmt.Printf("%s%-20s <--%s\n", colorGreen, option, colorReset)
			} else {
				fmt.Printf("%s\n", option)
			}
		}

		_, key, err := keyboard.GetKey()
		if err != nil {
			return "", fmt.Errorf("could not get key that was activated: %v", err.Error())
		}

		if key == keyboard.KeyArrowDown {
			selected = (selected + 1) % len(options)
		} else if key == keyboard.KeyArrowUp {
			selected = (selected - 1 + len(options)) % len(options)
		} else if key == keyboard.KeyEnter {
			break
		}

		fmt.Print(clearScreen)
	}

	return options[selected], nil
}