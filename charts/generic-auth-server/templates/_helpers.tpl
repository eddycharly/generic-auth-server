{{/* vim: set filetype=mustache: */}}

{{- define "generic-auth-server.name" -}}
{{ template "kyverno.lib.names.name" . }}
{{- end -}}

{{- define "generic-auth-server.labels" -}}
{{- template "kyverno.lib.labels.merge" (list
  (include "kyverno.lib.labels.common" .)
  (include "generic-auth-server.labels.selector" .)
) -}}
{{- end -}}

{{- define "generic-auth-server.labels.selector" -}}
{{- template "kyverno.lib.labels.merge" (list
  (include "kyverno.lib.labels.common.selector" .)
  (include "kyverno.lib.labels.component" "authz-server")
) -}}
{{- end -}}

{{- define "generic-auth-server.service-account.name" -}}
{{- if .Values.rbac.create -}}
  {{- default (include "generic-auth-server.name" .) .Values.rbac.serviceAccount.name -}}
{{- else -}}
  {{- required "A service account name is required when `rbac.create` is set to `false`" .Values.rbac.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{- define "generic-auth-server.image" -}}
{{- printf "%s/%s:%s" .registry .repository (default "latest" .tag) -}}
{{- end -}}
