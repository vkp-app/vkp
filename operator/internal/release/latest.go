package release

import (
	"context"
	"fmt"
	"gitlab.dcas.dev/k8s/kube-glass/operator/apis/paas/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logging "sigs.k8s.io/controller-runtime/pkg/log"
)

func GetLatest(ctx context.Context, c client.Client, track v1alpha1.ReleaseTrack) (*v1alpha1.ClusterVersion, error) {
	log := logging.FromContext(ctx).WithValues("track", track)
	log.Info("fetching latest cluster version")

	versionList := &v1alpha1.ClusterVersionList{}
	if err := c.List(ctx, versionList, client.MatchingLabels{v1alpha1.LabelTrackRef: string(track)}); err != nil {
		log.Error(err, "failed to list cluster versions")
		return nil, err
	}
	if len(versionList.Items) == 0 {
		return nil, fmt.Errorf("failed to find cluster versions for track: %s", track)
	}
	// find the largest minor version
	maxMinor := versionList.Items[0]
	for _, v := range versionList.Items {
		if v.Status.VersionNumber.Minor > maxMinor.Status.VersionNumber.Minor {
			maxMinor = v
			continue
		}
	}
	// find the largest patch version of
	// the largest minor version
	max := maxMinor
	for _, v := range versionList.Items {
		if v.Status.VersionNumber.Minor != maxMinor.Status.VersionNumber.Minor {
			continue
		}
		if v.Status.VersionNumber.Patch > max.Status.VersionNumber.Patch {
			max = v
		}
	}
	log.Info("located latest version", "version", max.Status.VersionNumber)
	return &max, nil
}
