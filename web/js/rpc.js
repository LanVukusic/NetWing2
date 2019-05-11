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
  })
);