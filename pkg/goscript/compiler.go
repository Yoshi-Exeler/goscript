package goscript

import "fmt"

type Compiler struct{}

type CompileJob struct {
	MainFilePath       string
	VendorPath         string
	LocalWorkspaceRoot string
	StandardLibPath    string
}

func (c *Compiler) Compile(job CompileJob) (*Program, error) {
	appSource, err := discoverSources(job.MainFilePath, job.LocalWorkspaceRoot)
	if err != nil {
		return nil, err
	}
	fqsc, err := generateFQSC(appSource)
	if err != nil {
		return nil, err
	}
	intermediate := parse(fqsc)
	fmt.Println(intermediate)
	return c.generateProgram(intermediate)
}

/*
generateProgram generates a program from the intermediary program representation
The following steps will be performed:
- Replace Symbol Placeholders in expressions
- Replace Function Placeholders in expressions
- Eliminate Dead code
- Optimize:
	- Resolve constant expressions as far as possible
- Generate the actual bytecode
*/
func (c *Compiler) generateProgram(intermediate *IntermediateProgram) (*Program, error) {
	return nil, nil
}
