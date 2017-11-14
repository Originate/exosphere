package matchers

import (
	"fmt"

	"github.com/Originate/exosphere/src/types/hcl"
	"github.com/onsi/gomega/types"
)

type haveHCLVariableMatcher struct {
	name string
}

func HaveHCLVariable(name string) types.GomegaMatcher {
	return &haveHCLVariableMatcher{name: name}
}

func (m *haveHCLVariableMatcher) Match(actual interface{}) (bool, error) {
	hclFile, ok := actual.(*hcl.File)
	if !ok {
		return false, fmt.Errorf("HaveHCLVariable matcher expects an HCL File")
	}
	_, hasKey := hclFile.Variable[m.name]
	return hasKey, nil
}

func (m *haveHCLVariableMatcher) FailureMessage(actual interface{}) string {
	variables := actual.(hcl.File).GetVariableNames()
	return fmt.Sprintf("Expected hcl file to contain the variable %q, but it only has:\n%q", m.name, variables)
}

func (m *haveHCLVariableMatcher) NegatedFailureMessage(actual interface{}) string {
	variables := actual.(hcl.File).GetVariableNames()
	return fmt.Sprintf("Expected hcl file to not contain the variable %q, but it does:\n%q", m.name, variables)
}
