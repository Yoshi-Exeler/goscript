package goscript

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ModuleSourceRootType byte

const (
	VENDOR   ModuleSourceRootType = 1 // modules that have this type will be searched in the vendor directory
	STANDARD ModuleSourceRootType = 2 // modules that have this type will be searched in the standard library
	LOCAL    ModuleSourceRootType = 3 // modules that have this type will be searched locally in your project
)

type ModuleSource struct {
	Name       string               // name of the module
	Path       string               // actual path of the module relative to its root
	ImportPath string               // the path using which this module can be imported (without its prefix)
	RootType   ModuleSourceRootType // type of the root this module can be found under
	Files      []string             // paths to all the files in the module
	Hash       string               // the hash uniquely identifying this module
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

var EXTERNAL_DIRECTIVE_REGEX = regexp.MustCompile(`(?mU)external ([a-zA-Z]*) from "(.*)"`)
var VERSION_TAG_REGEX = regexp.MustCompile(`(?m)@version=([\^0-9\.]*)`)
var BRNACH_TAG_REGEX = regexp.MustCompile(`(?m)@branch=([0-9a-zA-Z\.]*)`)

// GetRequiredExternals returns the list of all external modules required by the application at main path
func GetRequiredExternals(mainPath string) ([]*ExternalModuleSource, error) {
	// read the main file into memory
	content, err := os.ReadFile(mainPath)
	if err != nil {
		return nil, fmt.Errorf("could not read main application file from %v with error %v", mainPath, err)
	}
	ret := []*ExternalModuleSource{}
	// regex for our external definitions
	matches := EXTERNAL_DIRECTIVE_REGEX.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) != 3 {
			return nil, fmt.Errorf("matched invalid external definition")
		}
		// if a version and or branch tag exist, process them
		version := ""
		branch := ""
		// check for a version tag
		versionMatch := VERSION_TAG_REGEX.FindStringSubmatch(match[2])
		if len(versionMatch) == 2 {
			version = versionMatch[1]
		}
		// check for a branch tag
		branchMatch := BRNACH_TAG_REGEX.FindStringSubmatch(match[2])
		if len(branchMatch) == 2 {
			branch = branchMatch[1]
		}
		// cleanup the url
		if strings.Contains(match[2], "@") {
			parts := strings.Split(match[2], "@")
			if len(parts) == 0 {
				return nil, fmt.Errorf("failed to cleanup external package url %v", match[2])
			}
			// use the base url
			match[2] = parts[0]
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

// SourceWalk will discover source files and parse the imports required by the application
// at the specified path  If a required external import is not present in the vendor directory,
// this function will treat that as an error.
func SourceWalk(mainPath string) (*ApplicationSource, error) {
	// index the vendor directory
	vendorIndex, err := indexSourceRoot(VENDORPATH, "ext")
	if err != nil {
		return nil, fmt.Errorf("failed to index vendor directory with error %v", err)
	}
	// index the local directory (mainPath)
	localIndex, err := indexSourceRoot(mainPath, "loc")
	if err != nil {
		return nil, fmt.Errorf("failed to index local directory with error %v", err)
	}
	// index the standard directory
	standardIndex, err := indexSourceRoot(STDPATH, "std")
	if err != nil {
		return nil, fmt.Errorf("failed to index standard directory with error %v", err)
	}
	fmt.Println(vendorIndex, localIndex, standardIndex)
	return nil, nil
}

func indexSourceRoot(path string, relativeRoot string) (map[string]ModuleSource, error) {
	// index the directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory with error %v", err)
	}
	// create a collection for vendor modules
	modules := make(map[string]ModuleSource)
	// filter the entries
	for _, entry := range entries {
		// skip everything that isnt a directory
		if !entry.IsDir() {
			continue
		}
		// add this module and all its submodules to the vendor module collection
		err = dfsAddModule(filepath.Join(path, entry.Name()), relativeRoot+"/"+entry.Name(), entry.Name(), modules)
		if err != nil {
			return nil, fmt.Errorf("failed to index directory %v", err)
		}
	}
	return modules, nil
}

func dfsAddModule(path string, importPath string, name string, out map[string]ModuleSource) error {
	this := ModuleSource{
		Name:       name,
		Path:       path,
		ImportPath: importPath,
	}
	// list the specified directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("cannot read module directory %v with error %v", path, err)
	}
	// prepare the module hash
	hash := fnv.New32a()
	// iterate over the local dir entries
	for _, entry := range entries {
		// if we have found a submodule recursively add it too
		if entry.IsDir() {
			dfsAddModule(filepath.Join(path, entry.Name()), importPath+"/"+entry.Name(), entry.Name(), out)
			continue
		}
		// otherwise, we add this entry to the file list
		this.Files = append(this.Files, filepath.Join(path, entry.Name()))
		// and write its name and path into the hash
		hash.Write([]byte(filepath.Join(path, entry.Name())))
	}
	// set the module hash
	this.Hash = hex.EncodeToString(hash.Sum(nil))
	// append the entry to the output list
	out[name] = this
	return nil
}
