package resources

import (
	"bytes"
	"context"
	"testing"

	integreatlyv1alpha1 "github.com/integr8ly/integreatly-operator/pkg/apis/integreatly/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func getBuildScheme() (*runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	err := corev1.SchemeBuilder.AddToScheme(scheme)
	return scheme, err
}

func TestCopyDefaultPullSecretToNameSpace(t *testing.T) {
	defPullSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      DefaultOriginPullSecretName,
			Namespace: DefaultOriginPullSecretNamespace,
		},
		Data: map[string][]byte{
			"test": {'t', 'e', 's', 't'},
		},
	}

	scheme, err := getBuildScheme()
	if err != nil {
		t.Fatalf("failed to build scheme: %s", err.Error())
	}

	scenarios := []struct {
		Name         string
		FakeClient   k8sclient.Client
		Installation *integreatlyv1alpha1.RHMI
		Verify       func(client k8sclient.Client, err error, t *testing.T)
	}{
		{
			Name: "Test Default Pull Secret is successfully copied over to target namespace",
			FakeClient: fakeclient.NewFakeClientWithScheme(scheme, defPullSecret, &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test-namespace",
					Name:      "test-namespace",
					Labels:    map[string]string{"webapp": "true"},
				},
			}),
			Installation: &integreatlyv1alpha1.RHMI{},
			Verify: func(c k8sclient.Client, err error, t *testing.T) {
				if err != nil {
					t.Fatalf("unexpected error: %s", err.Error())
				}

				s := &corev1.Secret{}
				err = c.Get(context.TODO(), k8sclient.ObjectKey{Name: "new-name-of-secret", Namespace: "test-namespace"}, s)

				if bytes.Compare(s.Data["test"], defPullSecret.Data["test"]) != 0 {
					t.Fatalf("expected data %v, but got %v", defPullSecret.Data["test"], s.Data["test"])
				}
			},
		},
		{
			Name: "Test Get Default Pull Secret error when trying to copy",
			FakeClient: fakeclient.NewFakeClientWithScheme(scheme, &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "test-namespace",
					Name:      "test-namespace",
					Labels:    map[string]string{"webapp": "true"},
				},
			}),
			Installation: &integreatlyv1alpha1.RHMI{},
			Verify: func(c k8sclient.Client, err error, t *testing.T) {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
			},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			err := CopyDefaultPullSecretToNameSpace(context.TODO(), "test-namespace", "new-name-of-secret", scenario.Installation, scenario.FakeClient)
			scenario.Verify(scenario.FakeClient, err, t)
		})
	}
}
