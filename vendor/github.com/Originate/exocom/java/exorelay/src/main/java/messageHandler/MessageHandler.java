package messageHandler;

import events.EventListener;
import events.Event;
import events.ExorelayError;
import org.json.simple.JSONObject;
import RequestData.RequestData;

import java.util.function.Consumer;


public class MessageHandler implements EventListener {

    private HandlerRegistry commandHandlers;
    private HandlerRegistry replyHandlers;
    private EventListener upstreamListener;

    public MessageHandler(EventListener upstreamListener) {
        this.upstreamListener = upstreamListener;

        EventListener commandHandlersListener = this;
        commandHandlers = new HandlerRegistry("message-handler", commandHandlersListener);
        EventListener replyHandlersListener = this;
        replyHandlers = new HandlerRegistry("reply-handler", replyHandlersListener);

    }

    public boolean hasHandler(String messageName) {
        return commandHandlers.hasHandler(messageName);
    }

    public void registerHandler(String messageName, Consumer<JSONObject> handler) throws ExorelayError {
        commandHandlers.registerHandler(messageName, handler);
    }

    public void registerReplyHandler(String messageName, Consumer<JSONObject> handler) throws ExorelayError {
        replyHandlers.registerHandler(messageName, handler);
    }

    public String handleRequest(RequestData requestData) {
        if(requestData.getId() == null) return "missing message id";
        if(commandHandlers.handleCommand(requestData.getName(), requestData.getPayload())) return "success";

        return "unknown message";

    }


    @Override
    public void notify(Event type, Object data) throws ExorelayError {
        switch(type) {
            default:
            case ERROR :
                upstreamListener.notify(Event.ERROR, (String)data);
            break;

        }

    }
}