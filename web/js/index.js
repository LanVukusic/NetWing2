function addDevice (id, name){
  $(".devices").append('<li class="device"><span>'+id.toString()+'</span><span>'+name.toString()+'</span><input type="checkbox" name="" id="'+id.toString()+'"></li>');
}

function clearDevices (){
  $(".devices").html("");
}

function cliLog (level, type, msg){
  let time = new Date()
  let timeFormated = time.getHours()+":"+time.getMinutes()+":"+time.getSeconds()+"."+time.getMilliseconds().toString()
  let message = $('<div class="cli_line"><div class="cli_time_stamp">'+timeFormated+'</div><div class="cli_type">'+type+'</div><div class="cli_body">'+msg+'</div></div>');
  if(level == 0){
    //ok
    message.addClass("err_ok");
  }
  else if (level == 1){
    //warn
    message.addClass("err_warn");
  }else{
    //error
    message.addClass("err_err");
  }
  $(".cli").append(message);
}