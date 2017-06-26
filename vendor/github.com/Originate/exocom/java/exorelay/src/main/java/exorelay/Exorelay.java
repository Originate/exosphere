package exorelay;

import RequestData.RequestData;
import events.*;

import java.util.Map;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.function.*;

import messageHandler.MessageHandler;
import messageSender.MessageSender;
import org.json.simple.JSONObject;


public class Exorelay implements EventListener {

  private MessageSender messageSender;
  private MessageHandler messageHandler;
  private ZmqListener zmqListener;

  public Exorelay(Map<String, Object> config) throws ExorelayError {
    if(config.get("exocom-port") == null)
      throw new ExorelayError("exocomPort not provided");
    if(config.get("service-name") == null)
      throw new ExorelayError("serviceName not provided");

    EventListener messageHandlerListener = this;
    messageHandler = new MessageHandler(messageHandlerListener);
    EventListener messageSenderListener = this;
    messageSender = new MessageSender(config, messageSenderListener);
    EventListener zmqEventListener = this;
    zmqListener = new ZmqListener(this, zmqEventListener);

  }

  /**
   * Sends a message to ExoCom
   * @param messageName name of message
   * @param payload in JSON format
   * @return generated id of message being sent
   */
  public String send(String messageName, String payload) throws ExorelayError {
    return messageSender.send(messageName, payload);
  }


  void onIncomingMessage(RequestData requestData) throws ExorelayError {
    if(requestData.getName().equals("__status")) messageSender.send("__status-ok", null);

    System.out.println("inside exorelay: " + requestData.getName());
    String response = messageHandler.handleRequest(requestData);
    if(!response.equals("success")) throw new ExorelayError(response);
  }



  @Override
  public void notify(Event type, Object data) throws ExorelayError {
    switch(type) {
      case ERROR:
        throw new ExorelayError(data.toString());
    }
  }

  public Integer getExocomPort() {
    return messageSender.getExocomPort();
  }

  public void closePort() {
    if(messageSender != null) messageSender.closePort();
  }

  public void close() {
    if(zmqListener != null) zmqListener.close();
  }

  public void listen(Integer port) {
    zmqListener.setPort(port);
    ExecutorService es = Executors.newSingleThreadExecutor();
    es.submit(zmqListener);
  }

  public void registerHandler(String messageName, Consumer<JSONObject> handler) throws ExorelayError {
    messageHandler.registerHandler(messageName, handler);
  }

}
