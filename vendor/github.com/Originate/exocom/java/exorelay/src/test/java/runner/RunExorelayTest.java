package runner;

import cucumber.api.CucumberOptions;
import cucumber.api.junit.Cucumber;
import org.junit.runner.RunWith;

@RunWith(Cucumber.class)
@CucumberOptions(
        features = {"src/test/resourcese"},
        glue = {"steps"},
        format={"pretty"}
)
public class RunExorelayTest {
}
