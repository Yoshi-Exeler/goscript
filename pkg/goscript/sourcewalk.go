package goscript

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type ModuleSource struct {
	Name  string   // name of the module
	Files []string // paths to all the files in the module
	Hash  string   // the hash uniquely identifying this module
}

type ApplicationSource struct {
	Name            string         // application name
	ApplicationFile string         // path to the main file
	Modules         []ModuleSource // list of local and standard libarary modules
}

// ExternalModuleSource will be used by the depgraph resolver to load this module
type ExternalModuleSource struct {
	Name    string // module name
	URL     string // url from which to clone the module
	Version string // version expression that must match the git tag version (empty=latest)
	Branch  string // branch from which to get the version (empty=master)
}

var extRegex = regexp.MustCompile(`(?mU)external ([a-zA-Z]*) from "(.*)"`)

// GetRequiredExternals returns the list of all external modules required by the application at main path
func GetRequiredExternals(mainPath string) ([]*ExternalModuleSource, error) {
	// read the main file into memory
	content, err := os.ReadFile(mainPath)
	if err != nil {
		return nil, fmt.Errorf("could not read main application file from %v with error %v", mainPath, err)
	}
	ret := []*ExternalModuleSource{}
	// regex for our external definitions
	matches := extRegex.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) != 3 {
			return nil, fmt.Errorf("matched invalid external definition")
		}
		// if a version and or branch tag exist, process them
		version := ""
		branch := ""
		// split around the @ tag if we have one
		if strings.Contains(match[2], "@") {
			parts := strings.Split(match[2], "@")
			if len(parts) != 2 {
				return nil, fmt.Errorf("external url %v may only contain one @ sign", match[2])
			}
			// subsplit around the / for the branch
			if strings.Contains(parts[1], "/") {
				subparts := strings.Split(parts[1], "/")
				version = subparts[0]
				branch = subparts[1]
			} else {
				version = parts[1]
			}
		}
		// save this module requirement
		ret = append(ret, &ExternalModuleSource{
			Name:    match[1],
			URL:     match[2],
			Version: version,
			Branch:  branch,
		})
	}
	return ret, nil
}

// SourceWalk will parse the imports required by the application at the specified path
// If a required external import is not present in the vendor directory, this function
// will treat that as an error.
func sourceWalk(mainPath string) (*ApplicationSource, error) { return nil, nil }
