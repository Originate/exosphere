package steps;


import cucumber.api.PendingException;
import cucumber.api.java.After;
import cucumber.api.java.Before;
import cucumber.api.java.en.*;
import static org.junit.Assert.*;

import events.ExorelayError;
import exocomMock.ExocomMock;
import exorelay.Exorelay;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.junit.Rule;
import org.junit.rules.ExpectedException;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

public class Steps {

    private ExocomMock exocom;
    private Exorelay exorelay;
    private String statusCode;

    private final ByteArrayOutputStream outContent = new ByteArrayOutputStream();
    private final ByteArrayOutputStream errContent = new ByteArrayOutputStream();

    @Before
    public void setUpStreams() {
//        System.setOut(new PrintStream(outContent));
//        System.setErr(new PrintStream(errContent));
    }



    /* -------------- ExoCom -------------- */

    @Given("^ExoCom runs at port (\\d+)$")
    public void exocomRunsAtPort(int port) {
        exocom = new ExocomMock(port);
        new Thread(exocom).start();
    }

    /* -------------- ExoRelay -------------- */

    @Given("^an ExoRelay instance listening on port (\\d+)$")
    public void anExoRelayInstanceListeningOnPort(int port) throws ExorelayError {
        exocom.registerService("test-service", port);
        Map<String, Object> config = new HashMap<>();
        config.put("exocom-port", exocom.getPort());
        config.put("service-name", "test-service");
        exorelay = new Exorelay(config);
        exorelay.listen(port);
    }

    @Given("^an ExoRelay instance$")
    public void anExoRelayInstance() throws Throwable {
        Map<String, Object> config = new HashMap<>();
        config.put("exocom-port", exocom.getPort());
        config.put("service-name", "test-service");
        exorelay = new Exorelay(config);
    }

    @And("^an ExoRelay instance running inside the \"([^\"]*)\" service at port (\\d+)$")
    public void anExoRelayInstanceCalledRunningInsideTheServiceAtPort(String service, int port) throws ExorelayError {
        exocom.registerService(service, port);
        Map<String, Object> config = new HashMap<>();
        config.put("exocom-port", exocom.getPort());
        config.put("service-name", service);
        exorelay = new Exorelay(config);
        exorelay.listen(port);
    }

    /* -------------- status.feature -------------- */

    @When("^I check the status$")
    public void iCheckTheStatus() throws InterruptedException{
        JSONObject message = new JSONObject();
        message.put("service", "test-service");
        message.put("name", "__status");
        try {
            exocom.send(message);
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }
        TimeUnit.MILLISECONDS.sleep(50);
        assertTrue(exocom.getReceievedMessages().size() > 0);
        statusCode = (String) exocom.getReceievedMessages().get(0).get("name");
    }

    @Then("^it signals it is online$")
    public void itSignalsItIsOnline() {
        assertEquals(statusCode, "__status-ok");

    }

    /* -------------- listen.feature -------------- */


    @When("^I take it online at port (\\d+)$")
    public void iTakeItOnlineAtPort(int port) {
        exorelay.listen(port);
    }

    @Then("^it is online at port (\\d+)$")
    public void itIsOnlineAtPort(int port) throws InterruptedException {
        exocom.registerService("test-service", port);
        JSONObject message = new JSONObject();
        message.put("service", "test-service");
        message.put("name", "__status");
        try {
            exocom.send(message);
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }
        TimeUnit.MILLISECONDS.sleep(50);
        assertTrue(exocom.getReceievedMessages().size() > 0);
        statusCode = (String) exocom.getReceievedMessages().get(0).get("name");
        assertEquals(statusCode, "__status-ok");
    }

    /* -------------- exocom-port.feature -------------- */
    private String thrownMessage;

    @When("^I try to create an ExoRelay without providing the ExoCom port$")
    public void iTryToCreateAnExoRelayWithoutProvidingTheExoComPort() {
        Map<String, Object> config = new HashMap<>();
        config.put("service-name", "test-service");
        try {
            exorelay = new Exorelay(config);
        } catch (ExorelayError exorelayError) {
            thrownMessage = exorelayError.getMessage();
        }
    }


