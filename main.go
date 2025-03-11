package main

import (
	"fmt"
	"log"

	"github.com/GLobyNew/gator/internal/config"
)

func main() {
	userConfig, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	err = userConfig.SetUser("test")
	if err != nil {
		log.Fatalln(err)
	}

	userConfig, err = config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%v\n%v\n", userConfig.DbURL, userConfig.CurrentUserName)
}
