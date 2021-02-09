package cmd

func imageRefMatchesImageName(imageRef string, imageName string) bool {
	// imageRef is a direct match to the imageName
	return imageRef == imageName ||
	// imageRef is a tag or a digest of the image name
	( len(imageRef) >= len(imageName)+1 && ( imageRef[:len(imageName)+1] == imageName+":" || imageRef[:len(imageName)+1] == imageName+"@" ) )
}

