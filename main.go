//go:build js && wasm

package main

import (
	"encoding/json"
	"strings"
	"syscall/js"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/engine"
	"sigs.k8s.io/yaml"
)

func renderV3(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return map[string]interface{}{"error": "expected at least 2 arguments: filesJSON, valuesStr"}
	}

	filesJSON := args[0].String()
	valuesStr := args[1].String()

	var files map[string]string
	if err := json.Unmarshal([]byte(filesJSON), &files); err != nil {
		return map[string]interface{}{"error": "invalid files JSON: " + err.Error()}
	}

	// Searches the files for a Chart.yaml to extract the metadata
	var chartYamlContent string
	var chartYamlExists bool
	if content, ok := files["Chart.yaml"]; ok {
		chartYamlContent = content
		chartYamlExists = true
	} else if content, ok := files["Chart.yml"]; ok {
		chartYamlContent = content
		chartYamlExists = true
	}

	var meta chart.Metadata
	if chartYamlExists {
		if err := yaml.Unmarshal([]byte(chartYamlContent), &meta); err != nil {
			return map[string]interface{}{"error": "invalid Chart.yaml: " + err.Error()}
		}
	}

	if meta.Name == "" {
		meta.Name = "chart"
	}
	if meta.Version == "" {
		meta.Version = "0.1.0"
	}

	//Manual overrides via args
	if len(args) > 2 && !args[2].IsUndefined() && !args[2].IsNull() && args[2].String() != "" {
		meta.Name = args[2].String()
	}
	if len(args) > 3 && !args[3].IsUndefined() && !args[3].IsNull() && args[3].String() != "" {
		meta.Version = args[3].String()
	}

	releaseName := "release-name"
	if len(args) > 4 && !args[4].IsUndefined() && !args[4].IsNull() && args[4].String() != "" {
		releaseName = args[4].String()
	}

	namespace := "default"
	if len(args) > 5 && !args[5].IsUndefined() && !args[5].IsNull() && args[5].String() != "" {
		namespace = args[5].String()
	}

	values, err := chartutil.ReadValues([]byte(valuesStr))
	if err != nil {
		return map[string]interface{}{"error": "invalid values format: " + err.Error()}
	}

	c := &chart.Chart{
		Metadata:  &meta,
		Templates: []*chart.File{},
		Files:     []*chart.File{},
	}

	for path, content := range files {
		f := &chart.File{
			Name: path,
			Data: []byte(content),
		}

		switch {
		case strings.HasPrefix(path, "templates/"):
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
	js.Global().Set("helmV3Template", js.FuncOf(renderV3))

	select {}
}
