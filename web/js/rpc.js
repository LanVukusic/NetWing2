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
        console.log(evt.Interf);
        switch(evt.Interf){
          case 0:
            setMIDILearnFader(evt.DeviceID, evt.ChannelID);
            break
          case 3:
            setMIDILearnExec(evt.DeviceID, evt.ChannelID);
            break
        }
        break

      case "MappingsResponse":
        if(evt.Interface == 0){
          addFaderInstance(evt.FaderID, `${evt.DeviceID}.${evt.ChannelID}`);
        }else{
          if( evt.Interface == 3){
            updateExecInstance(evt.ExecID, evt.ExecPage, `${evt.DeviceID}.${evt.ChannelID}`);
          }
        }
        break

      case "UpdateFader":
        if(evt.Type == 0){  // update fader
          $("#fader"+evt.FaderID).val(parseInt(evt.Value));
        }else if (evt.Type == 3){  // update exec
          let val = Math.floor(evt.Value * 100 / 127)
          $("#exec_page_"+evt.PageID).find("#exec_item"+evt.FaderID).css("background", "linear-gradient(90deg, #ffffff57 "+val+"%, #ffffff1e "+val+"%)")
          $("#fader"+evt.FaderID).val(parseInt(evt.Value));
        }
        break
      
        case "newExecPage":
          add_exec_page(evt.Page,evt.Width, evt.Height)
      

        
    }
  }
} else {
  alert("Your browser does not support websocket connection. Interface is not operational.")
}
