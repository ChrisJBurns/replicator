/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	mydomainv1alpha1 "replicator/api/v1alpha1"
)

// ImageReconciler reconciles a Image object
type ImageReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=my.domain,resources=images,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=my.domain,resources=images/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=my.domain,resources=images/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Image object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *ImageReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("Image", req.NamespacedName)

	// Fetch the Image instance
	instance := &mydomainv1alpha1.Image{}
	err := r.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if this Image already exists
	result, err := r.ensureImage(req, instance)
	if result != nil {

		keychain := &InMemoryKeyChain{
			credentials: map[string]userCredentials{
				"ghcr.io": {
					username: "insert username here",
					password: "insert password here",
				},
			},
		}

		log.Info("Image doesn't exist in destination")
		log.Info("Downloading image...")
		ref, err := name.ParseReference(instance.Spec.ImageSource.URL)
		img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
		check(err, log)
		log.Info("Downloaded image")

		log.Info("Pushing image")
		parts := strings.Split(instance.Spec.ImageSource.URL, ":")
		tag, err := name.ParseReference(instance.Spec.ImageDestination.URL + ":" + parts[1])
		err = remote.Write(tag, img, remote.WithAuthFromKeychain(keychain))
		check(err, log)
		log.Info("Pushed image")
		return ctrl.Result{}, nil
	}

	// // Check if this Service already exists
	// result, err = r.ensureService(req, instance, r.backendService(instance))
	// if result != nil {
	// 	log.Error(err, "Service Not ready")
	// 	return *result, err
	// }

	// Deployment and Service already exists - don't requeue
	log.Info("Skip reconcile: Image already exists")

	return ctrl.Result{}, nil
}

// struct to hold username and password
type userCredentials struct {
	username string
	password string
}

// A keychain that is able to return the creds based on the registry url
type InMemoryKeyChain struct {
	// An in-memory map of username/passwords
	credentials map[string]userCredentials
}

// Returns a function that is able to produce an AuthConfig struct. This is basically a factory factory
func (k *InMemoryKeyChain) Resolve(resource authn.Resource) (authn.Authenticator, error) {

	// Ask for the registry URL
	registryHost := resource.RegistryStr()

	// Find the user credentials by the registry URL
	userCreds, ok := k.credentials[registryHost]

	if !ok {
		return authn.Anonymous, fmt.Errorf("unable to find credentials for host %s", registryHost)
	}

	// Return an authConfig that contains the username and password we looked up
	return authn.FromConfig(authn.AuthConfig{
		Username: userCreds.username,
		Password: userCreds.password,
	}), nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ImageReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mydomainv1alpha1.Image{}).
		Complete(r)
}

func check(e error, log logr.Logger) {
	if e != nil {
		log.Error(e, "Error")
	}
}
