package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fatih/camelcase"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
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

	// loop over service's endpoints
	for serv, meta := range service.Spec.Components.Schemas {

		parts := camelcase.Split(serv)
		if parts[len(parts)-1] == "Request" {
			endpoint := strings.Join(parts[:len(parts)-1], "")
			endpointDesc := meta.Value.Description
			// fmt.Println("endpoint:", endpoint)
			// fmt.Println("description:", endpointDesc)

			// applying paresd html template to m3o services
			b_html := bytes.Buffer{}
			buf_html := bufio.NewWriter(&b_html)
			err = tempHTML.Execute(buf_html, map[string]interface{}{
				"service":    service,
				"endpoint":   endpoint,
				"epdesc":     endpointDesc,
				"properties": meta.Value.Properties,
			})
			if err != nil {
				fmt.Println("Failed to unmarshal", err)
				os.Exit(1)
			}

			// applying paresd js template to m3o services
			b_js := bytes.Buffer{}
			buf_js := bufio.NewWriter(&b_js)
			err = tempJS.Execute(buf_js, map[string]interface{}{
				"service":    service,
				"endpoint":   endpoint,
				"epdesc":     endpointDesc,
				"properties": meta.Value.Properties,
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
	// curl example
	templ, err := template.New("cli" + serviceName + endpoint).Funcs(funcMap()).Parse(cliExampleTemplate)
	if err != nil {
		fmt.Println("Failed to unmarshal", err)
		os.Exit(1)
	}
	b := bytes.Buffer{}
	buf := bufio.NewWriter(&b)
	err = templ.Execute(buf, map[string]interface{}{
		"service":  service,
		"example":  example,
		"endpoint": endpoint,
		"funcName": strcase.UpperCamelCase(title),
	})

	if err != nil {
		fmt.Println("Failed to apply a parsed template to the specified data object", err)
		os.Exit(1)
	}

	err = os.MkdirAll(filepath.Join(examplesPath, "cli", serviceName, endpoint), FOLDER_EXECUTE_PERMISSION)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cliExampleFile := filepath.Join(examplesPath, "cli", serviceName, endpoint, title+".sh")
	f, err := os.OpenFile(cliExampleFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, FILE_EXECUTE_PERMISSION)
	if err != nil {
		fmt.Println("Failed to open schema file", err)
		os.Exit(1)
	}

	buf.Flush()
	_, err = f.Write(b.Bytes())
	if err != nil {
		fmt.Println("Failed to append to schema file", err)
		os.Exit(1)
	}
}

func schemaToWebExample(exampleJSON map[string]interface{}) string {
	// type jsonObj map[string]interface{}
	s := ""
	for key, value := range exampleJSON {
		// fmt.Println(value)
		switch value.(type) {
		case float64:
			val := value.(float64)
			s += "\t--" + key + "=" + fmt.Sprint(val) + " \\\n"
		case int64:
			val := value.(int64)
			s += "\t--" + key + "=" + fmt.Sprint(val) + " \\\n"
		case string:
			s += "\t--" + key + "=" + "\"" + value.(string) + "\"" + " \\\n"
		case interface{}:
			bs, _ := json.MarshalIndent(value, "", "  ")
			jsonList := strings.Split(string(bs), "\n")
			s += "\t--" + key + "=" + "'" + jsonList[0] + "\n"
			spacer := strings.Repeat(" ", 4+len(key))
			for index, line := range jsonList[1:] {
				if index == len(jsonList)-2 {
					s += "\t" + spacer + line + "' \\\n"
				} else {
					s += "\t" + spacer + line + "\n"
				}
			}
		}
	}
	return strings.TrimSuffix(s, "\\\n") + "\n"
}
