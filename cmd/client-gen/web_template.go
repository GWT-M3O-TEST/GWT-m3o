package main

const webHTMLServiceTemplate = `{{ $reqType := requestType .endpoint }}{{ $service := .service -}}
<html>
  <body>
    <div id="client">
      <form id="client-call" onsubmit="call()">
        <div>
          <input name="token" id="token" placeholder="token">
        </div>
        <div>
          <input name="service" id="service" placeholder="{{ $service.Name }}">
        </div>
        <div>
          <input name="endpoint" name="endpoint" placeholder="{{ .endpoint }}">
        </div>
        <div>
          <textarea rows=5 cols=30 name="request" id="request">{}</textarea>
        </div>
        <button>Submit</button>
      </form>
    </div>
    <div id="response"></div>
  </body>
  <script src="index.js"></script>
</html>`

const webJSServiceTemplate = `{{ $reqType := requestType .endpoint }}{{ $service := .service -}}
class {{ title $service.Name }} {
	constructor(token) {
	  this.token = token;
	}
  
	call({{ $service.Name }}, {{ .endpoint }}, request, callback) {
	  // e.g /v1/helloworld/Call
	  var path = "/v1/" + {{ $service.Name }} + "/" + {{ .endpoint }}
  
	  var xmlHttp = new XMLHttpRequest();
	  xmlHttp.onreadystatechange = function () {
		if (xmlHttp.readyState == 4);
		callback(xmlHttp.responseText, xmlHttp.status);
	  };
	  xmlHttp.open("POST", "https://api.m3o.com" + path, true); // true for asynchronous
	  xmlHttp.setRequestHeader("Authorization", "Bearer " + this.token);
	  xmlHttp.setRequestHeader("Content-Type", "application/json");
	  xmlHttp.send(request);
	}
  }
  
  function call() {
		var form = document.getElementById("client-call");
		var token = form.elements["token"].value;
		var service = form.elements["service"].value;
		var endpoint = form.elements["endpoint"].value;
		var request = form.elements["request"].value;
  
		var m3o = new Client(token);
  
		m3o.call({{ $service.Name }}, {{ .endpoint }}, request, function(response) {
		  document.getElementById("response").innerText = response;
		});
  }  
`
