package messageHandler;


import java.util.HashMap;
import java.util.Map;
import java.util.function.Consumer;

import RequestData.RequestData;
import events.EventListener;
import events.ExorelayError;
import org.json.simple.JSONObject;

class HandlerRegistry {

    private Map<String, Consumer<RequestData>> replyCommandHandlers;
    private Map<String, Consumer<RequestData>> sendCommandHandlers;
    private Map<String, Consumer<JSONObject>> handlers;
    private String debugName;
    private EventListener upstreamListener;

    HandlerRegistry(String debugName, EventListener upstreamListener) {
        replyCommandHandlers = new HashMap<>();
        sendCommandHandlers = new HashMap<>();
        handlers = new HashMap<>();

        this.debugName = debugName;
        this.upstreamListener = upstreamListener;
    }


    boolean handleCommand(String messageName, JSONObject payload) {
        Consumer<JSONObject> handler = getHandler(messageName);
        if(handler != null) {
            System.out.println("exorelay:" + debugName + " handling message '" + messageName + "'");
            handler.accept(payload);
            return true;
        } else {
            return false;
        }
    }

    void registerHandler(String messageName, Consumer<JSONObject> handler) throws ExorelayError {
        if (hasHandler(messageName)) throw new ExorelayError("There is already a handler for message '" + messageName + "'");
        System.out.println("exorelay:" + debugName + "  registering handler for id '" + messageName + "'");

        handlers.put(messageName, handler);
    }

    public boolean hasHandler(String messageName) {
        return getHandler(messageName) != null;
    }
    public Consumer<JSONObject> getHandler(String messageName) {
        return handlers.get(messageName);
    }


}
