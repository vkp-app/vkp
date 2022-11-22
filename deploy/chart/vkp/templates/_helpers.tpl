{{- define "vkp.labels" -}}
app.kubernetes.io/name: vkp
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "dex.labels" }}
app.kubernetes.io/component: dex
{{ include "vkp.labels" . }}
{{- end }}

{{- define "api.labels" }}
app.kubernetes.io/component: apiserver
{{ include "vkp.labels" . }}
{{- end }}

{{- define "web.labels" }}
app.kubernetes.io/component: web
{{ include "vkp.labels" . }}
{{- end }}

{{- define "pki.caPath" -}}
/var/run/secrets/paas.dcas.dev/tls:/etc/ssl/certs:/etc/pki/tls/certs:/etc/pki/tls:/etc/pki/ca-trust/extracted:/etc/ssl
{{- end }}

{{- define "dex.secret" }}
{{- printf "%s-idp" .Release.Name }}
{{- end }}