package goscript

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// this package will contain the preprocessor for goscript

var STRUCT_NAME_REGEX = regexp.MustCompile(`(?m)struct (.*) {`)
var FUNC_NAME_REGEX = regexp.MustCompile(`(?m)func (.*)\(`)

// this regex matches any access to an external properties (for example db.Connect() will match with CG1 being the module name and CG2 being the property)
var EXTERNAL_SYMBOL_REGEX = regexp.MustCompile(`(?m)((?:[a-zA-Z]{1}[a-zA-Z0-9]?)*)\.((?:[a-zA-Z]{1}[a-zA-Z0-9]?)*)`)

var SIMPLE_STRING_REGEX = regexp.MustCompile(`(?mU)".*[^\\]"`)

var MULTILINE_STRING_REGEX = regexp.MustCompile(`(?Us)\x60.*\x60`)

/*
	 generateFQSC will apply preprocessing steps to the applications source code
	   The following steps will be applied:
		- delete import directives
		- delete application directive
		- delete external directives
		- prefix all symbols with the hash of their module (a function getJWT from the package jwt would become hash_getJWT)
		- replace external symbol references with their expected name (jwt.getJWT becomes hash_getJWT)
			- dots can be contained in property access, member access, numbers and strings
			- numbers shouldnt be a problem since we can just ignore them based on required property names
			- we need to detect all strings, remember their starting and ending indices in the string and then bounds check wether or not the
			  match we have is inside a string zone
		- merge all source files
*/
func generateFQSC(source *ApplicationSource) (string, error) {
	start := time.Now()
	fmt.Println("[GSC][genFQSC] begin generation of fqsc, stripping directives and generating module blobs")
	source.ApplicationFile.Content = stripDirectives(source.ApplicationFile.Content)
	for _, mod := range source.Modules {
		for _, file := range mod.Files {
			mod.Content += file.Content
		}
		mod.Content = stripDirectives(mod.Content)
		mod.Content = prefixSymbolNames(mod.Content, mod.Hash, mod.Name)
		fqsc, err := fixReferences(mod.Content, mod, source.Modules)
		if err != nil {
			return "", fmt.Errorf("failed to fix references for module %v with error %v", mod.ImportPath, err)
		}
		mod.Content = fqsc
		fmt.Printf("[GSC][genFQSC] completed module %v\n", mod.ImportPath)
	}
	fmt.Println("[GSC][genFQSC] modules stripped, module symbols normalized")
	fqsc, err := fixMainFileReferences(source.ApplicationFile.Content, &source.ApplicationFile, source.Modules)
	if err != nil {
		return "", fmt.Errorf("failed to fix references for main file with error %v", err)
	}
	source.ApplicationFile.Content = fqsc
	fmt.Println("[GSC][genFQSC] main stripped, main symbols normalized, merging now")
	// now we just merge all the sources and return them
	fullFQSC := source.ApplicationFile.Content
	for _, mod := range source.Modules {
		fullFQSC += "\n" + mod.Content
	}
	fmt.Println("[GSC][genFQSC] blobs merged into FQSC successfully")
	fmt.Printf("[GSC][STAGE_COMPLETION] fqsc generation completed in %s\n", time.Since(start))
	fullFQSC += "\n>"
	// dump fqsc to a file if FQSC debugging is enabled
	if DEBUG_DUMP_FQSC {
		os.WriteFile("debug_dump.fqsc", []byte(fullFQSC), 0600)
	}
	return fullFQSC, nil
}

var APPLICATION_REGEX = regexp.MustCompile(`(?m)application (.*)$`)
var EMPTY_LINE_REGEX = regexp.MustCompile(`(?m)^\n`)

func trimWhitespace(code string) string {
	lines := strings.Split(code, "\n")
	ret := bytes.NewBuffer([]byte{})
	for idx, line := range lines {
		ret.Write([]byte(strings.TrimSpace(line)))
		if idx != len(lines)-1 {
			ret.Write([]byte("\n"))
		}
	}
	return ret.String()
}