    @Then("^it throws the error \"([^\"]*)\"$")
    public void itThrowsTheError(String errorMessage) {
        assertEquals(errorMessage, thrownMessage);
    }

    @When("^I create an ExoRelay instance that uses ExoCom port (\\d+)$")
    public void iCreateAnExoRelayInstanceThatUsesExoComPort(int port) throws Throwable {
        Map<String, Object> config = new HashMap<>();
        config.put("exocom-port", port);
        config.put("service-name", "test-service");
        exorelay = new Exorelay(config);
    }

    @Then("^this instance uses the ExoCom port (\\d+)$")
    public void thisInstanceusesTheExoComPort(int port) throws Throwable {
        assertEquals(exorelay.getExocomPort(), new Integer(port));
    }

    /* -------------- sending.feature -------------- */

    private String messageId;
    @When("^sending the message:$")
    public void sendingTheMessage(String msg) throws ParseException, ExorelayError {
        JSONParser parser = new JSONParser();
        JSONObject root = (JSONObject) parser.parse(msg);
        String name = (String) root.get("name");
        String payload = null;
        if(root.get("payload") != null) {
            payload = ((JSONObject) root.get("payload")).toJSONString();
        }
        messageId = exorelay.send(name, payload);
    }

    @Then("^ExoRelay makes the ZMQ request:$")
    public void exorelayMakesTheZMQRequest(String msg) throws ParseException, InterruptedException {
        TimeUnit.MILLISECONDS.sleep(10);
        assertTrue(exocom.getReceievedMessages().size() > 0);
        String receivedRequest = exocom.getReceievedMessages().get(0).toJSONString();
        JSONParser parser = new JSONParser();
        JSONObject sentRequest = (JSONObject) parser.parse(msg);
        sentRequest.put("id", messageId);
        TimeUnit.MILLISECONDS.sleep(50);
        assertEquals(sentRequest.toJSONString(), receivedRequest);
    }

    @When("^trying to send an empty message:$")
    public void tryingToSendAnEmptyMessage(String msg) throws ParseException {
        JSONParser parser = new JSONParser();
        JSONObject root = (JSONObject) parser.parse(msg);
        String name = (String) root.get("name");

        try {
            messageId = exorelay.send(name, null);
        } catch (ExorelayError e) {
            thrownMessage = e.getMessage();
        }

    }

    @Rule
    private ExpectedException expectedException = ExpectedException.none();

    @Then("^ExoRelay throws an ExorelayError exception with the message \"([^\"]*)\"$")
    public void exorelayEmitsAnEventWithTheError(String message) {
        expectedException.expect(ExorelayError.class);

        assertEquals(thrownMessage, message);
    }


    @Given("^I register a handler for the \"([^\"]*)\" message:$")
    public void iRegisterAHandlerForTheMessage(String messageName, String code) throws ExorelayError {
        exorelay.registerHandler(messageName, param -> {System.out.println("Hello world!");});
    }

    @When("^receiving this message:$")
    public void receivingThisMessage(String message) throws Exception {
        JSONParser parser = new JSONParser();
        JSONObject data = (JSONObject) parser.parse(message);
        data.put("service", "test-service");

        System.out.println("receiving this message");
        exocom.send(data);
    }

    @Then("^ExoRelay runs the registered handler, in this example calling print with \"([^\"]*)\"$")
    public void exorelayRunsTheRegisteredHandlerInThisExampleCallingWith(String str) {
        assertEquals(str, outContent.toString());
        //TODO: setUpStreamns() and cleanUpStreams() commented out for debugging
        //un-comment them to make this step definition accurate
    }

    @Then("^ExoRelay emits an \"([^\"]*)\" event with the error \"([^\"]*)\"$")
    public void exorelayEmitsAnEventWithTheError(String arg0, String arg1) throws Throwable {
        // Write code here that turns the phrase above into concrete actions
        throw new PendingException();
    }


    @After
    public void closePorts() {
        if(exorelay != null) {
            exorelay.closePort();
            exorelay.close();
        }
        if(exocom != null) exocom.close();
    }


    @After
    public void cleanUpStreams() {
//        outContent.reset();
//        errContent.reset();
//        System.setOut(null);
//        System.setErr(null);
    }
}
