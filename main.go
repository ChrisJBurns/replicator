package main

import (
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"gopkg.in/yaml.v3"
)

type Image struct {
	Name     string `yaml:"name"`
	Url      string `yaml:"url"`
	Tag      string `yaml:"tag"`
	Digest   string `yaml:"digest"`
	Cosigned struct {
		Enabled   string `yaml:"enabled"`
		Signature string `yaml:"signature"`
	}
	PushLocation struct {
		Registry string `yaml:"registry"`
		Path     string `yaml:"path"`
	} `yaml:"push-location"`
}

type Registry struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

// correctly populate the data.
type T struct {
	Images     []Image    `yaml:"images"`
	Registries []Registry `yaml:"registries"`
}

func buildImageUrl(imageUrl string, imageTag string, imageDigest string) string {
	if imageDigest == "" {
		return fmt.Sprintf("%s:%s", imageUrl, imageTag)
	}
	return fmt.Sprintf("%s:%s@%s", imageUrl, imageTag, imageDigest)
}

func main() {
	t := readFile()

	for _, image := range t.Images {
		imageUrl := buildImageUrl(image.Url, image.Tag, image.Digest)
		ref, err := name.ParseReference(imageUrl)
		check(err)

		fmt.Printf("Downloading image: %s... \n", image.Name)
		img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
		check(err)
		fmt.Printf("Downloaded image: %s \n", image.Name)
		// fmt.Printf("image details: %v \n", img)

		hash, err := img.Digest()
		check(err)

		if image.Digest == hash.String() {
			fmt.Println("digest verified")
		}

		newPlace := getRegistryUrl(t, image)
		newUrl := fmt.Sprintf("%s/%s:%s", newPlace, image.PushLocation.Path, image.Tag)
		tag, err := name.ParseReference(newUrl)
		check(err)

		fmt.Printf("Pushing image: %s ... \n", newUrl)

		remote.Write(tag, img, remote.WithAuthFromKeychain(authn.DefaultKeychain))
		fmt.Printf("Pushed image: %s \n", newUrl)
	}
}

func getRegistryUrl(t T, image Image) string {
	for _, r := range t.Registries {
		if r.Name == image.PushLocation.Registry {
			return r.Url
		}
	}
	panic(fmt.Errorf("no registries found for image push registry"))
}

func readFile() T {
	dat, err := os.ReadFile("./config.yaml")
	check(err)
	t := T{}
	err = yaml.Unmarshal(dat, &t)
	check(err)
	return t
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// unused for now
func createFile(image Image, img v1.Image) {
	fp, err := os.CreateTemp(".", fmt.Sprintf("%s-%s", image.Name, image.Tag))
	check(err)

	newTag, err := name.NewTag(fmt.Sprintf("%s:%s", image.Name, image.Tag))
	check(err)

	defer fp.Close()
	fmt.Println("Writing to tarball...")
	if err := tarball.Write(newTag, img, fp); err != nil {
		panic(fmt.Errorf("error writing image to tarball: %v", err))
	}
	fmt.Println("Written")
}
