package main

import (
	"flag"
	"github.com/fubarhouse/golang-drush/makeupdater"
	"strings"
)

func main() {
	var strMake = flag.String("makes", "", "Comma-separated list of absolute paths to make files to update.")
	flag.Parse()
	if *strMake != "" {
		Makes := strings.Split(*strMake, ",")
		for _, Makefile := range Makes {
			makeupdater.UpdateMake(Makefile)
		}
	} else {
		flag.Usage()
	}
}
