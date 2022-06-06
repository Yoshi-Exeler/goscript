package goscript

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

// GetRequiredExternals returns the list of all external modules required by the application at main path
func getRequiredExternals(mainPath string) (*[]ExternalModuleSource, error) {

}

// SourceWalk will parse the imports required by the application at the specified path
// If a required external import is not present in the vendor directory, this function
// will treat that as an error.
func sourceWalk(mainPath string) (*ApplicationSource, error) { return nil, nil }
