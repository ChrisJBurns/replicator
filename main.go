package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"gopkg.in/yaml.v3"
)

type Image struct {
	Name   string `yaml:"name"`
	Source Source `yaml:"source"`
	Target Target `yaml:"target"`
	// Url      string `yaml:"url"`
	// Tag      string `yaml:"tag"`
	// Digest   string `yaml:"digest"`
	Cosigned struct {
		Enabled   string `yaml:"enabled"`
		Signature string `yaml:"signature"`
	}
	TargetLocation struct {
		Registry string `yaml:"registry"`
		Path     string `yaml:"path"`
	} `yaml:"target-location"`
}

type Source struct {
	Registry string `yaml:"registry"`
	Image    string `yaml:"image"`
}

type Target struct {
	Registry       string   `yaml:"registry"`
	Image          string   `yaml:"image"`
	AdditionalTags []string `yaml:"additionalTags"`
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

func getImageTag(image string) string {
	if strings.ContainsAny(image, "@sha256:") {
		tag := strings.Split(image, ":")[1]
		tag = strings.Split(tag, "@")[0]
		fmt.Printf("Tag is: %s \n", tag)
		return tag
	}

	tag := strings.Split(image, ":")[1]
	fmt.Printf("Tag is: %s \n", tag)
	return tag
}

func main() {
	t := readFile()

	for _, image := range t.Images {
		getImageTag(image.Source.Image)
		// imageUrl := buildImageUrl(image.Url, image.Tag, image.Digest)
		ref, err := name.ParseReference(image.Source.Image)
		check(err)

		fmt.Printf("Downloading image: %s... \n", image.Name)
		img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
		check(err)
		fmt.Printf("Downloaded image: %s \n", image.Name)
		fmt.Printf("image details: %v \n", img)

		// hash, err := img.Digest()
		// check(err)

		// if image.Digest == hash.String() {
		// 	fmt.Println("digest verified")
		// }

		// newPlace := getRegistryUrl(t, image)
		// newUrl := fmt.Sprintf("%s/%s:%s", newPlace, image.TargetLocation.Path, image.Tag)
		tag, err := name.ParseReference(image.Target.Image)
		check(err)

		fmt.Printf("Pushing image: %s ... \n", image.Target.Image)

		remote.Write(tag, img, remote.WithAuthFromKeychain(authn.DefaultKeychain))
		fmt.Printf("Pushed image: %s \n", image.Target.Image)
	}
}

func getRegistryUrl(t T, image Image) string {
	for _, r := range t.Registries {
		if r.Name == image.TargetLocation.Registry {
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
// func createFile(image Image, img v1.Image) {
// 	fp, err := os.CreateTemp(".", fmt.Sprintf("%s-%s", image.Name, image.Tag))
// 	check(err)

// 	newTag, err := name.NewTag(fmt.Sprintf("%s:%s", image.Name, image.Tag))
// 	check(err)

// 	defer fp.Close()
// 	fmt.Println("Writing to tarball...")
// 	if err := tarball.Write(newTag, img, fp); err != nil {
// 		panic(fmt.Errorf("error writing image to tarball: %v", err))
// 	}
// 	fmt.Println("Written")
// }
