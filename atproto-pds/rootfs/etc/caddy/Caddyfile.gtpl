{
  http_port 3000
  auto_https off
  admin off
  default_bind 0.0.0.0
}

*.{{ .hostname }}:3000, {{ .hostname }}:3000 {
  log matched {
  }

  {{- if .gatekeeper.enabled }}
  @gatekeeper {
    path /xrpc/com.atproto.server.getSession
    path /xrpc/com.atproto.server.describeServer
    path /xrpc/com.atproto.server.updateEmail
    path /xrpc/com.atproto.server.createSession
    path /xrpc/com.atproto.server.createAccount
    path /@atproto/oauth-provider/~api/sign-in
    path /gate/*
  }

  handle @gatekeeper {
    reverse_proxy http://127.0.0.1:3002 {
      header_up X-Forwarded-For {http.request.header.CF-Connecting-IP}
      header_up X-Forwarded-Proto {http.request.header.X-Forwarded-Proto}
      header_up X-Forwarded-Host {http.request.host}
      header_up X-Real-IP {http.request.header.CF-Connecting-IP}
    }
  }
  {{- end }}

  {{- if .customize.enabled }}
  import /share/{{ .customize.caddyfile_filename }};
  {{- end }}

  reverse_proxy http://127.0.0.1:3001 {
    header_up X-Forwarded-For {http.request.header.CF-Connecting-IP}
    header_up X-Forwarded-Proto {http.request.header.X-Forwarded-Proto}
    header_up X-Forwarded-Host {http.request.host}
    header_up X-Real-IP {http.request.header.CF-Connecting-IP}
  }
}

:3000 {
  log unmatched {
  }
  respond 421
}
