package events;

public interface EventListener {

    void notify(Event type, Object data) throws ExorelayError;
}
