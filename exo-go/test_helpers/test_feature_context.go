package testHelpers

// nolint gocyclo
import (
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

func TestFeatureContext(s *godog.Suite) {

	s.Step(`^I eventually see the following snippets:$`, func(table *gherkin.DataTable) error {
		for _, row := range table.Rows {
			if err := childCmdPlus.WaitForText(row.Cells[0].Value, time.Minute*2); err != nil {
				return err
			}
		}
		return nil
	})

}
