package graph

import (
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const labelTenant = "paas.dcas.dev/tenant"

var ErrUnauthorised = errors.New("unauthorised")

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client.Client
	Scheme *runtime.Scheme
}
