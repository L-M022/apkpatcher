package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"command failed: %s %s\n%s\nerror: %w",
			name,
			strings.Join(args, " "),
			string(output),
			err,
		)
	}

	return nil
}

func javaCommand(args ...string) error {
	return runCommand("java", args...)
}

func DecompileAPK(apk string, outputDir string) error {

	if _, err := os.Stat(apk); err != nil {
		return fmt.Errorf("APK not found: %s", apk)
	}

	if err := os.RemoveAll(outputDir); err != nil {
		return fmt.Errorf("cannot clean output directory: %w", err)
	}

	err := javaCommand(
		"-jar",
		ApktoolPath(),
		"d",
		"-f",
		"--keep-broken-res",
		apk,
		"-o",
		outputDir,
	)

	if err != nil {
		return fmt.Errorf("apktool decode failed: %w", err)
	}

	return nil
}
func ApktoolPath() string {
	base := GetBaseDir()
	return filepath.Join(base, "tools", "apktool_3.0.2.jar")
}
func GetBaseDir() string {

	exePath, err := os.Executable()
	if err != nil {
		return "."
	}

	exeDir := filepath.Dir(exePath)

	// Detecta go run (ejecutable temporal)
	if strings.Contains(exeDir, "go-build") {

		current, err := os.Getwd()

		if err == nil {
			return current
		}
	}

	return exeDir
}
func CheckAPKTools() error {

	required := []string{
		ToolPath("apktool_3.0.2.jar"),
		ToolPath("apksigner.jar"),
		ToolPath("zipalign.exe"),
		ToolPath("ApkToolkit_Key.pk8"),
		ToolPath("ApkToolkit_Certificate.pem"),
	}

	for _, file := range required {
		if _, err := os.Stat(file); err != nil {
			return fmt.Errorf("missing required tool: %s", file)
		}
	}

	return nil
}

func BuildAPK(projectDir string, outputAPK string) error {

	if _, err := os.Stat(projectDir); err != nil {
		return fmt.Errorf("apktool project not found: %s", projectDir)
	}

	err := javaCommand(
		"-jar",
		ApktoolPath(),
		"b",
		projectDir,
		"-o",
		outputAPK,
	)

	if err != nil {
		return fmt.Errorf("apktool build failed: %w", err)
	}

	return nil
}

func ZipalignAPK(inputAPK string, outputAPK string) error {

	if _, err := os.Stat(inputAPK); err != nil {
		return fmt.Errorf("APK to align not found: %s", inputAPK)
	}

	err := runCommand(
		ToolPath("zipalign.exe"),
		"-f",
		"4",
		inputAPK,
		outputAPK,
	)

	if err != nil {
		return fmt.Errorf("zipalign failed: %w", err)
	}

	return nil
}
func ToolPath(name string) string {

	return filepath.Join(
		GetBaseDir(),
		"tools",
		name,
	)
}
func CleanApktoolFramework() error {

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("cannot get user home directory: %w", err)
	}

	frameworkDir := filepath.Join(
		home,
		"AppData",
		"Local",
		"apktool",
		"framework",
	)

	if _, err := os.Stat(frameworkDir); os.IsNotExist(err) {
		return nil
	}

	err = os.RemoveAll(frameworkDir)
	if err != nil {
		return fmt.Errorf("failed removing apktool framework directory: %w", err)
	}

	return nil
}
func SignAPK(apk string) error {

	if _, err := os.Stat(apk); err != nil {
		return fmt.Errorf("APK to sign not found: %s", apk)
	}

	err := javaCommand(
		"-jar",
		ToolPath("apksigner.jar"),
		"sign",
		"--key",
		ToolPath("ApkToolkit_Key.pk8"),
		"--cert",
		ToolPath("ApkToolkit_Certificate.pem"),
		apk,
	)

	if err != nil {
		return fmt.Errorf("APK signing failed: %w", err)
	}

	return nil
}

func ReplaceAPK(original string, modified string) error {

	if _, err := os.Stat(modified); err != nil {
		return fmt.Errorf("modified APK does not exist: %s", modified)
	}

	backup := original + ".backup"

	if err := os.Rename(original, backup); err != nil {
		return fmt.Errorf("cannot create APK backup: %w", err)
	}

	if err := os.Rename(modified, original); err != nil {

		// restaurar si falla
		os.Rename(backup, original)

		return fmt.Errorf("cannot replace APK: %w", err)
	}

	os.Remove(backup)

	return nil
}

func ProcessAPKTools(projectDir string, outputDir string, finalAPK string) error {

	unsignedAPK := filepath.Join(
		outputDir,
		"unsigned.apk",
	)

	alignedAPK := filepath.Join(
		outputDir,
		"aligned.apk",
	)

	if err := BuildAPK(projectDir, unsignedAPK); err != nil {
		return err
	}

	if err := ZipalignAPK(unsignedAPK, alignedAPK); err != nil {
		return err
	}

	if err := SignAPK(alignedAPK); err != nil {
		return err
	}

	if err := os.Rename(alignedAPK, finalAPK); err != nil {
		return fmt.Errorf("cannot move final APK: %w", err)
	}

	return nil
}
