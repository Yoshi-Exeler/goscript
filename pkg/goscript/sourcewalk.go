package goscript

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type ModuleSourceRootType byte

const (
	VENDOR   ModuleSourceRootType = 1 // modules that have this type will be searched in the vendor directory
	STANDARD ModuleSourceRootType = 2 // modules that have this type will be searched in the standard library
	LOCAL    ModuleSourceRootType = 3 // modules that have this type will be searched locally in your project
)

type ModuleSource struct {
	Name       string                      // name of the module
	Path       string                      // actual path of the module relative to its root
	ImportPath string                      // the path using which this module can be imported (without its prefix)
	RootType   ModuleSourceRootType        // type of the root this module can be found under
	Files      []SourceFile                // paths to all the files in the module
	Hash       string                      // the hash uniquely identifying this module
	Content    string                      // merged content of all files (happens in preprocessor)
	Imports    map[string]*ImportDirective // map of the aliased module name to the module import
}

type ImportDirective struct {
	RootType     ModuleSourceRootType // which root path the module should be imported from
	Name         string               // name of the module to be imported
	RawDirective string               // the unprocessed directive
	ImportPath   string               // the import path of the module without the root prefix
	Alias        string               // the alias in the import directive
}

type SourceFile struct {
	Path    string                      // path to the source file
	Content string                      // content of the source file
	Imports map[string]*ImportDirective // aliased imports only set in the main file
}

