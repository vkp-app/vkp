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