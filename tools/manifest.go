package tools

import (
	"os"
	"path/filepath"
	"strings"
)

type ManifestOptions struct {
	AddLeanback bool

	AddBanner  bool
	BannerName string
	BannerPath string
}

func addBanner(xml string, banner string) string {

	if strings.Contains(xml, "android:banner=") {
		return xml
	}

	return strings.Replace(
		xml,
		"<application ",
		`<application android:banner="@drawable/`+banner+`" `,
		1,
	)
}
func addLeanbackFeature(xml string) string {

	if strings.Contains(xml, "android.software.leanback") {
		return xml
	}

	feature :=
		"\n    <uses-feature android:name=\"android.software.leanback\" android:required=\"false\"/>\n"

	return strings.Replace(
		xml,
		"<application",
		feature+"<application",
		1,
	)
}
func addLeanbackLauncher(xml string) string {

	if strings.Contains(xml, "android.intent.category.LEANBACK_LAUNCHER") {
		return xml
	}

	return strings.Replace(
		xml,
		`<category android:name="android.intent.category.LAUNCHER"/>`,
		`<category android:name="android.intent.category.LAUNCHER"/>
                <category android:name="android.intent.category.LEANBACK_LAUNCHER"/>`,
		1,
	)
}
func ModifyManifestFile(manifestPath string, opt ManifestOptions) error {

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return err
	}

	manifest := string(data)

	if opt.AddBanner && opt.BannerName != "" {
		manifest = addBanner(manifest, opt.BannerName)
	}

	if opt.AddLeanback {
		manifest = addLeanbackFeature(manifest)
		manifest = addLeanbackLauncher(manifest)
	}

	return os.WriteFile(
		manifestPath,
		[]byte(manifest),
		0644,
	)
}

func ModifyManifest(apk string, opt ManifestOptions) error {

	if err := CheckAPKTools(); err != nil {
		return err
	}
	if err := CleanApktoolFramework(); err != nil {
		return err
	}
	tempDir := apk + "-manifest-temp"

	defer os.RemoveAll(tempDir)

	err := DecompileAPK(
		apk,
		tempDir,
	)
	if err != nil {
		return err
	}

	err = FixMorpheAttrs(tempDir)

	if err != nil {
		return err
	}

	manifestPath := filepath.Join(
		tempDir,
		"AndroidManifest.xml",
	)

	if opt.AddBanner &&
		opt.BannerPath != "" &&
		opt.BannerName != "" {

		err = CopyBanner(
			tempDir,
			opt.BannerPath,
		)
		if err != nil {
			return err
		}
	}

	err = ModifyManifestFile(
		manifestPath,
		opt,
	)

	if err != nil {
		return err
	}

	outputDir := apk + "-build"

	defer os.RemoveAll(outputDir)

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	unsignedAPK := filepath.Join(
		outputDir,
		"unsigned.apk",
	)

	alignedAPK := filepath.Join(
		outputDir,
		"aligned.apk",
	)

	err = BuildAPK(
		tempDir,
		unsignedAPK,
	)

	if err != nil {
		return err
	}

	err = ZipalignAPK(
		unsignedAPK,
		alignedAPK,
	)

	if err != nil {
		return err
	}

	err = SignAPK(
		alignedAPK,
	)

	if err != nil {
		return err
	}

	backup := apk + ".old"

	err = os.Rename(
		apk,
		backup,
	)

	if err != nil {
		return err
	}

	err = os.Rename(
		alignedAPK,
		apk,
	)

	if err != nil {

		os.Rename(
			backup,
			apk,
		)

		return err
	}

	os.Remove(backup)

	return nil
}
