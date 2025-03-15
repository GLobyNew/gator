package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/GLobyNew/gator/internal/config"
	"github.com/GLobyNew/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	userConfig, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sql.Open("postgres", userConfig.DbURL)
	if err != nil {
		log.Fatalln("can't open db")
	}
	dbQueries := database.New(db)
	s := state{db: dbQueries, cfg: &userConfig}

	cmds := commands{cmd: make(map[string]func(*state, command) error)}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	cmds.register("addfeed", handleAddFeed)

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
