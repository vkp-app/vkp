package model

import paasv1alpha1 "gitlab.dcas.dev/k8s/kube-glass/operator/api/v1alpha1"

func FromDAO(t paasv1alpha1.ReleaseTrack) Track {
	switch t {
	case paasv1alpha1.TrackStable:
		return TrackStable
	default:
		fallthrough
	case paasv1alpha1.TrackRegular:
		return TrackRegular
	case paasv1alpha1.TrackRapid:
		return TrackRapid
	case paasv1alpha1.TrackBeta:
		return TrackBeta
	}
}

func (e Track) ToDAO() paasv1alpha1.ReleaseTrack {
	switch e {
	case TrackStable:
		return paasv1alpha1.TrackStable
	default:
		fallthrough
	case TrackRegular:
		return paasv1alpha1.TrackRegular
	case TrackRapid:
		return paasv1alpha1.TrackRapid
	case TrackBeta:
		return paasv1alpha1.TrackBeta
	}
}
