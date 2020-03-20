function addDevice(id, name) {
  $(".devices").append('<li class="device"><span>' + id.toString() + '</span><span>' + name.toString() + '</span><input type="checkbox" name="" id="' + id.toString() + '"></li>');
}

function addInterfaceInstance (id, Hname, FriendlyName){
  $(".interfaces_workplace").append('<div class="interface_inst" id ="'+id.toString()+'"><div class="inst_front"><input type="checkbox" value="1" name="" id=""></div><div class="inst_back"><div class="title">'+FriendlyName+'</div><div>'+Hname+' : '+id.toString()+'</div></div></div>');
}

function clearDevices() {
  $(".devices").html("");
}

function cliLog(level, type, msg) {
  let time = new Date()
  let timeFormated = time.getHours() + ":" + time.getMinutes() + ":" + time.getSeconds() + "." + time.getMilliseconds().toString()
  let cli = $("#cliLog")
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

function updateMIDItable(data) {
  //populate ins

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
gridster.add_widget("<div class='widget'>Test</div>", [5], [6], [0], [0])




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


  $("#RefreshDevice").click(function () {
    let data = {
      event: "getMidiDevices",
      data: ""
    }
    data = JSON.stringify(data)
    conn.send(data)
  }),

  $("#addInterfaceGenericMIDI").click(function () {
    $(".modal").removeClass("disabled");
  }),

  $("#closeModal").click(function () {
    $(".modal").addClass("disabled");
  }),

  $('.devList').on('click', '#MidiListDevice', function () {
    $(this).parent().children('div').each(function (i, obj) {
      $(obj).removeClass("selectedDevice")
    });
    $(this).toggleClass("selectedDevice");
  }),

  $("#cli_clear").click(function () {
    $(".cli").html("");
  }),

  $("#test").click(function () {
    cliLog(1, "test", "this is a tasty test")
  }),

  $("#applyDevice").click(function(){
    
    let inDev = null;
    let outDev = null;
    let hNameIn = null;
    let hNameOut = null;

    //get in device
    $("#TableMidiIns").children().each(function (i, obj) {
      //console.log($(obj).attr('class'), $(obj).hasClass("selectedDevice"));
      if($(obj).hasClass("selectedDevice")){
        inDev = i;
        hNameIn = $(obj).children().eq(1).text()

        //return false; // breaks
      }
    });
    //get out device
    $("#TableMidiOuts").children().each(function (i, obj) {
      if($(obj).hasClass("selectedDevice")){
        outDev = i;
        hNameOut = $(obj).children().eq(1).text()
        //outDev = $(obj).children('div').eq(1).text();
        //return false; // breaks
      }
    });

    // device types : 0 MIDI, 1 OSC , 2 ART-NET
    data = {
      event : "addInterface",
      inDevice : inDev,
      outDevice : outDev,
      deviceType : 0,
      HardwareName: hNameIn,
      FriendlyName: $("#dDisplayName").val()
    }
    // alerts user to select the device
    if(inDev == null || outDev==null){
      $("#noDeviceAlert").removeClass("disabled");
    }else{
      conn.send(JSON.stringify(data));
      $(".modal").addClass("disabled");
    }
  })
);