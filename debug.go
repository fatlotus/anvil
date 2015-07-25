package anvil

import (
	"fmt"
)

// Prints a user-parseable description of the tree in sorted order.
func (t Tree) ToDebugOutput() {
	prev := make([]rune, 0)

	for b := range t {
		prefix := " "
		if b.Contents() == nil {
			prefix = "-"
		}
		idx := 0
		curr := []rune(b.Name())

		for i := range prev {
			if prev[i] != curr[i] {
				break
			} else if prev[i] == '/' {
				idx = i + 1
				prefix += "  "
			}
		}

		basename := string(b.Name()[idx:])
		prev = curr

		fmt.Printf("%s %s (\033[33m%s\033[0m)\n", prefix, basename, b.Source())

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

	result += fmt.Sprintf(", mode=%s", b.Mode().String())

	return result + "}"
}
