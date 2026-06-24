{{/*
Common helpers for the basic-chart fixture.
*/}}
{{- define "basic-chart.fullname" -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
