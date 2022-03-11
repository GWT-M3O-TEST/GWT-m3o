package main

// Range over endpoint attributes
// for property, val := range meta.Value.Properties {
// 	propDescription := val.Value.Description
// 	fmt.Println("attribute:", property)
// 	fmt.Println("placeholder:", propDescription)
// }

const webHTMLServiceTemplate = `
{{ $service := .service -}}
<html>
  <body>
    <div id="{{ $service.Name }}">
      <form id="{{ .endpoint }}" onsubmit="call()">
        <div>
            <label for="service"><b>{{ $service.Name }}</b></label>
            <input type="hidden" name="service" id="service">
        </div>
        <div>
            <label for="endpoint"><b>{{ .endpoint }}</b></label>
            <input type="hidden" name="endpoint" id="endpoint">
        </div>
		<i>{{ .epdesc }}</i>
		</br>
		</br>
		<label for="token">token </label>
        <div>
            <input name="token" id="token" placeholder="token">
        </div>
		</br>
		{{- range $property, $val := .properties }}
		{{- if not (eq $val.Value.Type "object") }}
		<label for="{{ $property }}">{{ $property }} </label>
        <div>
            <input name="{{ $property }}" id="{{ $property }}" placeholder="{{ $val.Value.Description }}">
        </div>
		{{- end }}
		{{- if eq $val.Value.Type "object" }}
		<label for="{{ $property }}">{{ $property }} </label>
        <div>
            <textarea rows=5 cols=30 name="{{ $property }}" id="{{ $property }}" placeholder="{{ $val.Value.Description }}">{}</textarea>
        </div>
		{{- end }}
		{{- end }}
        <button>Submit</button>
      </form>
    </div>
    <div id="response"></div>
  </body>
  <script src="index.js"></script>
</html>
`

const webJSServiceTemplate = `
{{ $service := .service -}}
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
		var form = document.getElementById("{{ .endpoint }}");
		var token = form.elements["token"].value;
		var service = form.elements["service"].value;
		var endpoint = form.elements["endpoint"].value;
		{{ range $property, $val := .properties -}}
		var {{ $property }} = form.elements["{{ $property }}"].value;
		{{ end }}
  
		var m3o = new Client(token);
  
		m3o.call({{ $service.Name }}, {{ .endpoint }}, request, function(response) {
		  document.getElementById("response").innerText = response;
		});
  }
`
