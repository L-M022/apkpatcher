package tools

import (
	"os"
	"path/filepath"
	"strings"
)

func FixMorpheAttrs(tempDir string) error {

	attrsPath := filepath.Join(
		tempDir,
		"res",
		"values",
		"attrs.xml",
	)

	data, err := os.ReadFile(attrsPath)
	if err != nil {
		return err
	}

	content := string(data)

	replacements := map[string]string{
		`<item type="attr" name="morphe_header" format="string" />`: `<attr name="morphe_header" format="reference" />`,

		`<item type="attr" name="morphe_header_custom" format="string" />`: `<attr name="morphe_header_custom" format="reference" />`,
	}

	for old, new := range replacements {
		content = strings.ReplaceAll(content, old, new)
	}

	return os.WriteFile(
		attrsPath,
		[]byte(content),
		0644,
	)
}
