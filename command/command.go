package command

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os/exec"
	"path/filepath"
	"strings"
)

type Command struct {
	// Our structured data/object for Command
	alias   string
	command string
	verbose bool
}

const PATH_DRUSH = "/usr/local/bin/drush"

// Creates a new container for []Command objects
func NewDrushCommand() *Command {
	return &Command{}
}

// Changes all given values for a Drush command object.
func (drush *Command) Set(alias string, command string, verbose bool) {
	drush.alias = alias
	drush.command = command
	drush.verbose = verbose
}

// Returns the alias used to executed Drush commands.
func (drush *Command) GetAlias() string {
	return drush.alias
}

// Changes the alias used to executed Drush commands.
func (drush *Command) SetAlias(value string) {
	drush.alias = value
}

// Returns the command string on executed Drush commands.
func (drush *Command) GetCommand() string {
	return drush.command
}

// Changes the command string on executed Drush commands.
func (drush *Command) SetCommand(value string) {
	drush.command = value
}

// Returns the verbosity setting on executed Drush commands.
func (drush *Command) GetVerbose() bool {
	return drush.verbose
}

// Changes the verbosity setting on executed Drush commands.
func (drush *Command) SetVerbose(value bool) {
	drush.verbose = value
}

// Gets the output from a single Command object, does not support []Command items.
func (drush *Command) Output() ([]string, error) {
	comm, err := drush.Run()
	response := filepath.SplitList(string(comm))
	return response, err
}

// Run an individual Command object, does not support []Command items.
func (drush *Command) Run() ([]byte, error) {
	if strings.Contains(drush.alias, "@") == true {
		drush.alias = strings.Replace(drush.alias, "@", "", -1)
	}
	if drush.alias != "" {
		drush.alias = fmt.Sprintf("@%v", drush.alias)
	}
	if drush.verbose == true {
		drush.alias = fmt.Sprintf("%v --verbose", drush.alias)
	}
	args := fmt.Sprintf("%v %v", drush.alias, drush.command)
	comm, err := exec.Command("sh", "-c", PATH_DRUSH+" "+args).Output()
	return comm, err
}

// Executes a database synchronisation task from a source to destination with the use of Drush.
func DrushDatabaseSync(srcAlias, destAlias string) {
	/* So our binary and this function combined support two-way traffic...  */
	x := NewDrushCommand()
	srcAlias = strings.Replace(srcAlias, "@", "", -1)
	destAlias = strings.Replace(destAlias, "@", "", -1)
	x.Set("", fmt.Sprintf("sql-sync @%v @%v -y", srcAlias, destAlias), true)
	_, err := x.Output()
	if err == nil {
		log.Infoln("Syncronised databases complete.")
	} else {
		log.Errorln("Could not syncronise databases.")
	}
}

// Executes a file synchronisation task from a source to destination with the use of Drush.
func DrushFilesSync(srcAlias, destAlias string) {
	x := NewDrushCommand()
	srcAlias = strings.Replace(srcAlias, "@", "", -1)
	destAlias = strings.Replace(destAlias, "@", "", -1)
	x.Set("", fmt.Sprintf("--yes rsync --exclude-other-sites --exclude-conf @%v:%%files @%v:%%files", srcAlias, destAlias), true)
	_, err := x.Output()
	if err == nil {
		log.Infoln("Synced public file system.")
	} else {
		log.Warnln("Public file system has not been synced.")
	}
	x.Set("", fmt.Sprintf("--yes rsync --exclude-other-sites --exclude-conf @%v:%%private @%v:%%private", srcAlias, destAlias), true)
	_, err = x.Output()
	if err == nil {
		log.Infoln("Synced private file system.")
	} else {
		log.Warnln("Private file system has not been synced.")
	}
}

// Performs a cache clear task on an input site alias with the use of Drush.
func DrushClearCache(alias string) {
	drushCommand := NewDrushCommand()
	drushCommand.Set(alias, "cc all", false)
	_, err := drushCommand.Output()
	if err != nil {
		log.Warnln("Could not clear caches.", err)
	} else {
		log.Infoln("Caches cleared.")
	}
}

// Performs a registry rebuild task on an input site alias with the use of Drush.
func DrushRebuildRegistry(alias string) {
	drushCommand := NewDrushCommand()
	drushCommand.Set(alias, "rr", false)
	_, err := drushCommand.Output()
	if err != nil {
		log.Warnln("Could not rebuild registry.", err)
	} else {
		log.Infoln("Rebuilt registry.")
	}
}

// Performs a database update task on an input site alias with the use of Drush.
func DrushUpdateDatabase(alias string) {
	drushCommand := NewDrushCommand()
	drushCommand.Set(alias, "updb -y", false)
	_, err := drushCommand.Output()
	if err != nil {
		log.Warnln("Could not update database.", err)
	} else {
		log.Infoln("Updated database where possible.")
	}
}
