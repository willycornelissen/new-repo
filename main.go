package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const repoURL = "https://github.com/willycornelissen/ai-template"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: new-repo <directory-name>")
		os.Exit(1)
	}

	dirName := os.Args[1]

	if dirName == "" {
		fmt.Fprintln(os.Stderr, "Error: directory name cannot be empty")
		os.Exit(1)
	}

	absPath, err := filepath.Abs(dirName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}

	if _, err := os.Stat(absPath); err == nil {
		fmt.Fprintf(os.Stderr, "Error: directory %q already exists\n", dirName)
		os.Exit(1)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created directory %q\n", dirName)

	fmt.Println("Downloading template...")
	clone := exec.Command("git", "clone", "--depth", "1", repoURL, absPath)
	clone.Stdout = os.Stdout
	clone.Stderr = os.Stderr
	if err := clone.Run(); err != nil {
		os.RemoveAll(absPath)
		fmt.Fprintf(os.Stderr, "Error downloading template: %v\n", err)
		os.Exit(1)
	}

	gitDir := filepath.Join(absPath, ".git")
	if err := os.RemoveAll(gitDir); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing template .git: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Initializing new git repository...")
	init := exec.Command("git", "init")
	init.Dir = absPath
	init.Stdout = os.Stdout
	init.Stderr = os.Stderr
	if err := init.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing git repo: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nDone! New project created at %q\n", absPath)
}
