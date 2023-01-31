{{/*
Expand the name of the chart.
*/}}
{{- define "transaction-card-validator.name" -}}
{{- default .Chart.Name $.Values.validator.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "transaction-card-validator.fullname" -}}
{{- if $.Values.validator.fullnameOverride }}
{{- $.Values.validator.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name $.Values.validator.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "transaction-card-validator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "transaction-card-validator.labels" -}}
helm.sh/chart: {{ include "transaction-card-validator.chart" . }}
{{ include "transaction-card-validator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "transaction-card-validator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "transaction-card-validator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "transaction-card-validator.serviceAccountName" -}}
{{- if $.Values.validator.serviceAccount.create }}
{{- default (include "transaction-card-validator.fullname" .) $.Values.validator.serviceAccount.name }}
{{- else }}
{{- default "default" $.Values.validator.serviceAccount.name }}
{{- end }}
{{- end }}
