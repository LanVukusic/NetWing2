var socket = io();

socket.on('refreshMidiRet', function (msg) {
  //populate ins
  data = JSON.parse(msg);

  if (data.Ins == null) {
    $("#TableMidiIns").empty();
    $("#TableMidiIns").append('<div class="deviceTableDevice" ><div>No devices found</div></div>');
  } else {
    data.Ins.forEach(function (element) {
      $("#TableMidiIns").empty();
      $("#TableMidiIns").append('<div class="deviceTableDevice" id="MidiListDevice"><div>' + element.ID + '</div><div>' + element.Name + '</div></div>');
    });
  }

  if (data.Outs == null) {
    $("#TableMidiOuts").empty();
    $("#TableMidiOuts").append('<div class="deviceTableDevice"><div>No devices found</div></div>');
  } else {
    data.Outs.forEach(function (element) {
      $("#TableMidiOuts").empty();
      $("#TableMidiOuts").append('<div class="deviceTableDevice" id="MidiListDevice"><div>' + element.ID + '</div><div>' + element.Name + '</div></div>');
    });
  }
});

socket.on("CliLog", function(msg){
  data = JSON.parse(msg);
  let threatLevel = "err_ok"
  if(data.level == 1){
    threatLevel = "err_warn"
  }else if (data.level == 2){
    threatLevel ="err_err"
  }
  $(".cli").append('<div class="cli_line '+threatLevel+'+><div class="cli_time_stamp">'+Date().getTime()+'</div><div class="cli_type">'+data.cause+'</div><div class="cli_body">'+data.msg+'</div></div>')
  });