type ApplicationSource struct {
	ApplicationFile SourceFile      // path to the main file
	Modules         []*ModuleSource // list of local and standard libarary modules
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
var BRANCH_TAG_REGEX = regexp.MustCompile(`(?m)@branch=([0-9a-zA-Z\.]*)`)

// SourceWalk will discover source files and parse the imports required by the application
// at the specified path  If a required external import is not present in the vendor directory,
// this function will treat that as an error.
func SourceWalk(mainPath string, workspace string) (*ApplicationSource, error) {
	start := time.Now()
	fmt.Println("[GSC][sourceWalk] begin sourcewalk")
	// get the external modules required by our app
	dependencies, err := getRequiredExternals(mainPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get required external modules with error %v", err)
	}
	fmt.Printf("[GSC][sourceWalk] application requires %v direct external dependencies\n", len(dependencies))
	// TODO: call external dependency resolver here
	fmt.Printf("[GSC][sourceWalk] begin indexing. local=%v vendor=%v standard=%v\n", workspace, VENDORPATH, STDPATH)
	// index the vendor directory
	vendorIndex, err := indexModuleCollection(VENDORPATH, "ext")
	if err != nil {
		return nil, fmt.Errorf("failed to index vendor directory with error %v", err)
	}
	// index the local directory (mainPath)
	localIndex, err := indexModuleCollection(workspace, "loc")
	if err != nil {
		return nil, fmt.Errorf("failed to index local directory with error %v", err)
	}
	// index the standard directory
	standardIndex, err := indexModuleCollection(STDPATH, "std")
	if err != nil {
		return nil, fmt.Errorf("failed to index standard directory with error %v", err)
	}
	fmt.Printf("[GSC][sourceWalk] finished indexing. local=%v vendor=%v standard=%v\n", len(localIndex), len(vendorIndex), len(standardIndex))
	// now ensure all packages that the application uses actually exist locally
	for _, dependency := range dependencies {
		if vendorIndex["ext/"+dependency.Name] == nil {
			return nil, fmt.Errorf("required external module %v not found", dependency.Name)
		}
	}
	// resolve our dependencies
	flatDeps, err := findMinimalResolution(mainPath, vendorIndex, localIndex, standardIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies with error %v", err)
	}
	fmt.Println("[GSC][sourceWalk] dependency graph resolution completed.")
	for _, dep := range flatDeps {
		fmt.Printf("[GSC][sourceWalk]    %v\n", dep.ImportPath)
	}
	// grab the main file
	mainContent, err := os.ReadFile(mainPath)
	if err != nil {
		return nil, fmt.Errorf("could not read main application file with error %v", err)
	}
	// stip comments
	mainContent = []byte(stripComments(string(mainContent)))
	// declare app source
	src := &ApplicationSource{
		ApplicationFile: SourceFile{
			Path:    mainPath,
			Content: string(mainContent),
			Imports: make(map[string]*ImportDirective),
		},
		Modules: flatDeps,
	}
	// add import information to the main file
	mainImports, err := getImportsFromSourceText(string(mainContent))
	if err != nil {
		return nil, fmt.Errorf("could not get main file imports with error %v", err)
	}
	for _, imp := range mainImports {
		alloc := *imp
		src.ApplicationFile.Imports[imp.Alias] = &alloc
	}
	fmt.Printf("[GSC][STAGE_COMPLETION] sourcewalk completed in %s\n", time.Since(start))
	// return the app source struct
	return src, nil
}

var IMPORT_REGEX = regexp.MustCompile(`import "(.*)"\s*\n`)
var IMPORT_ALIAS_REGEX = regexp.MustCompile(`import "(.*)" as (.*)`)

func findModuleInIndices(imp *ImportDirective, vendorIndex map[string]*ModuleSource, localIndex map[string]*ModuleSource, standardIndex map[string]*ModuleSource) (*ModuleSource, error) {
	switch imp.RootType {
	case LOCAL:
		mod := localIndex[imp.ImportPath]
		if mod == nil {
			return nil, fmt.Errorf("module %v not found in local workspace", imp.ImportPath)
		}
		return mod, nil
	case STANDARD:
		mod := standardIndex[imp.ImportPath]
		if mod == nil {
			return nil, fmt.Errorf("module %v not found in standard library", imp.ImportPath)
		}
		return mod, nil
	case VENDOR:
		mod := vendorIndex[imp.ImportPath]
		if mod == nil {
			return nil, fmt.Errorf("module %v not found in vendor directory", imp.ImportPath)
		}
		return mod, nil
	default:
		return nil, fmt.Errorf("could not resolve import directive %v, unrecognized root directory", imp)
	}
}

// indexModuleCollection will create an index of a module collection directory such as $VENDORPATH or the local path
func indexModuleCollection(path string, relativeRoot string) (map[string]*ModuleSource, error) {
	// index the directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read directory with error %v", err)
	}
	// create a collection for vendor modules
	modules := make(map[string]*ModuleSource)
	// filter the entries
	for _, entry := range entries {
		// skip everything that isnt a directory
		if !entry.IsDir() {
			continue
		}
		// add this module and all its submodules to the vendor module collection
		err = recAddModule(filepath.Join(path, entry.Name()), relativeRoot+"/"+entry.Name(), entry.Name(), modules)
		if err != nil {
			return nil, fmt.Errorf("failed to index directory %v", err)
		}
	}
	return modules, nil
}

// recAddModule will add the module at the specified path and all its submodules to the out map
func recAddModule(path string, importPath string, name string, out map[string]*ModuleSource) error {
	this := ModuleSource{
		Name:       name,
		Path:       path,
		ImportPath: importPath,
		Imports:    make(map[string]*ImportDirective),
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
			recAddModule(filepath.Join(path, entry.Name()), importPath+"/"+entry.Name(), entry.Name(), out)
			continue
		}
		// read its content
		content, err := os.ReadFile(filepath.Join(path, entry.Name()))
		if err != nil {
			return fmt.Errorf("cannot read source file %v with error %v", filepath.Join(path, entry.Name()), err)
		}
		// stip comments
		content = []byte(stripComments(string(content)))
		// otherwise, we add this entry to the file list
		this.Files = append(this.Files, SourceFile{
			Path:    filepath.Join(path, entry.Name()),
			Content: string(content),
		})
		// and write its name and path into the hash
		hash.Write([]byte(filepath.Join(path, entry.Name())))
	}
	// set the module hash
	this.Hash = hex.EncodeToString(hash.Sum(nil))
	// append the entry to the output list
	out[importPath] = &this
	return nil
}

var COMMENT_REGEX = regexp.MustCompile(`(?m)^//.*`)

func stripComments(soure string) string {
	return COMMENT_REGEX.ReplaceAllString(soure, "")
}
