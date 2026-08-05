{
  "storage": {
    "rootDirectory": "/data/registry",
    "gc": {{ .storage.gc }},
    "gcDelay": "{{ .storage.gc_delay }}",
    "dedupe": {{ .storage.dedupe }}{{ if .storage.retention.enabled }},
    "retention": {
      "dryRun": false,
      "delay": "{{ .storage.retention.delay }}",
      "policies": [{{ range $i, $p := .storage.retention.policies }}{{ if $i }}, {{ end }}
        {
          "repositories": [{{ range $j, $r := $p.repositories }}{{ if $j }}, {{ end }}"{{ $r }}"{{ end }}],
          "deleteReferrers": {{ if $p.delete_referrers }}true{{ else }}false{{ end }},
          "deleteUntagged": {{ if $p.delete_untagged }}true{{ else }}false{{ end }},
          "keepTags": [{{ range $k, $t := $p.keep_tags }}{{ if $k }}, {{ end }}
            {
              "patterns": [{{ range $l, $pt := $t.patterns }}{{ if $l }}, {{ end }}"{{ $pt }}"{{ end }}]{{ if $t.pulled_within }},
              "pulledWithin": "{{ $t.pulled_within }}"{{ end }}{{ if $t.pushed_within }},
              "pushedWithin": "{{ $t.pushed_within }}"{{ end }}{{ if $t.most_recently_pushed_count }},
              "mostRecentlyPushedCount": {{ $t.most_recently_pushed_count }}{{ end }}
            }{{ end }}
          ]
        }{{ end }}
      ]
    }{{ end }}
  },
  "http": {
    "address": "0.0.0.0",
    "port": "5000",
    "realm": "zot",
    "externalUrl": "{{ .external_url }}",
    "compat": ["docker2s2"],
    "readTimeout": "{{ if .blob_upload_timeout }}{{ .blob_upload_timeout }}{{ else }}2m{{ end }}",
    "writeTimeout": "{{ if .blob_upload_timeout }}{{ .blob_upload_timeout }}{{ else }}2m{{ end }}",
    "auth": {
      "openid": {
        "providers": {
          "oidc": {
            "name": "{{ .oidc.provider_name }}",
            "clientid": "{{ .oidc.client_id }}",
            "clientsecret": "{{ .oidc.client_secret }}",
            "issuer": "{{ .oidc.issuer }}",
            "scopes": ["openid", "profile", "email", "groups"],
            "claimMapping": {
              "username": "{{ .oidc.username_claim }}",
              "groups": "groups"
            }
          }
        }
      },
      "apiKey": true,
      "sessionKeysFile": "/data/session-keys.json"
    },
    {{ if .rate_limit.enabled }}
    "ratelimit": {
      "rate": {{ .rate_limit.requests_per_second }}{{ if .rate_limit.per_method }},
      "methods": [{{ range $i, $m := .rate_limit.per_method }}{{ if $i }}, {{ end }}
        {
          "method": "{{ $m.method }}",
          "rate": {{ $m.requests_per_second }}
        }{{ end }}
      ]{{ end }}
    },
    {{ end }}
    "accessControl": {
      "repositories": {
        "**": {
          "policies": [
            {
              "users": [{{ range $i, $u := .oidc.admin_identities }}{{ if $i }}, {{ end }}"{{ $u }}"{{ end }}],
              "groups": [{{ range $i, $g := .oidc.admin_groups }}{{ if $i }}, {{ end }}"{{ $g }}"{{ end }}],
              "actions": ["read", "create", "update", "delete"]
            }{{ if .readonly_identities }},
            {
              "users": [{{ range $i, $u := .readonly_identities }}{{ if $i }}, {{ end }}"{{ $u }}"{{ end }}],
              "actions": ["read"]
            }{{ end }}{{ if .readonly_groups }},
            {
              "groups": [{{ range $i, $g := .readonly_groups }}{{ if $i }}, {{ end }}"{{ $g }}"{{ end }}],
              "actions": ["read"]
            }{{ end }}
          ],
          "defaultPolicy": [],
          "anonymousPolicy": [{{ if .anonymous_read }}"read"{{ end }}]
        }{{ range .repositories }},
        "{{ .pattern }}": {
          "policies": [{{ $first := true }}{{ if .readonly_identities }}
            {
              "users": [{{ range $i, $u := .readonly_identities }}{{ if $i }}, {{ end }}"{{ $u }}"{{ end }}],
              "actions": ["read"]
            }{{ $first = false }}
          {{ end }}{{ if .readonly_groups }}{{ if not $first }},{{ end }}
            {
              "groups": [{{ range $i, $g := .readonly_groups }}{{ if $i }}, {{ end }}"{{ $g }}"{{ end }}],
              "actions": ["read"]
            }
          {{ end }}],
          "defaultPolicy": [],
          "anonymousPolicy": [{{ if .anonymous_read }}"read"{{ end }}]
        }{{ end }}
      },
      "adminPolicy": {
        "users": [{{ range $i, $u := .oidc.admin_identities }}{{ if $i }}, {{ end }}"{{ $u }}"{{ end }}],
        "groups": [{{ range $i, $g := .oidc.admin_groups }}{{ if $i }}, {{ end }}"{{ $g }}"{{ end }}],
        "actions": ["read", "create", "update", "delete"]
      }
    }
  },
  "log": {
    "level": "{{ if .debug }}debug{{ else }}info{{ end }}"
  },
  "extensions": {
    "search": {
      "enable": true{{ if .cve.enabled }},
      "cve": {
        "updateInterval": "{{ .cve.update_interval }}"
      }{{ end }}
    },
    "ui": {
      "enable": true
    }{{ if .metrics_enabled }},
    "metrics": {
      "enable": true
    }{{ end }}{{ if .webhooks }},
    "events": {
      "enable": true,
      "sinks": [{{ range $i, $w := .webhooks }}{{ if $i }}, {{ end }}
        {
          "type": "http",
          "address": "{{ $w.url }}"{{ if $w.username }},
          "username": "{{ $w.username }}",
          "password": "{{ $w.password }}"{{ end }}{{ if $w.token }},
          "token": "{{ $w.token }}"{{ end }}
        }{{ end }}
      ]
    }{{ end }}{{ if .lint.enabled }},
    "lint": {
      "enable": true,
      "mandatoryAnnotations": [{{ range $i, $a := .lint.mandatory_annotations }}{{ if $i }}, {{ end }}"{{ $a }}"{{ end }}]
    }{{ end }}{{ if .scrub.enabled }},
    "scrub": {
      "enable": true,
      "interval": "{{ .scrub.interval }}"
    }{{ end }}{{ if .sync }},
    "sync": {
      "enable": true,
      "credentialsFile": "/etc/zot/sync-credentials.json",
      "registries": [{{ range $i, $r := .sync }}{{ if $i }}, {{ end }}
        {
          "urls": ["{{ $r.url }}"],
          "onDemand": {{ if $r.on_demand }}true{{ else }}false{{ end }},
          "tlsVerify": true{{ if $r.poll_interval }},
          "pollInterval": "{{ $r.poll_interval }}"{{ end }}{{ if $r.content_prefix }},
          "content": [
            {
              "prefix": "{{ $r.content_prefix }}"{{ if $r.content_destination }},
              "destination": "{{ $r.content_destination }}"{{ end }}
            }
          ]{{ end }}
        }{{ end }}
      ]
    }{{ end }}
  }
}
