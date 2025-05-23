package siri

import (
	"bytes"
	"embed"
	"io"
	"text/template"
)

//go:embed template/*.xml.template
var templatesFS embed.FS

func GenerateStopMonitoringRequest(stops []string) (string, error) {
	templ, err := templatesFS.Open("template/StopMonitoringRequest.xml.template")
	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(templ)
	if err != nil {
		return "", err
	}

	templateData := struct {
		RequestorRef string
		Stops        []string
	}{
		RequestorRef: RequestorRef,
		Stops:        stops,
	}

	tmpl, err := template.New("StopMonitoringRequest.xml.template").Parse(string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, templateData)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
