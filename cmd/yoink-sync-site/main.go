package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fubarhouse/golang-drush/command"
	"os"
)

func main() {

	var SourceAlias = flag.String("source-alias", "", "Alias of source site")
	var DestAlias = flag.String("dest-alias", "", "Alias of destination site")
	var SyncDB = flag.Bool("db", false, "Mark database for syncronization")
	var SyncFiles = flag.Bool("files", false, "Mark files for syncronization")

	// Usage:
	// -local-alias="mysite.dev" \
	// -remote-alias="mysite.dev" \
	// -db \
	// -files

	flag.Parse()

	if *SourceAlias == "" {
		log.Infoln("Source input is empty")
	}
	if *DestAlias == "" {
		log.Infoln("Destination input is empty")
	}
	if !*SyncDB {
		log.Infoln("Database flag is switched off")
	}
	if !*SyncFiles {
		log.Infoln("Files flag is switched off")
	}

	if *SourceAlias == "" || *DestAlias == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *SyncDB {
		log.Infoln("Database was marked for syncing, working now...")
		command.DrushDatabaseSync(*SourceAlias, *DestAlias)
	}
	if *SyncFiles {
		log.Infoln("Files were marked for syncing, working now...")
		command.DrushFilesSync(*SourceAlias, *DestAlias)
	}
	if *SyncDB || *SyncFiles {
		log.Infoln("Attempting to rebuild registries...")
		command.DrushRebuildRegistry(*DestAlias)
	}
}
