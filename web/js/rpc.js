// built in rpc calls
var rpc = {
  invoke : function(arg) { window.external.invoke(JSON.stringify(arg)); },
  init : function() { rpc.invoke({cmd : 'init'}); }
};


$(
  $("#submit").click(function(){
    rpc.invoke({type:"alert" , value:$("#ins").val()});
  }),
  
  $("#refresh_midi_devices").click(function(){
    rpc.invoke({type:"refresh_midi_devices", value:""});
  }),

  $("#clear_midi_devices").click(function(){
    clearDevices ()
  }),

  $("#listen_midi_devices").click(function(){
    rpc.invoke({type:"listen_debug_midi_devices", value:"3"});  // GET THE RIGHT DEVICE ID IN!!!!
  })

);

