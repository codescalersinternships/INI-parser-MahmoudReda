package main

import (
	"fmt"
)

func main() {
	// Parse the INI file
	config, err := ParseINIFile("config.ini")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Print the configuration data
	fmt.Println(config)
}
