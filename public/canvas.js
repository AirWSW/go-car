window.onload = function() {
  var dom = document.getElementById("canvas");
  dom.width = 640;
  dom.height = 480;
};

var dis = document.getElementById("dis");
var text = document.getElementById("text");

var ws = new WebSocket("ws://192.168.123.162:8080/subscribe");

var update = function() {
  ws.onmessage = function(event) {
    var str = event.data.toString().split("||");
    var msg_type = str[0];
    var msg_text = str[1];
    switch (msg_type) {
      case "DISTANCE":
        dis.textContent = msg_text;
        break;
      case "CAMINFO":
        text.textContent = msg_text;
        break;
      case "CONTOURS":
        canvas.width = canvas.width;
        var ctx = canvas.getContext("2d");
        eval(
          msg_text
            .replace(
              /\[\[/g,
              'ctx.strokeStyle="red";ctx.lineWidth=3;ctx.beginPath();ctx.moveTo'
            )
            .replace(
              /\] \[/g,
              ";ctx.closePath();ctx.stroke();ctx.beginPath();ctx.moveTo"
            )
            .replace(/\]\]/g, ";ctx.closePath();ctx.stroke();")
            .replace(/ /g, ";ctx.lineTo")
        );
        break;
      default:
        console.log(event.data.toString());
        break;
    }
  };
};

window.setTimeout(update);
