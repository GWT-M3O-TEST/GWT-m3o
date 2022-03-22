package main

// Range over endpoint attributes
// for property, val := range meta.Value.Properties {
// 	propDescription := val.Value.Description
// 	fmt.Println("attribute:", property)
// 	fmt.Println("placeholder:", propDescription)
// }

const webHTMLServiceTemplate = `
<!DOCTYPE html>
	<head>
	<!-- Required meta tags -->
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">

	<!-- Bootstrap CSS -->
	<link href="https://cdn.jsdelivr.net/npm/bootstrap@5.1.3/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-1BmE4kWBq78iYhFldvKuhfTAU6auU8tT94WrHftjDbrCEXSU1oBoqyl2QvZ6jIW3" crossorigin="anonymous">

	<title>M3O Web</title>
	</head>
	<body>
	<div id="{{ .service.Name }}" class="container">
		<div class="row">
    		<div class="col">
			<form id="{{ untitle .endpoint }}">
				<div class="mb-3">
					<label for="service" class="form-label fs-1 fw-bold">{{ .service.Name }}</label>
					<input type="hidden" class="form-control" id="service" value="{{ .service.Name }}">
		  		</div>
				<div class="mb-3">
					<label for="endpoint" class="form-label fs-3 fw-bold">{{ .endpoint }}</label>
				  	<input type="hidden" class="form-control" id="endpoint" value="{{ .endpoint }}">
				  	<div id="endpointDesc" class="form-text"><i>{{ endpointDescription .endpoint .schemas }}</i></div>
				</div>
				<div class="mb-3">
              		<label for="token" class="form-label fs-5">Token</label>
              		<input class="form-control" id="token">
            	</div>
				{{- range $reqp, $val := .reqps }}
				{{- if not (eq $val.Value.Type "object") }}
				<div class="mb-3">
              		<label for="{{ $reqp }}" class="form-label fs-5">{{ $reqp }}</label>
					<div class="form-text"><i>{{ $val.Value.Description }}</i></div>
              		<input class="form-control" id="{{ $reqp }}">
            	</div>
				{{- end }}
				{{- if eq $val.Value.Type "object" }}
				<div class="mb-3">
              		<label for="{{ $reqp }}" class="form-label fs-5">{{ $reqp }}</label>
					<div class="form-text"><i>{{ $val.Value.Description }}</i></div>
              		<textarea rows="4" class="form-control" id="{{ $reqp }}"></textarea>
            	</div>
				{{- end }}
				{{- end }}
				<button type="button" class="btn btn-primary" onclick="{{ .service.Name }}{{ .endpoint }}()">Submit</button>
			</form>
    		</div>
    		<div class="col-6">
				<p class="fs-1 fw-bold text-center">JSON</p>
				<div>
					<pre>
						<code id="json"></code>
					</pre>
				</div>
    		</div>
    		<div class="col">
				<p class="fs-1 fw-bold text-center">ViewTree</p>
				<div id="viewtree"></div>
    		</div>
  		</div>
  		</br>
  		<div class="row" id="restable"></div>
	</div>
	</body>
	<script type="module" src="{{ untitle .endpoint }}.js"></script>
</html>
`

const webJSServiceTemplate = `
import Client from '../../client/index.js';

window.{{ .service.Name }}{{ .endpoint }} = function () {
	let token = document.getElementById("token").value;
	let service = document.getElementById("service").value;
	let endpoint = document.getElementById("endpoint").value;
	{{- range $reqp, $val := .reqps }}
	let {{ $reqp }} = document.getElementById("{{ $reqp }}").value;
	{{- end }}
	let obj = new Object();
	{{- range $reqp, $val := .reqps }}
	obj.{{ $reqp }} = {{ $reqp }};
	{{ end }}
	let request = JSON.stringify(obj);

	let m3o = new Client(token);

	m3o.call(service, endpoint, request, function(response) {
		// resObj = JSON.parse(response);
		let res =` + "`" + `<table class="table">
		<thead>
		  <tr>
			{{- range $resp, $val := responsePropertiesFromSchemas .endpoint .schemas }}
			<th scope="col">{{ $resp }}</th>
			{{ end }}
		  </tr>
		</thead>
		<tbody>
		  <tr>
			<td>Mark</td>
			<td>Otto</td>
			<td>@mdo</td>
		  </tr>
		  <tr>
			<td>Jacob</td>
			<td>Thornton</td>
			<td>@fat</td>
		  </tr>
		  <tr>
			<td colspan="2">Larry the Bird</td>
			<td>@twitter</td>
		  </tr>
		</tbody>
	  </table>` + "`" + `
		document.getElementById("restable").innerHTML = res;
	});
}
`
