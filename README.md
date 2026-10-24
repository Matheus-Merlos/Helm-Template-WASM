<h1 align="center"> Helm Template in WASM </h1>
<p align="center">
  <a href="https://forthebadge.com">
    <img src="https://forthebadge.com/badges/made-with-go.svg" alt="Badge 1" />
  </a>
  <a href="https://forthebadge.com">
    <img src="https://forthebadge.com/badges/powered-by-black-magic.svg" alt="Badge 2" />
  </a>
</p>

<p align="center">
  Render Helm charts in the browser using WebAssembly.<br>
  No server required. No CLI required.
</p>

---

### Description

This project compiles Helm's template engine (`helm.sh/helm/v3/pkg/engine`) to WebAssembly so you can run `helm template` entirely client-side.

It is useful for:

- Building Helm chart playgrounds / editors;
- Previewing rendered manifests without installing Helm.
- Integrating chart rendering into web applications.

---

### Technologies

- **Go** (compiled to WebAssembly)
- **Helm v3** (`helm.sh/helm/v3`)
- WebAssembly (`GOOS=js GOARCH=wasm`)

> **Note:** Currently only **Helm v3** is supported.

---

### How to Build

Requirements:
- Go 1.21+

```bash
# Clone the repo
git clone https://github.com/Matheus-Merlos/Helm-Template-WASM.git
cd Helm-Template-WASM

# Build the WASM module
make build

# (Optional) Also copy wasm_exec.js
make release
```


The output will be in the `dist/` folder:

- `dist/helm-template.wasm`
- `dist/wasm_exec.js` (when using `make release`)

---

### How to Use

1. Serve `helm-template.wasm` and `wasm_exec.js` from your web app.
2. Load them and call the exposed function.

Example (simplified):

```js
const go = new Go();
const result = await WebAssembly.instantiateStreaming(
  fetch("/helm-template.wasm"),
  go.importObject
);
go.run(result.instance);

// files = { "templates/deployment.yaml": "...", "Chart.yaml": "..." }
const rendered = helmV3Template(
  JSON.stringify(files),
  valuesYamlString,
  "my-chart",     // optional chart name
  "0.1.0",        // optional chart version
  "my-release",   // optional release name
  "default"       // optional namespace
);

console.log(JSON.parse(rendered));
```

---

### Contributing

Contributions are welcome!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/my-contribution`)
3. Commit your changes
4. Open a Pull Request

Please make sure the project still builds with:

```bash
make build
```

---

### License

This project is licensed under the **Apache License 2.0**.

See the [LICENSE](LICENSE) file for the full text.

> This project uses code from [Helm](https://github.com/helm/helm), which is also licensed under Apache 2.0.
