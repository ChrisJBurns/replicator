package controller

import (
	"log"
	mydomainv1alpha1 "replicator/api/v1alpha1"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (r *ImageReconciler) ensureImage(request reconcile.Request,
	instance *mydomainv1alpha1.Image,
) (*reconcile.Result, error) {

	// we take the version tag (and digest) of the image to check if it exists in destination
	parts := strings.Split(instance.Spec.ImageSource.URL, ":")
	ref, err := name.ParseReference(instance.Spec.ImageDestination.URL + ":" + parts[1])
	log.Print(ref.Name())
	_, err = remote.Get(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return &reconcile.Result{}, err
	}

	return nil, nil
}
