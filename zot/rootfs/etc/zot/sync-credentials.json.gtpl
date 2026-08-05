{
{{ $first := true }}{{ range .sync }}{{ if .username }}{{ if not $first }},
{{ end }}  "{{ .url }}": {
    "username": "{{ .username }}",
    "password": "{{ .password }}"
  }{{ $first = false }}{{ end }}{{ end }}
}
