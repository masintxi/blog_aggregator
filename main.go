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
	db       *database.Queries
	cfg      *config.Config
	newFeeds chan database.Feed
	done     chan struct{}
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config file: %v", err)
		return
	}

	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error opening connexion to database: %v", err)
	}
	defer db.Close()

	programState := &state{
		db:       database.New(db),
		cfg:      &cfg,
		newFeeds: make(chan database.Feed, 100),
		done:     make(chan struct{}),
	}

	cmds := commands{
		regisCommands: make(map[string]func(*state, command) error),
	}
	cmds.register("login", handlerLogin)
	cmds.register("register", registerUser)
	cmds.register("reset", resetUsersTable)
	cmds.register("users", getUsersList)
	cmds.register("agg", handleAgg)
	cmds.register("addfeed", middlewareLoggedIn(handleNewFeed))
	cmds.register("delfeed", deleteFeed)
	cmds.register("feeds", getFeedsList)
	cmds.register("follow", middlewareLoggedIn(handleNewFollow))
	cmds.register("following", middlewareLoggedIn(handleFollowsForUser))
	cmds.register("unfollow", middlewareLoggedIn(handleUnfollow))
	cmds.register("browse", middlewareLoggedIn(handleBrowsePosts))

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
