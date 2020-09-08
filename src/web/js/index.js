// GLOBALS
var USED_FADER_IDS = [-1]
var USED_EXEC_PAGES = []
var curr_exec_page


// FUNCTIONS
function addDevice(id, name) {
  $(".devices").append('<li class="device"><span>' + id.toString() + '</span><span>' + name.toString() + '</span><input type="checkbox" name="" id="' + id.toString() + '"></li>');
}

function addInterfaceInstance(id, Hname, FriendlyName) {
  $("#interface-space").append('<div class="interface_inst" id ="' + id.toString() + '"><div class="inst_back"><div class="title">' + FriendlyName + '</div><div>' + Hname + ' : ' + id.toString() + '</div></div><div class="inst_last"><button>Settings</button></div></div>');
}

function addFaderInstance(fader_channel, midi_chan) {
  USED_FADER_IDS.push(parseInt(fader_channel));
  $(".faders_holder").append('<div class="fader" style="order:' + fader_channel + '"><button id="fader-edit-button" onclick="faderEdit(this)">edit</button><input disabled="" type="range" orient="vertical" max="127" min="0" class="slider" id="fader' + fader_channel + '"><div><span><span>MIDI:</span><i id="fader-label-midi">' + midi_chan + '</i></span><span>Exec:<i id="fader-label-exec">' + fader_channel + '</i></span></div></div>');
}

function updateExecInstance(fader_channel, exec_page, midi_chan) {
  $("#exec_page_" + exec_page).find("div[itemid='" + fader_channel + "']").find("#exec_mapping").html(midi_chan);
  $("#exec_page_" + exec_page).find("#exec_item" + fader_channel).attr("isset", "1")
  $(".modal_execs").addClass("disabled");
}

function clearDevices() {
  $(".devices").html("");
}

