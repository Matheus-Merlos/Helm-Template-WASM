//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
)

func render(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return map[string]interface{}{"error": "expected 3 arguments: filesJSON, releaseName, namespace"}
	}

	filesJSON := args[0].String()
	releaseName := args[1].String()
	namespace := args[2].String()

	var files map[string]string
	if err := json.Unmarshal([]byte(filesJSON), &files); err != nil {
		return map[string]interface{}{"error": "invalid files JSON: " + err.Error()}
	}

	c := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:    "chart",
			Version: "0.1.0",
		},
		Templates: []*chart.File{},
		Files:     []*chart.File{},
	}

	for path, content := range files {
		f := &chart.File{
			Name: path,
			Data: []byte(content),
		}

		switch {
		case path == "Chart.yaml" || path == "Chart.yml":
		case len(path) >= 10 && path[:10] == "templates/":
			c.Templates = append(c.Templates, f)
		default:
			c.Files = append(c.Files, f)
		}
	}

	values := chartutil.Values{}
	if v, ok := files["values.yaml"]; ok {
		vals, err := chartutil.ReadValues([]byte(v))
		if err != nil {
			return map[string]interface{}{"error": "invalid values.yaml: " + err.Error()}
		}
		values = vals
	}

	options := chartutil.ReleaseOptions{
		Name:      releaseName,
		Namespace: namespace,
		IsInstall: true,
		IsUpgrade: false,
	}

	vals, err := chartutil.ToRenderValues(c, values, options, nil)
	if err != nil {
		return map[string]interface{}{"error": "ToRenderValues: " + err.Error()}
	}

	rendered, err := engine.Render(c, vals)
	if err != nil {
		return map[string]interface{}{"error": "render: " + err.Error()}
	}

	out, _ := json.Marshal(rendered)
	return string(out)
}

func main() {
	js.Global().Set("helmTemplate", js.FuncOf(render))

	select {}
}
