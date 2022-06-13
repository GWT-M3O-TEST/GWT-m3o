package main

// const dartIndexTemplate = `library m3o;

// export 'src/client.dart';
// {{ range $service := .services }}export 'src/{{ $service.Name }}/{{ $service.Name}}.dart';
// {{ end }}
// `

// class {{title $service.Name}}Service {
// 	var _client;
//   	final String token;

// 	{{title $service.Name}}Service(String token) :token = token {
// 	  _client = Client(token: token);
// 	}
// {{ range $key, $req := $service.Spec.Components.RequestBodies }}{{ $reqType := requestType $key }}{{ $endpointName := requestTypeToEndpointName $key}}
// 	/{{ if endpointComment $endpointName $service.Spec.Components.Schemas }}{{ endpointComment $endpointName $service.Spec.Components.Schemas }}{{ end -}}
// 	{{ if isNotStream $service.Spec $service.Name $reqType }}Future<{{ $endpointName }}Response> {{untitle $endpointName}}({{ $endpointName }}Request req) async {
// 		Request request = Request(
// 			service: '{{$service.Name}}',
// 			endpoint: '{{$endpointName}}',
// 			body: req.toJson(),
// 		);

// 		try {
// 			Response res = await _client.call(request);
// 			if (isError(res.body)) {
// 			  final err = Merr(res.toJson());
// 			  return {{ $endpointName }}Response.Merr(body: err.b);
// 			}
// 			return {{ $endpointName }}ResponseData.fromJson(res.body);
// 		  } catch (e) {
// 			throw Exception(e);
// 		  }
// 	}{{end}}
// 	{{ if isStream $service.Spec $service.Name $reqType }}Stream<{{ $endpointName }}Response> {{untitle $endpointName}}({{ $endpointName }}Request req) async* {
// 		Request request = Request(
// 			service: '{{$service.Name}}',
// 			endpoint: '{{$endpointName}}',
// 			body: req.toJson(),
// 		);

// 		try {
// 			var webS = await _client.stream(request);
// 			await for (var value in webS!) {
// 				final vo = jsonDecode(value);
// 				if (isError(vo)) {
// 					yield {{ $endpointName }}Response.Merr(body: vo);
// 				} else {
// 					yield {{ $endpointName }}ResponseData.fromJson(vo);
// 				}
// 			}
// 		} catch (e) {
// 			throw Exception(e);
// 		}
// 	}{{end}}{{end}}
// }

// {{ range $typeName, $schema := $service.Spec.Components.Schemas }}
// {{ $isResponse := isResponse $typeName }}
// {{ if not $isResponse }}
// @Freezed()
// class {{ title $typeName }} with _${{ title $typeName }} {
// 	const factory {{ title $typeName }}({{ recursiveTypeDefinitionDart $service.Name $typeName $service.Spec.Components.Schemas }}) = _{{ title $typeName }};
// 	factory {{ title $typeName }}.fromJson(Map<String, dynamic> json) =>
//       _${{ title $typeName }}FromJson(json);
// }
// {{ end }}
// {{ if $isResponse }}
// @Freezed()
// class {{ title $typeName }} with _${{ title $typeName }} {
// 	const factory {{ title $typeName }}({{ recursiveTypeDefinitionDart $service.Name $typeName $service.Spec.Components.Schemas }}) = {{ title $typeName }}Data;
// 	const factory {{ title $typeName }}.Merr({Map<String, dynamic>? body}) =
// 	{{ title $typeName }}Merr;
// 	factory {{ title $typeName }}.fromJson(Map<String, dynamic> json) =>
//       _${{ title $typeName }}FromJson(json);
// }
// {{ end }}
// {{ end }}

