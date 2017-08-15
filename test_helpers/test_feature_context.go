package testHelpers

// nolint gocyclo
import (
	"strings"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

// TestFeatureContext defines the feature context for features/test
func TestFeatureContext(s *godog.Suite) {

	s.Step(`^I eventually see the following snippets:$`, func(table *gherkin.DataTable) error {
		return childCmdPlus.WaitForCondition(func(_, output string) bool {
			success := true
			for _, row := range table.Rows {
				if !strings.Contains(output, row.Cells[0].Value) {
					success = false
					break
				}
			}
			return success
		}, time.Minute*2)
	})

}
