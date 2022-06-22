package main

const kotlinServiceTemplate = `
{{- $service := .service }}
package com.m3o.m3okotlin.services.{{ $service.Name }}

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
  {{- $isEmptyRequest := isEmptyRequest $reqType $service.Spec.Components.Schemas}}
  {{- $isEmptyResponse := isEmptyResponse $reqType $service.Spec.Components.Schemas}}
    {{- if $isEmptyRequest }}
    suspend fun {{ untitle $endpointName }}(): {{ title $service.Name}}{{ $endpointName }}Response {
        return ktorHttpClient.post(getUrl(SERVICE, "{{ $endpointName }}")) 
    }
    {{- else if $isEmptyResponse }}
    suspend fun {{ untitle $endpointName }}(req: {{ title $service.Name}}{{ $endpointName }}Request){
      return ktorHttpClient.post(getUrl(SERVICE, "{{ $endpointName }}")) {
        body = req
      }
    }  
    {{- else if not $isEmptyRequest }}
    suspend fun {{ untitle $endpointName }}(req: {{ title $service.Name}}{{ $endpointName }}Request): {{ title $service.Name}}{{ $endpointName }}Response {
        return ktorHttpClient.post(getUrl(SERVICE, "{{ $endpointName }}")) {
          body = req
        }
    }
    {{- end }}
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
  {{- if isObject $typeName }}
@Serializable
data class {{ title $service.Name}}{{ title $typeName }}({{ recursiveTypeDefinitionKotlin $service.Name $typeName $service.Spec.Components.Schemas }})
  {{- else if (isEmptyRequest $typeName $service.Spec.Components.Schemas) }}
@Serializable
class {{ title $service.Name}}{{ title $typeName }}({{ recursiveTypeDefinitionKotlin $service.Name $typeName $service.Spec.Components.Schemas }})
  {{- else if (isEmptyResponse $typeName $service.Spec.Components.Schemas) }}
@Serializable
class {{ title $service.Name}}{{ title $typeName }}({{ recursiveTypeDefinitionKotlin $service.Name $typeName $service.Spec.Components.Schemas }})
  {{- else }}
@Serializable
data class {{ title $service.Name}}{{ title $typeName }}({{ recursiveTypeDefinitionKotlin $service.Name $typeName $service.Spec.Components.Schemas }})  
  {{- end }}
{{- end }}
`

const kotlinExampleTemplate = `
{{- $service := .service }}
{{- $endpoint := .endpoint }}
package examples.{{ $service.Name }}.{{ $endpoint }}

import com.m3o.m3okotlin.M3O
import com.m3o.m3okotlin.services.{{ $service.Name }}

{{- $reqType := requestType $endpoint }}
{{- if isNotStream $service.Spec $service.Name $reqType }}

suspend fun main() {
  M3O.initialize(System.getenv("M3O_API_TOKEN"))

  val req = {{ title $service.Name }}{{ title $endpoint }}Request(name = "Jone")
  
  try {
      val response = {{ title $service.Name }}Service.{{ $endpoint }}(req)
      println(response)
  } catch (e: Exception) {
      println(e)
  }
}
{{- end }}
{{- if isStream $service.Spec $service.Name $reqType }}
fun main() {
  M3O.initialize(System.getenv("M3O_API_TOKEN"))

  val req = val req = {{ title $service.Name }}{{ title $endpoint }}Request(messages = 2, name = "John")
  
  try {
      val socket = {{ title $service.Name }}Service.{{ $endpoint }}(req) { socketError, response ->
          if (socketError == null) {
              println(response)
          } else {
              println(socketError)
          }
      }
  } catch (e: Exception) {
      println(e)
  }
}
{{- end }}
`

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
