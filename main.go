package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"gopkg.in/yaml.v3"
)

type Config struct {
	RegistryURL     string `yaml:"registry_url"`
	CompartmentOCID string `yaml:"compartment_ocid"`
	ImageRepository string `yaml:"image_repository"`
	ImageTag        string `yaml:"image_tag"`
	ClientID        string `yaml:"client_id"`
	ClientSecret    string `yaml:"client_secret"`
	AccessToken     string `yaml:"access_token"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
}

func main() {
	configFile, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalf("Error opening config file: %v", err)
	}
	defer configFile.Close()

	var config Config
	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("Error decoding config file: %v", err)
	}

	username := config.Username
	password := config.Password
	auth := &authn.Basic{
		Username: username,
		Password: password,
	}

	// Create a reference to the image
	imageRef := config.RegistryURL + "/" + config.CompartmentOCID + "/" + config.ImageRepository
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		log.Fatalf("Failed to parse image reference: %v", err)
	}

	ctx := context.Background()

	// listing all tags in repo and printing hash
	tags, err := remote.List(ref.Context(), remote.WithAuth(auth))
	if err != nil {
		fmt.Println("Error pulling images:", err)
		os.Exit(1)
	}
	tagsWithHashes := make(map[string]string)
	for _, tag := range tags {
		imageRef := config.RegistryURL + "/" + config.CompartmentOCID + "/" + config.ImageRepository + ":" + tag
		ref, err := name.ParseReference(imageRef)
		if err != nil {
			log.Fatalf("Failed to parse image reference: %v", err)
		}
		desc, err := remote.Get(ref, remote.WithAuth(auth), remote.WithContext(ctx))
		if err != nil {
			fmt.Println("Error pulling descriptor:", err)
			os.Exit(1)
		}
		fmt.Printf("tag: %s\ndescriptor hex: %s\n", tag, desc.Digest.Hex)
		fmt.Println("tag signed? ", isSigned(tags, desc))
		fmt.Println("tag signed by get? ", isSignedByGet(desc.Digest.Hex, auth, config))
		tagsWithHashes[tag] = desc.Digest.Hex
		fmt.Println("---------------")
	}

	fmt.Println("#" + imageRef)

	fmt.Println("---------------")

	for tag := range tagsWithHashes {
		if strings.Contains(tag, ".sig") {
			//fmt.Println("Signature:", tag, "used?", isSignatureUsed(tagsWithHashes, tag))
			if !isSignatureUsed(tagsWithHashes, tag) {
				fmt.Println(tag)
			}
		}
	}
	fmt.Println("---------------")

	// // retrieving specific tag
	// desc, err := remote.Get(ref, remote.WithAuth(auth), remote.WithContext(ctx))
	// if err != nil {
	// 	fmt.Println("Error pulling image:", err)
	// 	os.Exit(1)
	// }
	// if desc.MediaType != types.OCIManifestSchema1 && desc.MediaType != types.DockerManifestSchema2 {
	// 	log.Fatalf("Unexpected media type: %s", desc.MediaType)
	// }

	// fmt.Printf("tag: %s\ndescriptor hex: %v\n", ref.Identifier(), desc.Digest.Hex)
	// fmt.Println("tag signed? ", isSigned(tags, desc))
	// fmt.Println("tag signed by get? ", isSignedByGet(desc.Digest.Hex, auth, config))

	// // creating image from descriptor
	// img, err := desc.Image()
	// if err != nil {
	// 	log.Fatalf("getting config img: %v", err)
	// }
	// // creating manifest from image
	// manifest, err := img.Manifest()
	// if err != nil {
	// 	log.Fatalf("getting manifest: %v", err)
	// }
	// fmt.Println("Manifest annotations:")
	// for annotation, value := range manifest.Annotations {
	// 	fmt.Printf("%s = %s\n", annotation, value)
	// }

	// // Fetch the config file
	// imgConfigFile, err := img.ConfigFile()
	// if err != nil {
	// 	log.Fatalf("Error getting config file: %v", err)
	// }

	// // Print config file labels
	// fmt.Println("Config File Labels:")
	// for key, value := range imgConfigFile.Config.Labels {
	// 	fmt.Printf("%s = %s\n", key, value)
	// }

	// // print the entire config
	// configJSON, err := json.MarshalIndent(imgConfigFile, "", "  ")
	// if err != nil {
	// 	log.Fatalf("Error marshaling config file to JSON: %v", err)
	// }
	// fmt.Printf("Full Config JSON:\n%s\n", string(configJSON))

	// retrieving non-existent tag
	/*
		imageRef = config.RegistryURL + "/" + config.CompartmentOCID + "/" + config.ImageRepository + ":" + "xxx"
		ref, err = name.ParseReference(imageRef)
		if err != nil {
			log.Fatalf("Failed to parse image reference: %v", err)
		}
		// should fail
		_, err = remote.Get(ref, remote.WithAuth(auth), remote.WithContext(ctx))
		if err != nil {
			fmt.Println("Error getting tag:", err)
			os.Exit(1)
		}
	*/

	// cosign part
	// Get the remote image reference.
	/*remoteRef, err := cosignremote.ResolveDigest(ref, cosignremote.WithRemoteOptions(cosignremote, auth))
	if err != nil {
		log.Fatalf("resolving image digest: %v", err)
	}*/

	// opts := options.RegistryOptions{
	// 	AuthConfig: authn.AuthConfig{
	// 		Username: username,
	// 		Password: password,
	// 	},
	// }

	// registryOpts, err := opts.ClientOpts(ctx)
	// if err != nil {
	// 	log.Fatalf("creating registry opts: %v", err)
	// }

	// co := &cosign.CheckOpts{
	// 	RegistryClientOpts: registryOpts,
	// 	IgnoreSCT:          true, // Ignore SCT checks for simplicity (adjust as needed)
	// 	IgnoreTlog:         true, //Ignore transparency log checks for simplicity(adjust as needed)
	// }

	// sigs, _, err := cosign.VerifyImageSignatures(ctx, ref, co)
	// if err != nil {
	// 	log.Fatalf("verifying signatures: %v", err)
	// }
	// fmt.Println(len(sigs))

	// add := utils.Add(1, 2)
	// fmt.Println(add)

	// utils.AddAndPrint(1, 2)

	// Accessing layers
	/*
		layers, err := img.Layers()
		if err != nil {
			log.Fatalf("getting layers: %v", err)
		}

		fmt.Println("\nLayers:")
		for _, layer := range layers {
			digest, err := layer.Digest()
			if err != nil {
				log.Printf("getting layer digest: %v", err)
				continue
			}
			mediaType, err := layer.MediaType()
			if err != nil {
				log.Printf("getting layer media type: %v", err)
				continue
			}

			fmt.Printf("  Digest: %v, MediaType: %v\n", digest, mediaType)
		}
	*/
}

func isSigned(tags []string, desc *remote.Descriptor) bool {
	for _, tag := range tags {
		if strings.Contains(tag, desc.Digest.Hex) {
			return true
		}
	}
	return false
}

func isSignedByGet(hash string, auth authn.Authenticator, config Config) bool {
	imageRef := config.RegistryURL + "/" + config.CompartmentOCID + "/" + config.ImageRepository + ":" + "sha256-" + hash + ".sig"
	ref, err := name.ParseReference(imageRef)
	if err != nil {
		log.Fatalf("Failed to parse image reference: %v", err)
	}
	_, err = remote.Get(ref, remote.WithAuth(auth))
	return err == nil
}

func isSignatureUsed(tagsWithHashes map[string]string, tag string) bool {
	signatures := make(map[string]bool)
	originals := make(map[string]bool)

	for tag, hash := range tagsWithHashes {
		if strings.HasPrefix(tag, "sha256-") && strings.HasSuffix(tag, ".sig") {
			signatures[tag] = true
		} else {
			originals[hash] = true
		}
	}
	base := strings.TrimSuffix(strings.TrimPrefix(tag, "sha256-"), ".sig")
	return originals[base]
}
