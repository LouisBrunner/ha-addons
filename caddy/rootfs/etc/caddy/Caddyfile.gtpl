{
	email {{ .email }}
	{{ if .debug }}
	debug
	{{ end }}
	{{ if .local_certs }}
	local_certs
	{{ end }}
	admin off
	servers {
		timeouts {
			read_header 60s
			read_body 60s
			write 60s
		}
		trusted_proxies static 172.30.33.0/24
		client_ip_headers Cf-Connecting-Ip
	}
}

{{ $hardening := .hardening }}
{{ if $hardening.block_scanner_probes }}
(block_probes) {
	@probe {
		{{ range $hardening.block_scanner_probes }}
		path {{ . }}
		{{ end }}
	}
	respond @probe 403
}
{{ end }}

{{ $base_domain := .base_domain }}
{{ $max_body_size := .max_body_size_megabytes }}
{{ range .subdomains }}
{{ $body_size := $max_body_size }}
{{ if hasKey . "max_body_size_megabytes" }}{{ $body_size = .max_body_size_megabytes }}{{ end }}
{{ .name }}.{{ $base_domain }} {
	log
	{{ if $hardening.block_scanner_probes }}
	import block_probes
	{{ end }}
	{{ if $body_size }}
	request_body {
		max_size {{ $body_size }}MB
	}
	{{ end }}
	{{ if $hardening.rate_limit.events }}
	rate_limit {
		zone {{ .name }}_limit {
			key {client_ip}
			events {{ $hardening.rate_limit.events }}
			window {{ $hardening.rate_limit.window }}
		}
	}
	{{ end }}
	header {
		{{ if $hardening.hsts.max_age_seconds }}
		Strict-Transport-Security "max-age={{ $hardening.hsts.max_age_seconds | int }}{{ if $hardening.hsts.include_subdomains }}; includeSubDomains{{ end }}"
		{{ end }}
		X-Content-Type-Options "nosniff"
		Referrer-Policy "same-origin"
		{{ if $hardening.permissions_policy }}
		Permissions-Policy "{{ range $i, $f := $hardening.permissions_policy }}{{ if $i }}, {{ end }}{{ $f }}=(){{ end }}"
		{{ end }}
		Cross-Origin-Opener-Policy "same-origin"
		Cross-Origin-Embedder-Policy "credentialless"
		X-Frame-Options "DENY"
		+Content-Security-Policy "frame-ancestors 'none'"
		-X-Powered-By
		{{ if .no_store }}
		Cache-Control "no-store"
		{{ end }}
	}
	reverse_proxy {{ .upstream }}
}
{{ end }}

:8080 {
	handle /healthz {
		respond 200
	}

	respond 421
}
