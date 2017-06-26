package exorelay;


import RequestData.RequestData;
import events.*;
import org.json.simple.JSONObject;
import org.json.simple.parser.JSONParser;
import org.json.simple.parser.ParseException;
import org.zeromq.ZMQ;
import org.zeromq.ZMQException;

import java.nio.charset.Charset;
import java.util.concurrent.Callable;

public class ZmqListener implements Callable<Void>{

    private ZMQ.Socket socket;
    private ZMQ.Context context;
    private Integer port;
    private Exorelay exorelay;
    private EventListener upstreamListener;

    public ZmqListener(Exorelay exorelay, EventListener upstreamListener) {
        context = ZMQ.context(1);
        socket = context.socket(ZMQ.PULL);
        this.upstreamListener = upstreamListener;

        this.exorelay = exorelay;
    }

    public void setPort(Integer port) {
        this.port = port;
    }

    public void close() {
        System.out.println("exorelay:zmq-listener no longer listening at port " + port);
        socket.close();
//        context.term();
        socket = null;
        port = null;
        try {
            upstreamListener.notify(Event.ZMQ_OFFLINE, null);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }


    @Override
    public Void call() throws ExorelayError, ZMQException {
        try {
            socket.bind("tcp://*:" + port);
            upstreamListener.notify(Event.ZMQ_ONLINE, port);

            while (!Thread.currentThread().isInterrupted()) {
                System.out.println("I AM HERE");
                context.getMaxSockets();
                String request = socket.recvStr(Charset.forName("UTF-8"));
                onZmqSocketMessage(request);
            }

        } catch(ZMQException e) {
            if(e.getErrorCode() == ZMQ.Error.EADDRINUSE.getCode()) throw new ExorelayError("port " + port + " for ExoRelay is already in use");
            else if(e.getErrorCode() != ZMQ.Error.ETERM.getCode()) throw e;
        }

        return null;
    }

    private void onZmqSocketMessage(String data) throws ExorelayError {
        RequestData requestData = parseRequest(data);
        log(requestData);

        exorelay.onIncomingMessage(requestData);
    }

    private RequestData parseRequest(String data) {
        JSONParser parser = new JSONParser();
        JSONObject root = null;
        RequestData request = null;
        try {
            root = (JSONObject) parser.parse(data);

            String messageName = (String) root.get("name");
            String id = (String) root.get("messageId");
            JSONObject payload = null;
            String responseTo = null;
            if(!messageName.equals("__status")) {
                payload = (JSONObject) parser.parse((String) root.get("payload"));
                responseTo = (String) root.get("responseTo");
            }
            request = new RequestData(messageName, payload, responseTo, id);
        } catch (ParseException e) {
            System.out.println(e.getMessage());
            e.printStackTrace();
        }

        return request;

    }

    private void log(RequestData requestData) {
        System.out.print("exorelay:zmq-listener received message '" + requestData.getName() +
                            "' with id '" + requestData.getId() + "'");

        String responseTo = requestData.getResponseTo();
        if(responseTo !=null) System.out.println(" in response to '" + responseTo + "'");

    }



}
