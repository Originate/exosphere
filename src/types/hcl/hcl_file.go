package hcl

import hashicorpHCL "github.com/hashicorp/hcl"

// File represents a parsed HCL Terraform file
type File struct {
	Terraform Terraform
	Provider  map[string]Provider
	Variable  map[string]Variable
	Module    map[string]Module
}

// GetVariableNames returns the list of variable names defined in the file
func (f *File) GetVariableNames() []string {
	keys := []string{}
	for v := range f.Variable {
		keys = append(keys, v)
	}
	return keys
}

// GetModuleNames returns the list of module names defined in the file
func (f *File) GetModuleNames() []string {
	keys := []string{}
	for m := range f.Module {
		keys = append(keys, m)
	}
	return keys
}

// GetHCLFileFromTerraform will return a File for the passed string of hcl terraform
func GetHCLFileFromTerraform(terraform string) (*File, error) {
	var hclFile File
	err := hashicorpHCL.Decode(&hclFile, terraform)
	return &hclFile, err
}
