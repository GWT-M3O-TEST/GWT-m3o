package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/camelcase"
	"github.com/getkin/kin-openapi/openapi3"
)

type webG struct {
	generator
	// add appropriate fields
}

func (w *webG) ServiceClient(serviceName, webPath string, service service) {

	// HTML template for m3o-web clients
	tempHTML, err := template.New("webHTML" + serviceName).Funcs(funcMap()).Parse(webHTMLServiceTemplate)
	if err != nil {
		fmt.Println("Failed to unmarshal", err)
		os.Exit(1)
	}

	// JS template for m3o-web clients
	tempJS, err := template.New("webJS" + serviceName).Funcs(funcMap()).Parse(webJSServiceTemplate)
	if err != nil {
		fmt.Println("Failed to unmarshal", err)
		os.Exit(1)
	}

	// create folder for m3o service e.g clients/web/hellworld
	err = os.MkdirAll(filepath.Join(webPath, serviceName), FOLDER_EXECUTE_PERMISSION)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// loop over schemas
	for schema, meta := range service.Spec.Components.Schemas {

		parts := camelcase.Split(schema)

		// handling the request
		if parts[len(parts)-1] == "Request" {
			endpoint := strings.Join(parts[:len(parts)-1], "")
			endPointDesc := meta.Value.Description
			reqProperties := meta.Value.Properties
			ServiceClientHelper(serviceName, webPath, endpoint, endPointDesc,
				reqProperties, nil, tempHTML, tempJS, service)
		}

		// handling the response
		if parts[len(parts)-1] == "Response" {
			endpoint := strings.Join(parts[:len(parts)-1], "")
			resProperties := meta.Value.Properties
			ServiceClientHelper(serviceName, webPath, endpoint, "",
				nil, resProperties, tempHTML, tempJS, service)
		}

	}
}

func (w *webG) schemaToType(serviceName, typeName string, schemas map[string]*openapi3.SchemaRef) string {
	return ""
}

func (w *webG) IndexFile(dartPath string, services []service) {
}

func (w *webG) TopReadme(serviceName, examplesPath string, service service) {
}

func (w *webG) ExampleAndReadmeEdit(examplesPath, serviceName, endpoint, title string, service service, example example) {
}

func schemaToWebExample(exampleJSON map[string]interface{}) string {
	return ""
}

func ServiceClientHelper(serviceName, webPath, endpoint, endPointDesc string,
	reqPro, resPro map[string]*openapi3.SchemaRef,
	tempHTML, tempJS *template.Template,
	service service) {

	// applying paresd html template to m3o services
	b_html := bytes.Buffer{}
	buf_html := bufio.NewWriter(&b_html)
	err := tempHTML.Execute(buf_html, map[string]interface{}{
		"service":       service,
		"endpoint":      endpoint,
		"epdesc":        endPointDesc,
		"reqProperties": reqPro,
	})
	if err != nil {
		fmt.Println("Failed to unmarshal", err)
		os.Exit(1)
	}

	// applying paresd js template to m3o services
	b_js := bytes.Buffer{}
	buf_js := bufio.NewWriter(&b_js)
	err = tempJS.Execute(buf_js, map[string]interface{}{
		"service":       service,
		"endpoint":      endpoint,
		"epdesc":        endPointDesc,
		"reqProperties": reqPro,
		"resProperties": resPro,
	})

	if err != nil {
		fmt.Println("Failed to unmarshal", err)
		os.Exit(1)
	}

	// lower case the endpoint name
	endpoint = strings.ToLower(endpoint)

	// create folder for endpoint
	err = os.MkdirAll(filepath.Join(webPath, serviceName, endpoint), FOLDER_EXECUTE_PERMISSION)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create html file
	htmlFile := filepath.Join(webPath, serviceName, endpoint, fmt.Sprint(endpoint, ".html"))
	f, err := os.OpenFile(htmlFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, FILE_EXECUTE_PERMISSION)
	if err != nil {
		fmt.Println("Failed to open schema file", err)
		os.Exit(1)
	}
	buf_html.Flush()
	_, err = f.Write(b_html.Bytes())
	if err != nil {
		fmt.Println("Failed to append to schema file", err)
		os.Exit(1)
	}

	// create js file
	jsFile := filepath.Join(webPath, serviceName, endpoint, fmt.Sprint(endpoint, ".js"))
	f, err = os.OpenFile(jsFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, FILE_EXECUTE_PERMISSION)
	if err != nil {
		fmt.Println("Failed to open schema file", err)
		os.Exit(1)
	}
	buf_js.Flush()
	_, err = f.Write(b_js.Bytes())
	if err != nil {
		fmt.Println("Failed to append to schema file", err)
		os.Exit(1)
	}
}
