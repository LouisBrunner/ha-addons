appConfig:
  auth:
    enableOidc: true
    oidc:
      clientId: {{ .oidc.client_id }}
      endpoint: {{ .oidc.issuer }}
      scope: openid profile email{{ if .oidc.admin_group }} groups{{ end }}
{{- if .oidc.admin_group }}
      adminGroup: {{ .oidc.admin_group }}
{{- end }}
      enableSilentRenew: true
