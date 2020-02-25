package client

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	fakesigs "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

//go:generate moq -out sigs_client_moq.go . SigsClientInterface
type SigsClientInterface interface {
	k8sclient.Reader
	k8sclient.Writer
	k8sclient.StatusClient
	GetSigsClient() k8sclient.Client
}

func NewSigsClientMoqWithScheme(clientScheme *runtime.Scheme, initObjs ...runtime.Object) *SigsClientInterfaceMock {
	sigsClient := fakesigs.NewFakeClientWithScheme(clientScheme, initObjs...)
	return &SigsClientInterfaceMock{
		GetSigsClientFunc: func() k8sclient.Client {
			return sigsClient
		},
		GetFunc: func(ctx context.Context, key k8sclient.ObjectKey, obj runtime.Object) error {
			return sigsClient.Get(ctx, key, obj)
		},
		CreateFunc: func(ctx context.Context, obj runtime.Object, opts ...k8sclient.CreateOption) error {
			return sigsClient.Create(ctx, obj)
		},
		UpdateFunc: func(ctx context.Context, obj runtime.Object, opts ...k8sclient.UpdateOption) error {
			return sigsClient.Update(ctx, obj)
		},
		DeleteFunc: func(ctx context.Context, obj runtime.Object, opts ...k8sclient.DeleteOption) error {
			return sigsClient.Delete(ctx, obj)
		},
		ListFunc: func(ctx context.Context, list runtime.Object, opts ...k8sclient.ListOption) error {
			return sigsClient.List(ctx, list, opts...)
		},
		StatusFunc: func() k8sclient.StatusWriter {
			return sigsClient.Status()
		},
	}
}