const kotlinServiceTemplate = `
{{- $service := .service }}
package com.m3o.m3okotlin.services

import com.m3o.m3okotlin.M3O.getUrl
import com.m3o.m3okotlin.M3O.ktorHttpClient
{{- if serviceHasStream $service.Spec $service.Name }}
import com.m3o.m3okotlin.WebSocket
{{- end }}

import io.ktor.client.request.*
import kotlinx.serialization.decodeFromString
import kotlinx.serialization.encodeToString
import kotlinx.serialization.json.Json
import kotlinx.serialization.Serializable

private const val SERVICE = "{{ $service.Name }}"

object {{ title $service.Name }}Service {
{{- range $key, $req := $service.Spec.Components.RequestBodies }}
{{- $reqType := requestType $key }}
{{- $endpointName := requestTypeToEndpointName $key}}
{{- if isNotStream $service.Spec $service.Name $reqType }}
    suspend fun {{ untitle $endpointName }}(req: {{ title $service.Name}}{{ $endpointName }}Request): {{ title $service.Name}}{{ $endpointName }}Response {
        return ktorHttpClient.post(getUrl(SERVICE, "{{ $endpointName }}")) {
          body = req
        }
    }
{{- end }}
{{- if isStream $service.Spec $service.Name $reqType }}
    fun {{ untitle $endpointName }}(req: {{ title $service.Name}}{{ $endpointName }}Request, action: (Exception?, {{ title $service.Name}}{{ $endpointName }}Response?) -> Unit) {
        val url = getUrl(SERVICE, "{{ $endpointName }}", true)
        WebSocket(url, Json.encodeToString(req)) { e, response ->
            action(e, if (response != null) Json.decodeFromString(response) else null)
        }.connect()
    }
{{- end }}
{{- end }}
}

{{- range $typeName, $schema := $service.Spec.Components.Schemas }}
{{- $isResponse := isResponse $typeName }}
{{- if not $isResponse }}
@Serializable
internal data class {{ title $service.Name}}{{ $typeName }}({{ recursiveTypeDefinitionKotlin $service.Name $typeName $service.Spec.Components.Schemas }})
{{- end }}
{{- if $isResponse }}
@Serializable
data class {{ title $service.Name}}{{ $typeName }}({{ recursiveTypeDefinitionKotlin $service.Name $typeName $service.Spec.Components.Schemas }})
{{- end }}
{{- end }}
`

const kotlinExampleTemplate = `{{ $service := .service }}import 'dart:io';

import 'package:m3o/src/{{ $service.Name }}/{{ $service.Name }}.dart';

void main() async {
  final ser = {{title $service.Name}}Service(Platform.environment['M3O_API_TOKEN']!);
 
  final payload = <String, dynamic>{{ dartExampleRequest .example.Request }};

  {{ title .endpoint }}Request req = {{ title .endpoint }}Request.fromJson(payload);

  
  try {
	  {{ $reqType := requestType .endpoint }}
	  {{ if isNotStream $service.Spec $service.Name $reqType }}
	  {{ title .endpoint }}Response res = await ser.{{ .endpoint }}(req);

    res.map((value) => print(value),
        Merr: ({{ title .endpoint }}ResponseMerr err) => print(err.body!['body']));

	  {{ end }}	
	  {{ if isStream $service.Spec $service.Name $reqType }}
	  final res = await ser.{{ .endpoint }}(req);
		await for (var sr in res) {
		sr.map((value) => print(value),
			Merr: ({{ title .endpoint }}ResponseMerr err) => print(err.body));
		}	
	  {{ end }}
  } catch (e) {
    print(e);
  } finally {
    exit(0);
  }
}`

const kotlinReadmeTopTemplate = `{{ $service := .service }}# {{ title $service.Name }}

An [m3o.com](https://m3o.com) API. For example usage see [m3o.com/{{ $service.Name }}/api](https://m3o.com/{{ $service.Name }}/api).

Endpoints:

`

const kotlinReadmeBottomTemplate = `{{ $service := .service }}## {{ title .endpoint }}

{{ endpointDescription .endpoint $service.Spec.Components.Schemas }}

[https://m3o.com/{{ $service.Name }}/api#{{ title .endpoint}}](https://m3o.com/{{ $service.Name }}/api#{{ title .endpoint}})

` + "```" + `dart
{{ $service := .service -}}import 'dart:io';

import 'package:m3o/src/{{ $service.Name }}/{{ $service.Name }}.dart';

void main() async {
  final ser = {{title $service.Name}}Service(Platform.environment['M3O_API_TOKEN']!);
 
  final payload = <String, dynamic>{{ dartExampleRequest .example.Request }};

  {{ title .endpoint }}Request req = {{ title .endpoint }}Request.fromJson(payload);

  {{ $reqType := requestType .endpoint }}
  {{ if isNotStream $service.Spec $service.Name $reqType -}}
  try {

	{{ title .endpoint }}Response res = await ser.{{ .endpoint }}(req);

    res.map((value) => print(value),
	  Merr: ({{ title .endpoint }}ResponseMerr err) => print(err.body!['body']));
  {{- end }}	
  {{ if isStream $service.Spec $service.Name $reqType -}}
  try {

    final res = await ser.{{ .endpoint }}(req);

	  await for (var sr in res) {
	  sr.map((value) => print(value),
		Merr: ({{ title .endpoint }}ResponseMerr err) => print(err.body));
	  }	
	{{- end }}
  } catch (e) {
    print(e);
  } finally {
    exit(0);
  }
}
` + "```" + `
`
