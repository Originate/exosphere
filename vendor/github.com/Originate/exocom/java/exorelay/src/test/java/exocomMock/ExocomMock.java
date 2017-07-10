package exocomMock;


import com.fasterxml.uuid.Generators;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.zeromq.ZMQ;

import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;


public class ExocomMock implements Runnable {

    private Integer pullSocketPort;
    private ZMQ.Context context;
    private ZMQ.Socket pullSocket;
    private Map<String, ZMQ.Socket> pushSockets;
    private List<JSONObject> receivedMessages;
    private Object receiveCallback;

    public ExocomMock(Integer pullSocketPort) {
        this.pullSocketPort = pullSocketPort;
        receivedMessages = new ArrayList<>();
        pushSockets = new HashMap<>();

        context = ZMQ.context(1);
    }

    public void close() {
        for (String service : pushSockets.keySet()) {
            pushSockets.get(service).close();
        }
        pushSockets.clear();
        if(pullSocket != null) pullSocket.close();
        context.close();
    }

    @Override
    public void run() {
        listen();

    }

    private void listen() {
        pullSocket = context.socket(ZMQ.PULL);
        pullSocket.bind("tcp://*:" + pullSocketPort);

        while (!Thread.currentThread().isInterrupted()) {
            try {
                JSONParser parser = new JSONParser();
                String message = pullSocket.recvStr(Charset.forName("UTF-8"));
                JSONObject request = (JSONObject) parser.parse(message);
                onPullSocketMessage(request);
            } catch (ParseException e) {
                System.out.println(e.getMessage());
                e.printStackTrace();
            }

        }

    }

    public void onReceive() {

    }

    public void registerService(String name, Integer port) {
        ZMQ.Socket pushSocket = context.socket(ZMQ.PUSH);
        pushSocket.connect("tcp://localhost:" + port);
        pushSockets.put(name, pushSocket);

    }

    public void reset() {
        receivedMessages.clear();
    }

    public void send(JSONObject data) throws Exception {
        String service = (String) data.get("service");
        if(!pushSockets.containsKey(service)) throw new Exception("unknown service:" + service);

        receivedMessages.clear();
        if(!data.containsKey("messageId")) {
            data.put("messageId", Generators.timeBasedGenerator().generate().toString());
        }
        pushSockets.get(service).send(data.toJSONString());
    }


    public void onPullSocketMessage(JSONObject data) {
        receivedMessages.add(data);
    }

    public List<JSONObject> getReceievedMessages() {
        return receivedMessages;
    }

    public Integer getPort() {
        return pullSocketPort;
    }
}
