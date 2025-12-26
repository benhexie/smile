package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	// Define source and build directories
	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}
	srcDir := filepath.Join(rootDir, "src")
	buildDir := filepath.Join(rootDir, "build")

	// Ensure build directory exists
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		if err := os.Mkdir(buildDir, 0755); err != nil {
			fmt.Printf("Error creating build directory: %v\n", err)
			os.Exit(1)
		}
	}

	// Change to src directory for go commands
	if err := os.Chdir(srcDir); err != nil {
		fmt.Printf("Error changing to src directory: %v\n", err)
		os.Exit(1)
	}

	// Install and run fileb0x to generate static assets
	fmt.Println("Installing fileb0x (static file generator)...")
	runCommand("go", "install", "github.com/UnnoTed/fileb0x@latest")

	fmt.Println("Generating static files...")
	// Try running fileb0x from PATH
	if err := exec.Command("fileb0x", "b0x.yaml").Run(); err != nil {
		// If not in PATH, try common Go bin location
		homeDir, _ := os.UserHomeDir()
		b0xPath := filepath.Join(homeDir, "go", "bin", "fileb0x")
		if runtime.GOOS == "windows" {
			b0xPath += ".exe"
		}
		
		fmt.Printf("fileb0x not found in PATH, trying %s...\n", b0xPath)
		if err := exec.Command(b0xPath, "b0x.yaml").Run(); err != nil {
			fmt.Printf("Warning: Failed to run fileb0x: %v. Build might fail if static files are missing.\n", err)
		}
	}

	fmt.Println("Installing dependencies...")
	runCommand("go", "mod", "tidy")

	fmt.Println("Building...")
	
	// Prepare build command
	outputName := "smile.exe"
	outputPath := filepath.Join(buildDir, outputName)
	
	// -s -w flags to strip debug information and reduce binary size (optional but good for 'release')
	// -H=windowsgui to hide console window
	args := []string{"build", "-ldflags", "-H=windowsgui -s -w", "-o", outputPath}

	cmd := exec.Command("go", args...)
	
	// Set environment variables for cross-compilation
	env := os.Environ()
	
	// Always target Windows since this tool is Windows-specific
	// But check if we are already on Windows to avoid redundancy (though setting it explicitly is safe)
	if runtime.GOOS != "windows" {
		fmt.Println("Detected non-Windows OS. Cross-compiling for Windows...")
		env = append(env, "GOOS=windows", "GOARCH=amd64")
		// Enable CGO_ENABLED=0 for easier cross-compilation usually, 
		// but this project uses 'modernc.org/sqlite' which is CGO-free (pure Go),
		// so CGO_ENABLED=0 is likely fine. 
		// However, check if any other dependencies need CGO.
		// github.com/billgraziano/dpapi calls Windows DLLs.
		// Cross-compiling Go with Windows syscalls usually works fine without CGO if no C code is involved.
		// modernc.org/sqlite is a CGO-free port of SQLite.
		// So we should be good.
		env = append(env, "CGO_ENABLED=0")
	}
	cmd.Env = env

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Build failed: %v\n", err)
		fmt.Println("Ensure you have Go installed and if on macOS/Linux, that you can cross-compile.")
		os.Exit(1)
	}

	// Copy config.prop to build directory
	srcConfig := filepath.Join(rootDir, "config.prop")
	dstConfig := filepath.Join(buildDir, "config.prop")
	copyFile(srcConfig, dstConfig)

	fmt.Println("Build success :)")
	fmt.Printf("smile.exe is in %s\n", buildDir)
}

func copyFile(src, dst string) {
	input, err := os.ReadFile(src)
	if err != nil {
		fmt.Printf("Warning: Could not copy config.prop: %v\n", err)
		return
	}
	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		fmt.Printf("Warning: Could not write config.prop to build dir: %v\n", err)
	}
}

func runCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Command %s failed: %v\n", name, err)
		os.Exit(1)
	}
}