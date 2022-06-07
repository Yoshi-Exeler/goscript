package goscript

import (
	"fmt"
	"regexp"
)

// this package will contain the preprocessor for goscript

var STRUCT_NAME_REGEX = regexp.MustCompile(`(?m)struct (.*) {`)
var FUNC_NAME_REGEX = regexp.MustCompile(`(?m)func (.*)\(`)

// this regex matches any access to an external properties (for example db.Connect() will match with CG1 being the module name and CG2 being the property)
var EXTERNAL_SYMBOL_REGEX = regexp.MustCompile(`(?m)[^a-zA-Z0-9]{1}((?:[a-zA-Z]{1}[a-zA-Z0-9]?)*)\.((?:[a-zA-Z]{1}[a-zA-Z0-9]?)*)`)

var SIMPLE_STRING_REGEX = regexp.MustCompile(`(?mU)".*[^\\]"`)

var MULTILINE_STRING_REGEX = regexp.MustCompile(`(?Us)\x60.*\x60`)

/* generateFQSC will apply preprocessing steps to the applications source code
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
	fmt.Println("[GSC][genFQSC] begin generation of fqsc, stripping directives and generating module blobs")
	source.ApplicationFile.Content = stipDirectives(source.ApplicationFile.Content)
	for _, mod := range source.Modules {
		for _, file := range mod.Files {
			mod.Content += file.Content
		}
		mod.Content = stipDirectives(mod.Content)
		mod.Content = prefixSymbolNames(mod.Content, mod.Hash, mod.Name)
		fqsc, err := fixReferences(mod.Content, mod, source.Modules)
		if err != nil {
			return "", fmt.Errorf("failed to fix references for module %v with error %v", mod.ImportPath, err)
		}
		mod.Content = fqsc
		fmt.Printf("[GSC][genFQSC] module=%v fqsc:\n%v\n", mod.ImportPath, mod.Content)
	}
	fqsc, err := fixReferences(source.ApplicationFile.Content, nil, source.Modules)
	if err != nil {
		return "", fmt.Errorf("failed to fix references for main file with error %v", err)
	}
	source.ApplicationFile.Content = fqsc
	fmt.Println("[GSC][genFQSC] module blobs generated")
	fmt.Println(source.ApplicationFile.Content)
	return "", nil
}

var APPLICATION_REGEX = regexp.MustCompile(`(?m)application (.*)$`)
var EMPTY_LINE_REGEX = regexp.MustCompile(`(?m)^\n`)

func fixReferences(source string, module *ModuleSource, modules []*ModuleSource) (string, error) {
	// get the string mask for the source code
	sourceMask := getStringMask(source)
	// get both index and full matches of symbols in the source code
	symbolIndexMatches := EXTERNAL_SYMBOL_REGEX.FindAllStringIndex(source, -1)
	fullMatches := EXTERNAL_DIRECTIVE_REGEX.FindAllStringSubmatch(source, -1)
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
		newSymbol := fmt.Sprintf("%v_%v_%v.%v", targetModule.Hash, targetModule.Name, fullMatches[i][1], fullMatches[i][2])
		for j := symbolIndexMatches[i][0]; j < symbolIndexMatches[i][1]; j++ {
			// replace the relevant section in the fqsc
			source = source[:symbolIndexMatches[i][0]] + newSymbol + source[symbolIndexMatches[i][1]:]
		}
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
	stepOne := STRUCT_NAME_REGEX.ReplaceAllString(source, fmt.Sprintf("struct %v_%v_$1 {", prefix, name))
	stepTwo := FUNC_NAME_REGEX.ReplaceAllString(stepOne, fmt.Sprintf("func %v_%v_$1(", prefix, name))
	return stepTwo
}

func stipDirectives(source string) string {
	stepOne := IMPORT_REGEX.ReplaceAllString(source, "")
	stepTwo := APPLICATION_REGEX.ReplaceAllString(stepOne, "")
	stepThree := EXTERNAL_DIRECTIVE_REGEX.ReplaceAllString(stepTwo, "")
	stepFour := EMPTY_LINE_REGEX.ReplaceAllString(stepThree, "")
	return stepFour
}
