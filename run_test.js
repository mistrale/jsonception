var http = require('http');
var WebSocketClient = require('websocket').client;

var options = {
  host: 'localhost',
  path: '/libraries/1/run?idLib=1',
  port:9000,
  method: 'POST',
  //This is the only line that is new. `headers` is an object with the headers to request
};

var uuid =""
var client = new WebSocketClient();

callback = function(response) {
  var str = ''
  response.on('data', function (chunk) {
	  var obj = JSON.parse(chunk);
	  if (obj["status"] == false) {
		process.exit(-1);
	  } else {
		uuid = obj["response"]
		test = client.connect('ws://localhost:9000/websocket/room?room_name=' + uuid, null, "http://localhost:9000");
	  }
  });
}

client.on('connect', function(connection) {
	console.log('WebSocket Client Connected');
	connection.on('error', function(error) {
		console.log("Connection Error: " + error.toString());
	});
	connection.on('close', function() {
		console.log('echo-protocol Connection Closed');
	});
	connection.on('message', function(message) {
		if (message.type === 'utf8') {
			var obj = JSON.parse(message.utf8Data);
			if (obj["status"] == false) {
				throw new Error("Something went badly wrong!"  + message.utf8Data);
			} else {
				console.log("Received event")
			}
		}
	});
});

client.on('connectFailed', function(connection) {
	console.log(connection)
});

var req = http.request(options, callback);
req.end();

