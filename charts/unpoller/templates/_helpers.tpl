{{/*
Expand the name of the chart.
*/}}
{{- define "unpoller.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "unpoller.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Name for the grafana CR secret.
*/}}
{{- define "unpoller.grafana-secret" -}}
{{-  printf "%s-%s"  "grafana-secret" (include "unpoller.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}


{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "unpoller.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "unpoller.labels" -}}
helm.sh/chart: {{ include "unpoller.chart" . }}
{{ include "unpoller.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "unpoller.selectorLabels" -}}
app.kubernetes.io/name: {{ include "unpoller.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{- define "unpoller.dashboadSelectorLabels" -}}
{{- if .Values.dashboards.grafana.create -}}
app.kubernetes.io/name: {{ include "unpoller.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- else -}}
{{- with .Values.dashboards.grafana.selectorLabels }}
{{- toYaml . }}
{{- end }}
{{- end -}}
{{- end}}

{{/*
Create the name of the service account to use
*/}}
{{- define "unpoller.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "unpoller.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
