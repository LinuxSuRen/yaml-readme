{{- range $key, $val := .}}
Year: {{$key}}
| Zh | En |
|---|---|
{{- range $item := $val}}
| {{$item.zh}} | {{$item.en}} |
{{- end}}
{{end}}