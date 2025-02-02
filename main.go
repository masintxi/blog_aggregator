package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/masintxi/blog_aggregator/internal/config"
	"github.com/masintxi/blog_aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
		return
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatal("error opening connexion to database")
	}
	defer db.Close()

	programState := &state{
		db:  database.New(db),
		cfg: &cfg,
	}

	cmds := commands{
		regisCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", registerUser)
	cmds.register("reset", resetUsersTable)
	cmds.register("users", getUsersList)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", handleNewFeed)
	cmds.register("feeds", getFeedsList)

	if len(os.Args) < 2 {
		log.Fatal("Not enough arguments")
		return
	}
	cmd := command{os.Args[1], os.Args[2:]}

	if err := cmds.run(programState, cmd); err != nil {
		log.Fatalf("error: %v\n", err)
		return
	}

}
