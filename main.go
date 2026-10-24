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
	if len(args) < 4 {
		return map[string]interface{}{"error": "expected at least 4 arguments: filesJSON, chartName, chartVersion, valuesStr"}
	}

	filesJSON := args[0].String()
	chartName := args[1].String()
	chartVersion := args[2].String()
	valuesStr := args[3].String()

	releaseName := "release-name"
	if len(args) > 4 && !args[4].IsUndefined() && !args[4].IsNull() && args[4].String() != "" {
		releaseName = args[4].String()
	}

	namespace := "default"
	if len(args) > 5 && !args[5].IsUndefined() && !args[5].IsNull() && args[5].String() != "" {
		namespace = args[5].String()
	}

	var files map[string]string
	if err := json.Unmarshal([]byte(filesJSON), &files); err != nil {
		return map[string]interface{}{"error": "invalid files JSON: " + err.Error()}
	}

	values, err := chartutil.ReadValues([]byte(valuesStr))
	if err != nil {
		return map[string]interface{}{"error": "invalid values format: " + err.Error()}
	}

	c := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:    chartName,
			Version: chartVersion,
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
