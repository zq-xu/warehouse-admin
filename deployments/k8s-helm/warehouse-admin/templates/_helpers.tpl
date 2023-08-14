{{- define "cicdWebserver.fullName" -}}
{{- default "cicd-webserver" .Values.cicdWebserver.nameOverride -}}
{{- end -}}

{{- define "namespace" -}}
{{- default "kube-system" .Release.Namespace -}}
{{- end -}}
