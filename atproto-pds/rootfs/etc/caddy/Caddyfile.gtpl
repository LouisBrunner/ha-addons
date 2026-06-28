{
  http_port 3000
  auto_https off
  admin off
  default_bind 0.0.0.0
}

*.{{ .hostname }}:3000, {{ .hostname }}:3000 {
  log {
  }

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
    reverse_proxy http://localhost:3002 {
      header_up X-Forwarded-For {http.request.header.CF-Connecting-IP}
    }
  }

  reverse_proxy http://localhost:3001
}

:3000 {
  log {
  }
  respond 421
}
