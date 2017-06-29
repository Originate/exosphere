package RequestData;


import org.json.simple.JSONObject;

public class RequestData {

    private String name;
    private JSONObject payload;
    private String responseTo;
    private String id;

    public RequestData(String name, JSONObject payload, String responseTo, String id){
        this.name = name;
        this.payload = payload;
        this.responseTo = responseTo;
        this.id = id;
    }

    public String getName() { return name; }

    public JSONObject getPayload() { return payload; }

    public String getResponseTo() { return responseTo; }

    public String getId() { return id; }
}
