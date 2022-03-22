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

	schemas := service.Spec.Components.Schemas

	for schema, meta := range schemas {

		parts := camelcase.Split(schema)
		endpoint := strings.Join(parts[:len(parts)-1], "")
		schemaSuffix := parts[len(parts)-1]
		// traversing schemas are not necessary in order, sometime the [endpoint]Response
		// comes before the [endpoint]Request in the api-[service].json file.
		if schemaSuffix == "Request" {

			// applying paresd HTML template to m3o services
			b_html := bytes.Buffer{}
			buf_html := bufio.NewWriter(&b_html)
			err = tempHTML.Execute(buf_html, map[string]interface{}{
				"service":  service,
				"schemas":  schemas,
				"endpoint": endpoint,
				"reqps":    meta.Value.Properties,
			})
			if err != nil {
				fmt.Println("Failed to unmarshal", err)
				os.Exit(1)
			}

			// applying paresd JS template to m3o services
			b_js := bytes.Buffer{}
			buf_js := bufio.NewWriter(&b_js)
			err = tempJS.Execute(buf_js, map[string]interface{}{
				"service":  service,
				"schemas":  schemas,
				"endpoint": endpoint,
				"reqps":    meta.Value.Properties,
			})
			if err != nil {
				fmt.Println("Failed to unmarshal", err)
				os.Exit(1)
			}

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

	fmt.Println("endpoint=>", endpoint)
	fmt.Println("endpointDesc=>", endPointDesc)
	for pro := range reqPro {
		fmt.Println("req=>", pro)
	}
	for pro := range resPro {
		fmt.Println("res=>", pro)
	}

}
