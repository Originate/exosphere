package helpers

// nolint gocyclo
import (
	"os"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

// TestFeatureContext defines the feature context for features/test
func TestFeatureContext(s *godog.Suite) {

	s.Step(`^I eventually see the following snippets:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows {
			if err := childCmdPlus.WaitForText(row.Cells[0].Value, time.Minute*2); err != nil {
				return err
			}
		}
		return nil
	})

	s.Step(`^I send an interrupt signal$`, func() error {
		return childCmdPlus.Cmd.Process.Signal(os.Interrupt)
	})

}
