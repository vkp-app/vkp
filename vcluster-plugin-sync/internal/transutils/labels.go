package transutils

import (
	"github.com/loft-sh/vcluster-sdk/syncer/translator"
	"github.com/loft-sh/vcluster-sdk/translate"
)

func TranslateLabels(labels map[string]string, namespace string) map[string]string {
	newLabels := map[string]string{}
	for k, v := range labels {
		if k == translate.NamespaceLabel {
			newLabels[k] = v
			continue
		}

		newLabels[translator.ConvertLabelKey(k)] = v
	}

	newLabels[translate.MarkerLabel] = translate.Suffix
	if namespace != "" {
		newLabels[translate.NamespaceLabel] = namespace
	}
	return newLabels
}