function cliLog(level, type, msg) {
  let time = new Date()
  let timeFormatted = time.getHours().toString().padStart(2, '0') + ":" + time.getMinutes().toString().padStart(2, '0') + ":" + time.getSeconds().toString().padStart(2, '0') + "." + time.getMilliseconds().toString().padStart(3, '0')
  let cli = $("#cliLog")
  let message = $('<div class="cli_line"><div class="cli_time_stamp">' + timeFormatted + '</div><div class="cli_type">' + type + '</div><div class="cli_body">' + msg + '</div></div>');

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

function faderEdit(elem) {
  let faderID = $(elem).parent().find('#fader-label-exec').html();
  let faderMIDI = $(elem).parent().find('#fader-label-midi').html();
  $("#fader_span_title").html(" : " + faderID);
  $("#fader_input_num").val(faderID);
  $("#fader_label_status").html("mapped");
  $("#fader_label_midi_chn").html(faderMIDI);
  $("#fader-update-button").html("update").attr("disabled", true);
  $("#fader-remove-button").removeClass("disabled").attr("disabled", false)
  $("#faders-modal").removeClass("disabled");
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

function setMIDILearnFader(device, channel) {
  $("#fader_label_status").html("mapped");
  $("#fader_label_midi_chn").html("" + device + "." + channel);
  $("#fader-update-button").prop('disabled', false);
}

function setMIDILearnExec(device, channel) {
  $("#exec_label_status").html("mapped");
  $("#exec_label_midi_chn").html("" + device + "." + channel);
  $("#exec-update-button").prop('disabled', false);
  //console.log(device, channel)
}

function exec_tile_func() {
  $("#exec_span_title").html($(this).attr('itemid'));
  if ($(this).attr('isset') == 0) {
    // is not set
    $("#exec-exists-error").html("Adding new exec mapping");
    $("#exec_label_status").html("unmapped");
    $("#exec_label_midi_chn").html("/");
    $("#exec-update-button").html("Create");
    $("#exec-remove-button").addClass("disabled");
  } else {
    // is set already
    $("#exec-exists-error").html("Update existing exec");
    $("#exec_label_status").html("Mapped");
    $("#exec_label_midi_chn").html("" + $(this).find("#exec_mapping").html());
    $("#exec-update-button").html("Update");
    $("#exec-update-button").prop('disabled', false);
    $("#exec-remove-button").removeClass("disabled");
    $("#exec-remove-button").prop('disabled', false);
  }

  $("#exec-update-button").attr('disabled', 'disabled');
  $("#execs-modal").removeClass("disabled");
}

function change_exec_page(obj) {
  curr_exec_page = parseInt(obj.target.innerText)
  // update main look
  $('.exec_page').each(function (i, obj) {
    $(obj).addClass("disabled")
  });

  $("#exec_page_" + $(this).text().toString().toLowerCase()).removeClass("disabled")

  //update menu look
  $('.page_holder').each(function (i, obj) {
    $(obj).removeClass("page_holder_active")
  });

  $(this).addClass("page_holder_active")
}

function add_exec_page(page, width, heigh) {
  let exec_page = $('<div>')
  exec_page.addClass('execs_up');
  exec_page.addClass('exec_page');
  exec_page.attr({
    id: 'exec_page_' + page
  });

  let counter = 1;

  for (let row = 0; row < parseInt(width); row++) {
    let rowHTML = $('<div>')
    rowHTML.addClass('execs_up_row');
    for (let col = 0; col < parseInt(heigh); col++) {
      let colHTML = $('<div class="execs_up_item" id="exec_item' + counter + '" itemid="' + counter + '" isset="0">' + counter + '<br><label id="exec_mapping"></label><br></div>')
      $(colHTML).click(exec_tile_func);
      rowHTML.append(colHTML)
      counter++;
    }
    exec_page.append(rowHTML);
  }

  curr_exec_page = page;
  USED_EXEC_PAGES.push(page)

  $(".main_execs").append(exec_page);


  // change to current page
  $('.exec_page').each(function (i, obj) {
    $(obj).addClass("disabled")
  });
  $("#exec_page_" + page).removeClass("disabled")

  //update menu look
  $('.page_holder').each(function (i, obj) {
    $(obj).removeClass("page_holder_active")
  });

  let page_change_button = $('<div class="page_holder page_holder_active" id="page_button_' + page + '">' + page + '</div>')
  page_change_button.click(change_exec_page)
  $(".execs_scrollbar").append(page_change_button)


}

function downloadObjectAsJson(exportObj, exportName){
  var dataStr = "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(exportObj));
  var downloadAnchorNode = document.createElement('a');
  downloadAnchorNode.setAttribute("href",     dataStr);
  downloadAnchorNode.setAttribute("download", exportName + ".json");
  document.body.appendChild(downloadAnchorNode); // required for firefox
  downloadAnchorNode.click();
  downloadAnchorNode.remove();
}

function loadSaveFile(obj){
  //console.log(obj);
  let data = {
    event: "loadSave",
    data: obj
  }
  data = JSON.stringify(data)
  conn.send(data)
}

$(
  $(".side_block").click(function () {

    // update main look
    $('.main_window').each(function (i, obj) {
      $(obj).addClass("disabled")
    });
    $(".main_" + $(this).text().toString().toLowerCase()).removeClass("disabled")

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

  $("#addInterfaceOSC").click(function () {
    let host = $("#OSC_ip_add").val();

    let data = {
      event: "restartOSC",
      host: host
    }
    data = JSON.stringify(data)
    conn.send(data)
    console.log(data);

  }),

  $("#addInterfaceGenericMIDI").click(function () {
    $(".modal_interfaces").removeClass("disabled");
  }),

  $("#closeModal").click(function () {
    $(".modal_faders").addClass("disabled");
  }),

  $("#closeModalInterf").click(function () {
    $(".modal_interfaces").addClass("disabled");
  }),

  $("#closeModalExecs").click(function () {
    $(".modal_execs").addClass("disabled");
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

  $("#applyDevice").click(function () {

    let inDev = null;
    let outDev = null;
    let hNameIn = null;
    let hNameOut = null;

    //get in device
    $("#TableMidiIns").children().each(function (i, obj) {
      //console.log($(obj).attr('class'), $(obj).hasClass("selectedDevice"));
      if ($(obj).hasClass("selectedDevice")) {
        inDev = i;
        hNameIn = $(obj).children().eq(1).text()

        //return false; // breaks
      }
    });
    //get out device
    $("#TableMidiOuts").children().each(function (i, obj) {
      if ($(obj).hasClass("selectedDevice")) {
        outDev = i;
        hNameOut = $(obj).children().eq(1).text()
        //outDev = $(obj).children('div').eq(1).text();
        //return false; // breaks
      }
    });

    // device types : 0 MIDI, 1 OSC , 2 ART-NET
    data = {
      event: "addInterface",
      inDevice: inDev,
      outDevice: outDev,
      deviceType: 0,
      HardwareName: hNameIn,
      FriendlyName: $("#dDisplayName").val()
    }

    // alerts user to select the device
    if (inDev == null || outDev == null) {
      $("#noDeviceAlert").removeClass("disabled");
    } else {
      conn.send(JSON.stringify(data));
      $(".modal").addClass("disabled");
    }
  }),

  $("#execs-add-fader").click(function () {
    $("#fader-exists-error").addClass("disabled")
    $("#fader_input_num").val(USED_FADER_IDS[USED_FADER_IDS.length - 1] + 1);
    $("#fader_span_title").html(": New");
    $("#fader_label_status").html("unmapped");
    $("#fader_label_midi_chn").html("/");
    $("#fader-update-button").html("Add new");
    $("#fader-update-button").attr('disabled', 'disabled');
    $("#fader-remove-button").addClass("disabled");
    $("#faders-modal").removeClass("disabled");
  }),

  $("#fader-remove-button").click(function () {
    let page = 0; // TO DO
    let faderId = $("#fader_input_num").val();

    data = {
      event: "removeMapping",
      extChn: parseInt(faderId),
      execPage: page,
      extType: 0, // 0 = fader, 3 = exec
    }
    conn.send(JSON.stringify(data));
  }),

  $("#fader-update-button").click(function () {
    let fader_channel = $("#fader_input_num").val(); // fader number on Magicq
    let midi_chan = $("#fader_label_midi_chn").html(); // MIDI channel that is mapped to the fader

    if (USED_FADER_IDS.includes(parseInt(fader_channel))) {
      // detected problem while trying to add an existing fader
      // display err
      $("#fader-exists-error").removeClass("disabled")
      return;
    }

    data = {
      event: "bindMIDIchannel", // type of event
      device: Math.floor(parseInt(midi_chan)), // number of midi device, basically a device id
      chn: parseInt(midi_chan.split('.')[1]), // MIDI channel
      extChn: parseInt(fader_channel), // number of fader to controll
      execPage: 0, // fader page
      typeFader: true, // should it fade or snap
      extType: 0 // 0 = fader, 3 = exec
    }

    conn.send(JSON.stringify(data))

    // close the window
    $("#faders-modal").addClass("disabled");
  }),

  $("#fader-learn-button").click(function () {
    // transmit the listening mode to the server
    data = {
      event: "changeMIDImode",
      interface: 0
    }
    conn.send(JSON.stringify(data));
  }),

  $("#exec-learn-button").click(function () {
    // transmit the listening mode to the server
    data = {
      event: "changeMIDImode",
      interface: 3,
    }
    conn.send(JSON.stringify(data));
  }),

  $("#exec-update-button").click(function () {
    let page = parseInt(curr_exec_page); // TOO DOOOOO
    let form = $(this).parent().parent().find(".left");
    let isFader = form.find("input[name='exec_type']:checked").val();
    let midiOut = form.find("#midiOut").val();
    let midi_chan = form.find("#exec_label_midi_chn").html();
    let execId = $("#exec_span_title").html();

    data = {
      event: "bindMIDIchannel",
      device: Math.floor(parseInt(midi_chan)),
      chn: parseInt(midi_chan.split('.')[1]),
      extChn: parseInt(execId),
      execPage: page,
      extType: 3, // 0 = fader, 4 = exec
      typeFader: (isFader == "fader") ? true : false,
      //feedback: midiOut
    }
    conn.send(JSON.stringify(data))

    // close the window
    $("#modal_window").addClass("disabled");
  }),

  $("#exec-remove-button").click(function () {
    let page = parseInt(curr_exec_page);
    let form = $(this).parent().parent().find(".left");
    let midi_chan = form.find("#exec_label_midi_chn").html();
    let execId = $("#exec_span_title").html();

    data = {
      event: "removeMapping",
      extChn: parseInt(execId),
      execPage: page,
      extType: 3, // 0 = fader, 3 = exec
    }
    conn.send(JSON.stringify(data));
  }),

  $(".execs_plus").click(function () {
    $(".modal_window").removeClass("disabled");
  }),

  $("#closeModalWindow").click(function () {
    $(".modal_window").addClass("disabled");
  }),

  $("#button_add_window").click(function (e) {
    let window = $("#exec_window_win").val();
    let page = parseInt($("#exec_window_page").val());

    if (window == "" || page == "" || Number.isNaN(page)) {
      $("#exec-win-exists-error").html("No fields can be left empty.").removeClass("disabled")
      return
    }
    if (parseInt(page) == NaN) {
      $("#exec-win-exists-error").html("Page input is not valid.").removeClass("disabled")
      return
    }
    if (Number.isNaN(parseInt(window.split("/")[0])) || Number.isNaN(parseInt(window.split("/")[1]))) {
      $("#exec-win-exists-error").html("Windows size input is not valid.").removeClass("disabled")
      return
    }

    if (USED_EXEC_PAGES.includes(page)) {
      // user is naughty and wants to add an existing page
      $("#exec-win-exists-error").html("This page number already exists.").removeClass("disabled")
      return
    } else {
      $("#exec-win-exists-error").addClass("disabled");
    }

    data = {
      event: "addNewPage",
      page: parseInt(page),
      width: parseInt(window.split("/")[0]),
      height: parseInt(window.split("/")[1])
    }

    //console.log(parseInt(window.split("/")[0]), parseInt(window.split("/")[1]), "penis penis")
    conn.send(JSON.stringify(data))
    $("#window-modal").addClass("disabled");
  }),

  // drag and drop start
  $("#drag_area").on('drag dragstart dragend dragover dragenter dragleave drop', function (e) {
    e.preventDefault();
    e.stopPropagation();
  })
  .on('dragover dragenter', function () {
    $("#drag_area").addClass('dr_active');
  })
  .on('dragleave dragend drop', function () {
    $("#drag_area").removeClass('dr_active');
  }).on('drop', function (e) {
    let droppedFiles = e.originalEvent.dataTransfer.files;
    var reader = new FileReader();
    reader.onload = (event) => {
      loadSaveFile(JSON.parse(event.target.result))
    };
    reader.readAsText(droppedFiles[0]);
  }),

  $('#file').change(function (event) {
    var reader = new FileReader();
    reader.onload = (event) => {
      loadSaveFile(JSON.parse(event.target.result))
    };
    reader.readAsText(event.target.files[0]);
  }),
  //drag and drop end

  $("#btn_save_dl").click(()=>{
    data = {
      event: "saveRequest",
      type: "local" 
    }
    conn.send(JSON.stringify(data))
  })




);