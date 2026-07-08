package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"new-repo/internal/config"
	"new-repo/internal/git"
	"new-repo/internal/scaffold"
	"new-repo/internal/skill"
)

var nameRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]*$`)

func main() {
	force := flag.Bool("force", false, "overwrite existing directory")
	skillsFile := flag.String("skills-file", "", "path to custom SKILLS.md")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "usage: new-repo [--force] [--skills-file <path>] <project-name | .>\n")
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
		if err := git.CloneTemplate(dir); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("installed ai-template at %s\n", dir)
		return
	}

	if !nameRe.MatchString(name) {
		fmt.Fprintf(os.Stderr, "error: invalid project name %q (must start with alphanumeric, containing alphanumeric, underscore, or hyphen)\n", name)
		os.Exit(1)
	}

	if !git.IsAvailable() {
		fmt.Fprintf(os.Stderr, "error: git is not installed\n")
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

	if err := scaffold.CreateOpenSpecDirs(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := scaffold.WriteOpenSpecConfig(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	if err := scaffold.WriteOpenSpecSkills(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var skillsContent string
	if *skillsFile != "" {
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

	skillNames, err := skill.ParseSkills(skillsContent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: parsing skills: %v\n", err)
		os.Exit(1)
	}

	if len(skillNames) > 0 {
		skillsDir := cfg.SkillsDir()
		skillsSrc, cleanup, err := git.CloneTemplateSkills()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: downloading skills: %v\n", err)
			os.Exit(1)
		}
		defer cleanup()

		entries, err := os.ReadDir(skillsSrc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: reading skills dir: %v\n", err)
			os.Exit(1)
		}
		available := make(map[string]bool, len(entries))
		for _, e := range entries {
			if e.IsDir() {
				available[e.Name()] = true
			}
		}
		toInstall := make([]string, 0, len(skillNames))
		for _, name := range skillNames {
			if available[name] {
				toInstall = append(toInstall, name)
			}
		}

		if err := skill.InstallSkills(skillsDir, []string{skillsSrc}, toInstall); err != nil {
			fmt.Fprintf(os.Stderr, "error: installing skills: %v\n", err)
			os.Exit(1)
		}
	}

	if err := git.Init(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("created project %q at %s\n", name, dir)
}
