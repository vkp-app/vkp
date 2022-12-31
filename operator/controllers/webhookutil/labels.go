package webhookutil

import "k8s.io/apimachinery/pkg/util/validation/field"

func RequireLabel(labels map[string]string, k, v string) *field.Error {
	return requireKV(labels, "labels", k, v)
}

func RequireAnnotation(annotations map[string]string, k, v string) *field.Error {
	return requireKV(annotations, "annotations", k, v)
}

func requireKV(data map[string]string, fieldName, k, v string) *field.Error {
	if data[k] != v {
		return field.Forbidden(field.NewPath("metadata").Child(fieldName).Key(k), "system fields cannot be changed")
	}
	return nil
}
