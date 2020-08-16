var conn

if (window["WebSocket"]) {
  conn = new WebSocket("ws://" + document.location.host + "/ws");

  
  conn.onclose = function (evt) {
    console.log("SERVER CONNECTION DROPPED")
    cliLog(2, "Server Connection", "Connection to backend dropped. You are offline")
    // ADD THE ERROR WARNING HERE
  };

  conn.onmessage = function (evt) {
    evt = JSON.parse(evt.data)
    console.log(evt)
    switch(evt.Event){
      case "cli":
        cliLog(evt.ThreatLevel, evt.Cause, evt.Body)
        break
      case "refreshMidiRet":
        updateMIDItable(evt)
        break
      case "UiAddDevice":
        console.log(evt.Data.replace(/'/g,"\""));
        evt = JSON.parse(evt.Data.replace(/'/g,"\""));
        addInterfaceInstance(evt.ID, evt.Hname, evt.FriendlyName);
        //addInterfaceInstance(12, "evt.Hname", "evt.FriendlyName")
        $(".modal").addClass("disabled");
        break
    }
  }
} else {
  alert("Your browser does not support websocket connection. Interface is not operational.")
}

