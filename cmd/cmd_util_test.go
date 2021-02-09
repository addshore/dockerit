package cmd

import "testing"

func Test_imageRefMatchesImageName(t *testing.T) {

	type test struct {
		match bool
		imageRef string
		imageName string
	}

	tests := []test{
		{match: true, imageRef: "node", imageName: "node"},
		{match: true, imageRef: "node:latest", imageName: "node"},
		{match: true, imageRef: "node@aDigest", imageName: "node"},
		{match: true, imageRef: "composer", imageName: "composer"},
		{match: true, imageRef: "composer:1", imageName: "composer"},
		{match: true, imageRef: "composer@foobar", imageName: "composer"},
		{match: true, imageRef: "composer:some-tag", imageName: "composer"},
		{match: false, imageRef: "compose", imageName: "composer"},
		{match: false, imageRef: "composerr", imageName: "composer"},
	}

	for _, tc := range tests {
		result := imageRefMatchesImageName(tc.imageRef, tc.imageName)
		if result != tc.match {
			t.Errorf("Ref match condition failure: ref: %s, image: %s.", tc.imageRef, tc.imageName)
		}
	}

}
