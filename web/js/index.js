function addDevice(id, name) {
  $(".devices").append('<li class="device"><span>' + id.toString() + '</span><span>' + name.toString() + '</span><input type="checkbox" name="" id="' + id.toString() + '"></li>');
}

function clearDevices() {
  $(".devices").html("");
}

function cliLog(level, type, msg) {
  let time = new Date()
  let timeFormated = time.getHours() + ":" + time.getMinutes() + ":" + time.getSeconds() + "." + time.getMilliseconds().toString()
  let cli = $(".cli")
  let message = $('<div class="cli_line"><div class="cli_time_stamp">' + timeFormated + '</div><div class="cli_type">' + type + '</div><div class="cli_body">' + msg + '</div></div>');

  if (level == 0) {
    //ok
    message.addClass("err_ok");
  } else if (level == 1) {
    //warn
    message.addClass("err_warn");
  } else {
    //error
    message.addClass("err_err");
  }
  cli.append(message);
  cli.scrollTop(cli.prop("scrollHeight"))
}



let gridSize = 30
console.log($(".gridster").width());
console.log($(".gridster").height());

let myGrid = $(".gridster").gridster({
  widget_margins: [0, 0],
  widget_base_dimensions: [gridSize, gridSize],
  widget_selector: ".widget",
  shift_widgets_up: false,
  min_cols: Math.floor($(".gridster").width() / gridSize),
  max_cols: Math.floor($(".gridster").width() / gridSize),
  min_rows: Math.floor($(".gridster").height() / gridSize),
  max_rows: Math.floor($(".gridster").height() / gridSize)
});

var gridster = $(".gridster").gridster().data('gridster');
gridster.add_widget("<div class='widget'>Test</div>", [5], [6], [0], [0] )

$(
  $(".side_block").click(function () {
    // update main look
    $('.main_window').each(function (i, obj) {
      $(obj).addClass("none")
    });
    $(".main_" + $(this).text().toString().toLowerCase()).removeClass("none")

    //update menu look
    $('.side_block').each(function (i, obj) {
      $(obj).removeClass("block_active")
    });
    $(this).addClass("block_active")
  }),


  $("#RefreshDevice").click(function(){
    socket.emit("refreshMidi", "");
  }),

  $("#addInterfaceGenericMIDI").click(function(){
    $(".modal").removeClass("disabled");
  }),

  $("#closeModal").click(function(){
    $(".modal").addClass("disabled");
  }),

  $('.devList').on('click', '#MidiListDevice', function() {
    $(this).parent().children('div').each(function (i, obj) {
      $(obj).removeClass("selectedDevice")
    });
    $(this).toggleClass("selectedDevice");
  }),

  $("#applyDevice").click(function(){
    
    let inDev = null;
    let outDev = null;

    //get in device
    $("#TableMidiIns").children().each(function (i, obj) {
      //console.log($(obj).attr('class'), $(obj).hasClass("selectedDevice"));
      if($(obj).hasClass("selectedDevice")){
        inDev = i;
        //return false; // breaks
      }
    });
    //get out device
    $("#TableMidiOuts").children().each(function (i, obj) {
      if($(obj).hasClass("selectedDevice")){
        outDev = i;
        //return false; // breaks
      }
    });

    // device types : 0 MIDI, 1 OSC , 2 ART-NET
    data = {
      inDevice : inDev,
      outDevice: outDev,
      deviceType : 0
    }

    console.log(inDev , outDev)
    // alerts user to select the device
    if(inDev == null || outDev==null){
      $("#noDeviceAlert").removeClass("disabled");
    }else{
      socket.emit("AddDevice", JSON.stringify(data));
    }
    
  })
);




