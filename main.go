package main

import (
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func main() {
	ref, err := name.ParseReference("gcr.io/google-containers/pause")
	if err != nil {
		panic(err)
	}

	fmt.Println("Downloading...")
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded")

	fp, err := os.CreateTemp(".", "img")
	defer fp.Close()
	fmt.Println("Writing to tarball...")
	if err := tarball.WriteToFile(fp.Name(), ref, img); err != nil {
		fmt.Errorf("Error writing image to tarball: %v", err)
	}
	fmt.Println("Written")
}
