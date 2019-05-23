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

