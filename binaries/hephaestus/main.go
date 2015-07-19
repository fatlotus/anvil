package main

import (
	"github.com/fatlotus/hephaestus"
)

func main() {
	image, err := hephaestus.LoadFromManifest("image/MANIFEST.json")
	if err != nil {
		panic(err)
	}
	
	state, err := hephaestus.LoadZipFile("Archive.zip")
	if err != nil {
		panic(err)
	}
	
	diff := hephaestus.Diff(state, image)
	hephaestus.PrintTree(diff)
}