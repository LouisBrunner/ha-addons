{
	http_port 3000
	auto_https off
	admin off
	default_bind 0.0.0.0
}

{{ .hostname }}:3000 {
	log

	route {
		handle /oauth2/* {
			reverse_proxy http://127.0.0.1:4180 {
				header_up X-Forwarded-Proto {http.request.header.X-Forwarded-Proto}
			}
		}

		reverse_proxy http://127.0.0.1:4180 {
			method GET
			rewrite /oauth2/auth
			header_up X-Forwarded-Proto {http.request.header.X-Forwarded-Proto}
			header_up X-Auth-Request-Redirect https://{host}{uri}

			@accepted status 2xx
			handle_response @accepted {
				request_header +X-Auth-Request-User "{http.reverse_proxy.header.X-Auth-Request-User}"
				request_header +X-Auth-Request-Email "{http.reverse_proxy.header.X-Auth-Request-Email}"

				reverse_proxy http://127.0.0.1:7070 {
					header_up X-Forwarded-Proto {http.request.header.X-Forwarded-Proto}
				}
			}

			handle_response {
				redir * https://{{ .hostname }}/oauth2/sign_in?rd=https://{host}{uri} 302
			}
		}
	}
}

:3000 {
	handle /healthz {
		respond 200
	}

	respond 421
}
