package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readConfig() map[string]string {
	file, err := os.Open("config/site.cfg")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	sites := make(map[string]string)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		str := strings.Split(line, ":")
		a := str[0]
		b := str[1]
		sites[a] = b
	}

	return sites
}
