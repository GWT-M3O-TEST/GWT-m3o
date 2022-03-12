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
      <form id="{{ untitle .endpoint }}" onsubmit="{{ untitle .endpoint }}()">
        <div>
            <label for="service"><b>{{ $service.Name }}</b></label>
            <input type="hidden" name="service" id="service" value="{{ $service.Name }}">
        </div>
        <div>
            <label for="endpoint"><b>{{ .endpoint }}</b></label>
            <input type="hidden" name="endpoint" id="endpoint" value="{{ .endpoint }}">
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
  <script src="{{ untitle .endpoint }}.js"></script>
</html>
`

const webJSServiceTemplate = `
{{ $service := .service -}}
class {{ title $service.Name }} {
	constructor(token) {
	  this.token = token;
	}
  
	call(service, endpoint, request, callback) {
	  // e.g /v1/helloworld/Call
	  var path = "/v1/" + service + "/" + endpoint
  
	  var xmlHttp = new XMLHttpRequest();
	  xmlHttp.open("POST", "https://api.m3o.com" + path, true); // true for asynchronous
	  xmlHttp.setRequestHeader("Authorization", "Bearer " + this.token);
	  xmlHttp.setRequestHeader("Content-Type", "application/json");
	  
	  xmlHttp.onreadystatechange = function () {
		if(xmlHttp.readyState === 4) {
			var status = xmlHttp.status;
			if (status === 0 || (status >= 200 && status < 400)) {
			  callback(xmlHttp.responseText);
			} else {
			
			}
		}
	  };

	  xmlHttp.send(request);
	}
  }
  
  function {{ untitle .endpoint }}() {
		var token = document.getElementById("token").value;
		var service = document.getElementById("service").value;
		var endpoint = document.getElementById("endpoint").value;
		{{- range $property, $val := .properties }}
		var {{ $property }} = document.getElementById("{{ $property }}").value;
		{{- end }}
		var obj = new Object();
		{{- range $property, $val := .properties }}
		obj.{{ $property }} = {{ $property }};
		{{ end }}
		var request = JSON.stringify(obj);
  
		var m3o = new {{ title $service.Name }}(token);
  
		m3o.call(service, endpoint, request, function(response) {
		  document.getElementById("response").innerText = response;
		});
  }
`
