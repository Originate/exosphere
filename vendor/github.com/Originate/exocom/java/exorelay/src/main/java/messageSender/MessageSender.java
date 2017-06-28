package messageSender;

import java.util.Map;


import events.Event;
import events.EventListener;

import events.ExorelayError;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.zeromq.ZMQ;
import com.fasterxml.uuid.Generators;

public class MessageSender {

  private Integer exocomPort;
  private String serviceName;
  private ZMQ.Socket socket;
  private ZMQ.Context context;
  private String lastSentId;
  private EventListener upstreamListener;


  public MessageSender(Map<String, Object> config, EventListener upstreamListener) throws ExorelayError {
    serviceName = (String) config.get("service-name");
    this.upstreamListener = upstreamListener;

    if (config.get("exocom-port") != null)
      exocomPort = (Integer) config.get("exocom-port");
    else
      throw new ExorelayError("ExoCom port not provided");


    context = ZMQ.context(1);
    socket = context.socket(ZMQ.PUSH);
    socket.connect("tcp://localhost:" + exocomPort);

    lastSentId = null;

  }

  public void closePort() {
    socket.close();
    context.close();
  }

  public String send(String messageName, String payload) throws ExorelayError {
    if(messageName == null || messageName.length() == 0)
      upstreamListener.notify(Event.ERROR, "ExoRelay#send cannot send empty messages");


    JSONObject message = new JSONObject();
    message.put("name", messageName);
    message.put("sender", serviceName);
    message.put("id", Generators.timeBasedGenerator().generate().toString());

    if(payload != null) {
      JSONParser parser = new JSONParser();
      try {
        JSONObject payloadObject = (JSONObject) parser.parse(payload);
        message.put("payload", payloadObject);
      } catch (ParseException e) {
        System.out.println("exorelay:message-sender Error parsing payload. Payload must be in JSON format");
      }
    }

    socket.send(message.toJSONString());
    lastSentId = (String) message.get("id");
    return lastSentId;
  }


  private void log(String messageName, Map<String, String> options) {
    if (options.get("responseTo") != null)
      System.out.println("exorelay:message-sender sending message '" + messageName + "' in response to '" + options.get("responseTo") + "'");
    else
      System.out.println("exorelay:message-sender sending message '" + messageName + "'");
  }

  public Integer getExocomPort(){
    return exocomPort;
  }

}
