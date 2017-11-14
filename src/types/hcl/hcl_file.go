package hcl

import hashicorpHCL "github.com/hashicorp/hcl"

type File struct {
	Terraform Terraform
	Provider  map[string]Provider
	Variable  map[string]Variable
	Module    map[string]Module
}

func (f File) GetVariableNames() []string {
	keys := []string{}
	for v := range f.Variable {
		keys = append(keys, v)
	}
	return keys
}

func (f File) GetModuleNames() []string {
	keys := []string{}
	for m := range f.Module {
		keys = append(keys, m)
	}
	return keys
}

func GetHCLFileFromTerraform(terraform string) (*File, error) {
	var hclFile File
	err := hashicorpHCL.Decode(&hclFile, terraform)
	return &hclFile, err
}
