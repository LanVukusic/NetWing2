// built in rpc calls
var rpc = {
  invoke: function (arg) {
    window.external.invoke(JSON.stringify(arg));
  },
  init: function () {
    rpc.invoke({
      cmd: 'init'
    });
  }
};

// handlers and event listeners
$(
  $("#submit").click(function () {
    //rpc.invoke({type:"alert" , value:$("#ins").val()});
    cliLog(1, "alert", $("#ins").val())
  }),

  $("#refresh_midi_devices").click(function () {
    rpc.invoke({
      type: "refresh_midi_devices",
      value: ""
    });
  }),

  $("#clear_midi_devices").click(function () {
    clearDevices()
  }),

  $("#listen_midi_devices").click(function () {
    rpc.invoke({
      type: "listen_debug_midi_devices",
      value: "3"
    }); // GET THE RIGHT DEVICE ID IN!!!!
  }),

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

  /* $("#addInterfaceGenericMIDI").click(function () {
    //add interface to the list


    //call backend

    //apply handlers from the backend to the object
  }) */

  $("#RefreshDevice").click(function () {
    //alert(counter.value)
    couter.Add(1)
    //add interface to the list


    //call backend

    //apply handlers from the backend to the object
  })


);

