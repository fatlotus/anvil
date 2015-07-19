package hephaestus

import (
	"fmt"
)

// Prints a user-parseable description of the tree in sorted order.
func PrintTree(t Tree) {
	for b := range t {
		prefix := " "
		if b.Contents() == nil {
			prefix = "-"
		}
		fmt.Printf("%s %-40s %40s\n", prefix, b.Name(), b.Source())

		if b.Error() != nil {
			fmt.Printf("Error: %s\n", b.Error())
			return
		}
	}
}

func formatBlob(b Blob) string {
	result := "{"

	if b.Contents() == nil {
		result += "remove: "
	}

	result += b.Name()

	if b.Error() != nil {
		result += ", err"
	}

	if b.Source() != "" {
		result += ", source=" + b.Source()
	}

	if b.ModTime().Unix() > 0 {
		result += ", modified=" + b.ModTime().String()
	}

	if b.Size() > 0 {
		result += fmt.Sprintf(", size=%d", b.Size())
	}

	return result + "}"
}
