package main

import (
	"os"

	"github.com/MrXMMM/E-commerce-Project/config"
	"github.com/MrXMMM/E-commerce-Project/modules/servers"
	"github.com/MrXMMM/E-commerce-Project/pkg/databases"
)

func envPath() string {
	if (len(os.Args)) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}

}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close() //it will do before main function return

	// start server
	servers.NewServer(cfg, db).Start()
}
