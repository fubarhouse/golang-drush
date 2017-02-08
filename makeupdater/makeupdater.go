package makeupdater

// Note this package is exclusively compatible with Drupal 7 make files.

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

func replaceTextInFile(fullPath string, oldString string, newString string) {
	read, err := ioutil.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}
	newContents := strings.Replace(string(read), oldString, newString, -1)
	err = ioutil.WriteFile(fullPath, []byte(newContents), 0)
	if err != nil {
		panic(err)
	}
}

func removeChar(input string, chars ...string) string {
	for _, value := range chars {
		input = strings.Replace(input, value, "", -1)
	}
	return input
}

func UpdateMake(fullpath string) {
	// Update the version numbers in a specified make file
	projects := GetProjectsFromMake(fullpath)
	count := 0
	for _, project := range projects {
		if project != "" {

			catCmd := "cat " + fullpath + " | grep \"projects\\[" + project + "\\]\" | grep version | cut -d '=' -f2"
			z, _ := exec.Command("sh", "-c", catCmd).Output()
			versionOld := removeChar(string(z), " ", "\"", "\n")
			x, _ := exec.Command("sh", "-c", "drush pm-releases --pipe "+project+" | grep Recommended | cut -d',' -f2").Output()
			versionNew := removeChar(string(x), " ", "7.x-", "\"", "\n", "[", "]")
			if versionOld != versionNew && versionOld != "" && versionNew != "" {
				fmt.Printf("Replacing %v v%v with v%v\n", project, versionOld, versionNew)
				count++
				replaceTextInFile(fullpath, fmt.Sprintf("projects[%v][version] = \"%v\"\n", project, versionOld), fmt.Sprintf("projects[%v][version] = \"%v\"\n", project, versionNew))
			}
		}
	}
	if count != 0 {
		fmt.Println(fullpath, "is already up to date.")
	}
}

func GetProjectsFromMake(fullpath string) []string {
	// Returns a list of projects from a given make file
	Projects := []string{}
	catCmd := "cat " + fullpath + " | grep projects | cut -d'[' -f2 | cut -d']' -f1 | uniq | sort"
	y, _ := exec.Command("sh", "-c", catCmd).Output()
	rawProjects := strings.Split(string(y), "\n")
	for _, project := range rawProjects {
		if project != "" {
			Projects = append(Projects, project)
		}
	}
	return Projects
}

func GenerateMake(Projects []string, File string) {
	// Takes a []string of projects and writes out a make file
	// Modules are added with the latest recommended version.
	headerLines := []string{}
	headerLines = append(headerLines, "; Generated by make-updater")
	headerLines = append(headerLines, "; Script created by Fubarhouse")
	headerLines = append(headerLines, "; Toolkit available at github.com/fubarhouse/golang-drush/...")
	headerLines = append(headerLines, "core = 7.x")
	headerLines = append(headerLines, "api = 2")
	headerLines = append(headerLines, "")

	// Rewrite core, if core is in the original Projects list.

	for _, Project := range Projects {
		coreAppended := 0
		if Project == "drupal" {
			if coreAppended == 0 {
				headerLines = append(headerLines, "; core")
				x, _ := exec.Command("sh", "-c", "drush pm-releases --pipe drupal | grep Recommended | cut -d',' -f2").Output()
				ProjectVersion := removeChar(string(x), " ", "7.x-", "\"", "\n", "[", "]")
				headerLines = append(headerLines, "projects[drupal][type] = \"core\"")
				headerLines = append(headerLines, fmt.Sprintf("projects[drupal][version] = \"%v\"", ProjectVersion))
				headerLines = append(headerLines, "projects[drupal][download][type] = \"get\"")
				headerLines = append(headerLines, fmt.Sprintf("projects[drupal][download][url] = \"https://ftp.drupal.org/files/projects/drupal-%v.tar.gz\"", ProjectVersion))
				headerLines = append(headerLines, "")
				coreAppended++
			}
		}
	}

	// Rewrite contrib
	headerLines = append(headerLines, "; modules")

	for _, Project := range Projects {

		if Project != "drupal" {
			x, y := exec.Command("sh", "-c", "drush pm-releases --pipe "+Project+" | grep Recommended | cut -d',' -f2").Output()
			if y == nil {
				ProjectVersion := removeChar(string(x), " ", "7.x-", "\"", "\n", "[", "]")
				headerLines = append(headerLines, fmt.Sprintf("projects[%v][version] = \"%v\"", Project, ProjectVersion))
			}
		}
	}

	// Print to path File

	newFile, _ := os.Create(File)
	for _, line := range headerLines {
		fmt.Fprintln(newFile, line)
	}
	newFile.Sync()
	defer newFile.Close()

}
