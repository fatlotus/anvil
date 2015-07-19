package main

import (
	"github.com/fatlotus/hephaestus"
)

// func loadSpec(path string) hephaestus.Tree {
// 	parts := strings.Split(path, ":")
// 	accum := loadTree(parts[0])
//
// 	for i := 1; i < len(parts); i++ {
// 		accum = hephaestus.Validate(hephaestus.Overlay(accum,
// 				loadTree(parts[i])))
// 	}
//
// 	return accum
// }
//
// func loadTree(path string) hephaestus.Tree {
// 	rdr, err := zip.OpenReader(path)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	prefix := filepath.Dir(path)
//
// 	return hephaestus.Validate(
// 		hephaestus.WithPrefix(
// 			hephaestus.Validate(hephaestus.TreeFromZip(rdr, path)), prefix))
// }

func main() {
	
	command := flag
	
	// fp, _ := os.Create("FS/MANIFEST.json")
	// hephaestus.SaveToManifest("FS/", fp)
	// fp.Close()
	
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
	
	// flag.Parse()
	//
	// a := loadSpec(flag.Arg(0))
	// b := loadSpec(flag.Arg(1))
	//
	// diff := hephaestus.Diff(a, b)
	//
	// hephaestus.PrintTree(diff)

	// fp, err := os.Create("FS2/root.zip")
	// if err != nil {
	// 	panic(err)
	// }
	//
	// w := zip.NewWriter(fp)
	//
	// if err := t.ToZip(w); err != nil {
	// 	panic(err)
	// }
	//
	// w.Close()
	
	// for blob := range d {
	// 	fmt.Printf(" : %s\n", hephaestus.FormatBlob(blob))
	//
	// 	if blob.Error() != nil {
	// 		fmt.Printf("error: %s\n", blob.Error().Error())
	// 		return
	// 	}
	// }
}