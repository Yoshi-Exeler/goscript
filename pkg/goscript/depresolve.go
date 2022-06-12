package goscript

import (
	"fmt"
	"os"
	"strings"
	"time"
)

func resolveDeps(mainPath string, vendorIndex map[string]*ModuleSource, localIndex map[string]*ModuleSource, standardIndex map[string]*ModuleSource) ([]*ModuleSource, error) {
	start := time.Now()
	fmt.Println("[GSC][resolveDeps] begin dependency graph resolution")
	// read the main app file into memory
	mainContent, err := os.ReadFile(mainPath)
	if err != nil {
		return nil, fmt.Errorf("could not read main application file with error %v", err)
	}
	// stip comments
	mainContent = []byte(stripComments(string(mainContent)))
	// get the direct imports from main
	directImports, err := getImportsFromSourceText(string(mainContent))
	if err != nil {
		return nil, fmt.Errorf("could not get imports from main application with error %v", err)
	}
	fmt.Printf("[GSC][resolveDeps] main file has %v direct dependencies\n", len(directImports))
	// resolve our direct dependencies
	dependencies := []*ModuleSource{}
	for _, directImport := range directImports {
		// get the direct module reference
		module, err := findModuleInIndices(directImport, vendorIndex, localIndex, standardIndex)
		if err != nil {
			return nil, err
		}
		// resolve transitive dependencies
		subModules, err := resolveUntilCompletion(module, vendorIndex, localIndex, standardIndex)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, subModules...)
		fmt.Printf("[GSC][resolveDeps] module %v adds %v transitive dependencies\n", module.ImportPath, len(subModules)-1)
	}
	fmt.Printf("[GSC][STAGE_COMPLETED] resolveDeps completed in %s\n", time.Since(start))
	for _, dep := range dependencies {
		fmt.Printf("[GSC][resolveDeps]    %v\n", dep.ImportPath)
	}
	return dependencies, nil
}

func resolveUntilCompletion(module *ModuleSource, vendorIndex map[string]*ModuleSource, localIndex map[string]*ModuleSource, standardIndex map[string]*ModuleSource) ([]*ModuleSource, error) {
	// combine all the sources of this module
	moduleBlob := ""
	for _, sourceFile := range module.Files {
		moduleBlob += sourceFile.Content
	}
	// get the imports of the current module
	imports, err := getImportsFromSourceText(moduleBlob)
	if err != nil {
		return nil, fmt.Errorf("failed to get imports from module %v", module.ImportPath)
	}
	// recursively resolve their dependencies
	modules := []*ModuleSource{}
	for _, imp := range imports {
		// grab this module from our indices
		module, err := findModuleInIndices(imp, vendorIndex, localIndex, standardIndex)
		if err != nil {
			return nil, err
		}
		subModules, err := resolveUntilCompletion(module, vendorIndex, localIndex, standardIndex)
		if err != nil {
			return nil, err
		}
		modules = append(modules, subModules...)
		// save this import in the module
		module.Imports[imp.Alias] = imp
	}
	modules = append(modules, module)
	return modules, nil
}

func getImportsFromSourceText(source string) ([]*ImportDirective, error) {
	ret := []*ImportDirective{}
	matches := IMPORT_REGEX.FindAllStringSubmatch(source, -1)
	for _, match := range matches {
		if len(match) != 2 {
			return nil, fmt.Errorf("invalid match for import statement")
		}
		parts := strings.SplitN(match[1], "/", 2)
		fullparts := strings.Split(match[1], "/")
		rt := ModuleSourceRootType(0)
		switch parts[0] {
		case "loc":
			rt = LOCAL
		case "std":
			rt = STANDARD
		case "ext":
			rt = VENDOR
		default:
			return nil, fmt.Errorf("unrecognized root type %v in import %v", parts[0], match[1])
		}
		ret = append(ret, &ImportDirective{
			RootType:     rt,
			Name:         fullparts[len(fullparts)-1],
			RawDirective: match[1],
			ImportPath:   match[1],
			Alias:        fullparts[len(fullparts)-1],
		})
	}
	aliasMatches := IMPORT_ALIAS_REGEX.FindAllStringSubmatch(source, -1)
	for _, match := range aliasMatches {
		if len(match) != 2 {
			return nil, fmt.Errorf("invalid match for import statement")
		}
		parts := strings.SplitN(match[1], "/", 2)
		fullparts := strings.Split(match[1], "/")
		rt := ModuleSourceRootType(0)
		switch parts[0] {
		case "loc":
			rt = LOCAL
		case "std":
			rt = STANDARD
		case "ext":
			rt = VENDOR
		default:
			return nil, fmt.Errorf("unrecognized root type %v in import %v", parts[0], match[1])
		}
		ret = append(ret, &ImportDirective{
			RootType:     rt,
			Name:         fullparts[len(fullparts)-1],
			RawDirective: match[1],
			ImportPath:   parts[1],
			Alias:        match[2],
		})
	}
	return ret, nil
}

// GetRequiredExternals returns the list of all external modules required by the application at main path
func getRequiredExternals(mainPath string) ([]*ExternalModuleSource, error) {
	// read the main file into memory
	content, err := os.ReadFile(mainPath)
	if err != nil {
		return nil, fmt.Errorf("could not read main application file from %v with error %v", mainPath, err)
	}
	// stip comments
	content = []byte(stripComments(string(content)))
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
		branchMatch := BRANCH_TAG_REGEX.FindStringSubmatch(match[2])
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
