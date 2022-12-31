package statusutil

import (
	"context"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func SetStatus(ctx context.Context, c client.Client, obj client.Object) error {
	log := logr.FromContextOrDiscard(ctx)
	if err := c.Status().Update(ctx, obj); err != nil {
		log.Error(err, "failed to update status")
		return err
	}
	return nil
}
