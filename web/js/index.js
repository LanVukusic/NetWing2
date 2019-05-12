function addDevice (id, name){
  $(".devices").append('<li class="device"><span>'+id.toString()+'</span><span>'+name.toString()+'</span><input type="checkbox" name="" id="'+id.toString()+'"></li>');
}

function clearDevices (){
  $(".devices").html("");
}