func fixReferences(source string, module *ModuleSource, modules []*ModuleSource) (string, error) {
	// get the string mask for the source code
	sourceMask := getStringMask(source)
	// get both index and full matches of symbols in the source code
	symbolIndexMatches := EXTERNAL_SYMBOL_REGEX.FindAllStringIndex(source, -1)
	fullMatches := EXTERNAL_SYMBOL_REGEX.FindAllStringSubmatch(source, -1)
	delta := 0
	for i := 0; i < len(symbolIndexMatches); i++ {
		// skip matches that are inside of a string
		if sourceMask[symbolIndexMatches[i][0]] {
			continue
		}
		// get the referenced import
		symbolImport := module.Imports[fullMatches[i][1]]
		if symbolImport == nil {
			return "", fmt.Errorf("symbol %v used but not imported", fullMatches[i][1])
		}
		// find the imported module in the module collection
		targetModule := &ModuleSource{}
		for _, mod := range modules {
			if mod.ImportPath == symbolImport.ImportPath {
				targetModule = mod
				break
			}
		}
		if targetModule == nil {
			return "", fmt.Errorf("imported module %v not found in local modules", fullMatches[i][1])
		}
		// generate the new symbol to replace the current one
		newSymbol := fmt.Sprintf("fn_%v_%v_%v", targetModule.Hash, targetModule.Name, fullMatches[i][2])
		source = source[:symbolIndexMatches[i][0]+delta] + newSymbol + source[symbolIndexMatches[i][1]+delta:]
		delta += len(newSymbol) - (symbolIndexMatches[i][1] - symbolIndexMatches[i][0])
	}
	return source, nil
}

func fixMainFileReferences(source string, file *SourceFile, modules []*ModuleSource) (string, error) {
	// get the string mask for the source code
	sourceMask := getStringMask(source)
	// get both index and full matches of symbols in the source code
	symbolIndexMatches := EXTERNAL_SYMBOL_REGEX.FindAllStringIndex(source, -1)
	fullMatches := EXTERNAL_SYMBOL_REGEX.FindAllStringSubmatch(source, -1)
	delta := 0
	for i := 0; i < len(symbolIndexMatches); i++ {
		// skip matches that are inside of a string
		if sourceMask[symbolIndexMatches[i][0]] {
			continue
		}
		// get the referenced import
		symbolImport := file.Imports[fullMatches[i][1]]
		if symbolImport == nil {
			return "", fmt.Errorf("symbol %v used but not imported", fullMatches[i][1])
		}
		// find the imported module in the module collection
		var targetModule *ModuleSource
		for _, mod := range modules {
			alloc := *mod
			if mod.ImportPath == symbolImport.ImportPath {
				targetModule = &alloc
				break
			}
		}
		if targetModule == nil {
			return "", fmt.Errorf("imported module %v not found in local modules", fullMatches[i][1])
		}
		// generate the new symbol to replace the current one
		newSymbol := fmt.Sprintf("fn_%v_%v_%v", targetModule.Hash, targetModule.Name, fullMatches[i][2])
		source = source[:symbolIndexMatches[i][0]+delta] + newSymbol + source[symbolIndexMatches[i][1]+delta:]
		delta += len(newSymbol) - (symbolIndexMatches[i][1] - symbolIndexMatches[i][0])
	}
	return source, nil
}

func getStringMask(source string) []bool {
	mask := make([]bool, len(source))
	simpleStringMatches := SIMPLE_STRING_REGEX.FindAllStringIndex(source, -1)
	for _, match := range simpleStringMatches {
		for i := match[0]; i < match[1]; i++ {
			mask[i] = true
		}
	}
	multilineStringRegex := MULTILINE_STRING_REGEX.FindAllStringIndex(source, -1)
	for _, match := range multilineStringRegex {
		for i := match[0]; i < match[1]; i++ {
			mask[i] = true
		}
	}
	return mask
}

func prefixSymbolNames(source string, prefix string, name string) string {
	// find all struct definitions
	stepOne := STRUCT_NAME_REGEX.ReplaceAllString(source, fmt.Sprintf("struct st_%v_%v_$1 {", prefix, name))
	stepTwo := FUNC_NAME_REGEX.ReplaceAllString(stepOne, fmt.Sprintf(">\nfunc fn_%v_%v_$1(", prefix, name))
	return stepTwo
}

func stripDirectives(source string) string {
	stepOne := IMPORT_REGEX.ReplaceAllString(source, "")
	stepTwo := APPLICATION_REGEX.ReplaceAllString(stepOne, "")
	stepThree := EXTERNAL_DIRECTIVE_REGEX.ReplaceAllString(stepTwo, "")
	stepFour := EMPTY_LINE_REGEX.ReplaceAllString(stepThree, "")
	return stepFour
}
