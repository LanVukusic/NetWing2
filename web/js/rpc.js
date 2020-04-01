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
    switch(evt.Event){
      case "cli":
        cliLog(evt.ThreatLevel, evt.Cause, evt.Body)
        break
      case "refreshMidiRet":
        updateMIDItable(evt)
        break
      case "UiAddDevice":
        evt = JSON.parse(evt.Data.replace(/'/g,"\""));
        addInterfaceInstance(evt.ID, evt.Hname, evt.FriendlyName);
        $(".modal").addClass("disabled");
        break
      case "learnMidiRet":
        //evt = JSON.parse(evt.Data.replace(/'/g,"\""));
        setMIDILearnFader(evt.DeviceID, evt.ChannelID);
        break
      case "MappingsResponse":
        addFaderInstance(evt.FaderID, `${evt.DeviceID}.${evt.ChannelID}`)
        break
      case "UpdateFader":
        $("#fader"+evt.FaderID).val(parseInt(evt.Value));
    }
  }
} else {
  alert("Your browser does not support websocket connection. Interface is not operational.")
}

