package main

import (
	"log"
	"os"

	"github.com/GLobyNew/gator/internal/config"
)

func main() {
	userConfig, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	s := state{config: &userConfig}
	cmds := commands{cmd: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)

	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatalln("no command provided")
	}

	cmd := command{name: args[0], args: args[1:]}
	err = cmds.run(&s, cmd)
	if err != nil {
		log.Fatalln(err)
	}
}
