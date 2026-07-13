package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func CopyBanner(apktoolDir string, bannerPath string) error {

	fmt.Println("BannerPath:", bannerPath)

	abs, _ := filepath.Abs(bannerPath)
	fmt.Println("Absolute:", abs)

	if _, err := os.Stat(bannerPath); err != nil {
		return fmt.Errorf("banner file not found: %s (%v)", bannerPath, err)
	}
	// if _, err := os.Stat(bannerPath); err != nil {
	// 	return fmt.Errorf(
	// 		"banner file not found: %s",
	// 		bannerPath,
	// 	)
	// }

	drawableDir := filepath.Join(
		apktoolDir,
		"res",
		"drawable-nodpi",
	)

	if err := os.MkdirAll(drawableDir, 0755); err != nil {
		return fmt.Errorf(
			"cannot create drawable directory: %w",
			err,
		)
	}

	ext := filepath.Ext(bannerPath)

	destination := filepath.Join(
		drawableDir,
		"tv_banner"+ext,
	)

	data, err := os.ReadFile(bannerPath)
	if err != nil {
		return fmt.Errorf(
			"cannot read banner: %w",
			err,
		)
	}

	if err := os.WriteFile(destination, data, 0644); err != nil {
		return fmt.Errorf(
			"cannot write banner: %w",
			err,
		)
	}

	return nil
}
