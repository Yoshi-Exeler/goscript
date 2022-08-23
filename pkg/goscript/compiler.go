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
	return nil, nil
}
