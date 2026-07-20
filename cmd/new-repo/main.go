package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"new-repo/internal/config"
	"new-repo/internal/embed"
	"new-repo/internal/git"
	"new-repo/internal/scaffold"
	"new-repo/internal/skill"
)

var nameRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

func main() {
	force := flag.Bool("force", false, "overwrite existing directory")
	skillsFile := flag.String("skills-file", "", "path to custom SKILLS.md")
	gstack := flag.Bool("gstack", false, "install skills from garrytan/gstack instead of AI Template")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: new-repo [--force] [--skills-file <path>] [--gstack] <project-name | .>\n")
		os.Exit(1)
	}

	name := flag.Arg(0)

	if name == "." {
		if !git.IsAvailable() {
			fmt.Fprintf(os.Stderr, "error: git is not installed\n")
			os.Exit(1)
		}
		dir, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := embed.ExtractTemplate(dir); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}

		if *gstack {
			skillsDir := filepath.Join(dir, ".opencode", "skills")
			os.RemoveAll(skillsDir)
			if err := os.MkdirAll(skillsDir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}

			// Remove commands directory if --gstack is true
			os.RemoveAll(filepath.Join(dir, ".opencode", "commands"))

			skillsContent, _, err := installGStackSkills(skillsDir, dir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}

			if err := scaffold.WriteSkillsMD(dir, skillsContent); err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				os.Exit(1)
			}
		}

		if err := git.Init(dir); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if *gstack {
			fmt.Printf("installed gstack ai-template at %s\n", dir)
		} else {
			fmt.Printf("installed ai-template at %s\n", dir)
		}
		return
	}

	if !nameRe.MatchString(name) {
		fmt.Fprintf(os.Stderr, "error: invalid project name %q (must start with alphanumeric, containing alphanumeric, underscore, or hyphen)\n", name)
		os.Exit(1)
	}

	cfg := config.New(name, *force)

	dir, err := scaffold.CreateProjectDir(name, *force)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := scaffold.WriteGitignore(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := scaffold.CreateOpenCodeDirs(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var skillsContent string
	var skillNames []string
	var gstackInstalled bool

	if *gstack {
		skillsDir := cfg.SkillsDir()
		var err error
		skillsContent, skillNames, err = installGStackSkills(skillsDir, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		gstackInstalled = true
	} else if *skillsFile != "" {
		data, err := os.ReadFile(*skillsFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: reading skills file: %v\n", err)
			os.Exit(1)
		}
		skillsContent = string(data)
	} else {
		var err error
		skillsContent, err = skill.ReadEmbedded()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	if err := scaffold.WriteSkillsMD(name, skillsContent); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := scaffold.WriteAgentsMD(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if !gstackInstalled {
		var err error
		skillNames, err = skill.ParseSkills(skillsContent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: parsing skills: %v\n", err)
			os.Exit(1)
		}

		if len(skillNames) > 0 {
			skillsDir := cfg.SkillsDir()
			toInstall := embed.ListAvailableSkillNames(skillNames)
			if err := embed.ExtractSkills(skillsDir, toInstall); err != nil {
				fmt.Fprintf(os.Stderr, "error: extracting skills: %v\n", err)
				os.Exit(1)
			}
		}
	}

	if !*gstack {
		commandsDir := filepath.Join(name, ".opencode", "commands")
		if err := embed.ExtractOpenCodeCommands(commandsDir); err != nil {
			fmt.Fprintf(os.Stderr, "error: extracting commands: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Ensure no empty commands folder exists
		os.RemoveAll(filepath.Join(name, ".opencode", "commands"))
	}

	docsDir := filepath.Join(name, ".opencode", "docs")
	if err := embed.ExtractOpenCodeDocs(docsDir); err != nil {
		fmt.Fprintf(os.Stderr, "error: extracting docs: %v\n", err)
		os.Exit(1)
	}

	if err := git.Init(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("created project %q at %s\n", name, dir)
}

func installGStackSkills(skillsDir string, targetDir string) (string, []string, error) {
	gstackTempDir, err := os.MkdirTemp("", "gstack-*")
	if err != nil {
		return "", nil, fmt.Errorf("creating temp dir: %w", err)
	}
	defer os.RemoveAll(gstackTempDir)

	fmt.Printf("cloning garrytan/gstack...\n")
	cmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/garrytan/gstack.git", gstackTempDir)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("git clone failed: %w\n%s", err, out)
	}

	skills, err := skill.FindGStackSkills(gstackTempDir)
	if err != nil {
		return "", nil, fmt.Errorf("finding gstack skills: %w", err)
	}

	skillsContent := skill.GenerateSkillsMD(skills)

	var skillNames []string
	for _, s := range skills {
		skillNames = append(skillNames, s.Name)
	}

	srcDirsMap := make(map[string]bool)
	for _, s := range skills {
		srcDirsMap[filepath.Dir(s.Path)] = true
	}
	var srcDirs []string
	for d := range srcDirsMap {
		srcDirs = append(srcDirs, d)
	}

	if err := skill.InstallSkills(skillsDir, srcDirs, skillNames); err != nil {
		return "", nil, fmt.Errorf("installing gstack skills: %w", err)
	}

	return skillsContent, skillNames, nil
}
