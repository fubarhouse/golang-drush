// Copyright © 2017 Karl Hepworth <Karl.Hepworth@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	c "github.com/fubarhouse/drubuild/composer"

	"errors"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Install or remove a project.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ComposerCmd, err := exec.LookPath("composer")
		if err != nil {
			log.Fatal("Composer was not found in your $PATH")
		}
		if remove {
			d := exec.Command(ComposerCmd, "remove", name)
			d.Dir = destination
			d.Stdout = os.Stdout
			d.Stderr = os.Stderr
			d.Run()
			d.Wait()
		}
		if add {
			var r string
			if preferSource {
				r = "require --prefer-source"
			} else {
				r = "require"
			}
			if version != "" {
				name += ":" + version
			}
			d := exec.Command(ComposerCmd, r, name)
			d.Dir = destination
			d.Stdout = os.Stdout
			d.Stderr = os.Stderr
			d.Run()
			d.Wait()
		}

		if !workingCopy {
			x, e := c.GetPath(destination, name)
			if e != nil {
				err := errors.New("could not find path associated to " + name)
				err.Error()
			}
			if x != "" {
				removeGitFromPath(x)
			}
		}

		if !add && !remove {
			fmt.Println("No action selected, add --add or --remove for effect.")
		}

	},
}

func init() {
	RootCmd.AddCommand(projectCmd)

	// Get the current working directory.
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	projectCmd.Flags().StringVarP(&name, "name", "n", "", "The human-readable name for this site")
	projectCmd.Flags().StringVarP(&destination, "path", "p", dir, "The path to where the site(s) exist.")
	projectCmd.Flags().StringVarP(&version, "version", "v", "", "Version of the package.")
	projectCmd.Flags().BoolVarP(&add, "add", "a", false, "Flag to trigger add action.")
	projectCmd.Flags().BoolVarP(&remove, "remove", "r", false, "Flag to trigger remove action.")
	projectCmd.Flags().BoolVarP(&workingCopy, "working-copy", "w", false, "Mark as a working-copy.")
	projectCmd.Flags().BoolVarP(&preferSource, "prefer-source", "s", false, "Build with preference to source packages.")

	projectCmd.MarkFlagRequired("name")
	projectCmd.MarkFlagRequired("path")

}
