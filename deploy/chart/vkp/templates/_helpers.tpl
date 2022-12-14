{{- define "vkp.labels" -}}
app.kubernetes.io/name: vkp
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "dex.labels" }}
app.kubernetes.io/component: dex
{{ include "vkp.labels" . }}
{{- end }}

{{- define "metrics-proxy.labels" }}
app.kubernetes.io/component: metrics-proxy
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

{{- define "image" }}
{{- if (.global.imageRegistry | default .image.registry) }}
{{- printf "%s/%s:%s" (.global.imageRegistry | default .image.registry) .image.repository (.image.tag | default .chart.AppVersion) }}
{{- else }}
{{- printf "%s:%s" .image.repository (.image.tag | default .chart.AppVersion) }}
{{- end }}
{{- end }}

{{- define "console.host" }}
{{- printf "%s.%s" .Values.vkp.consoleHost .Values.global.ingress.domain }}
{{- end }}


{{- define "dex.host" }}
{{- printf "%s.%s" .Values.dex.ingress.host .Values.global.ingress.domain }}
{{- end }}
