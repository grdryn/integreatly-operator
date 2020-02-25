package resources

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/operator-framework/operator-lifecycle-manager/pkg/api/apis/operators/v1alpha1"

	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func upgradeApproval(ctx context.Context, client k8sclient.Client, ip *v1alpha1.InstallPlan) error {
	if ip.Spec.Approved == false && len(ip.Spec.ClusterServiceVersionNames) > 0 {
		logrus.Infof("Approving %s resource version: %s", ip.Name, ip.Spec.ClusterServiceVersionNames[0])
		ip.Spec.Approved = true
		err := client.Update(ctx, ip)
		if err != nil {
			return fmt.Errorf("error approving installplan: %w", err)
		}

	}
	return nil
}
