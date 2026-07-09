package main

import (
	"flag"
	"fmt"
	"os"
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
		if err := embed.ExtractTemplate(dir); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if err := git.Init(dir); err != nil {
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
		toInstall := embed.ListAvailableSkillNames(skillNames)
		if err := embed.ExtractSkills(skillsDir, toInstall); err != nil {
			fmt.Fprintf(os.Stderr, "error: extracting skills: %v\n", err)
			os.Exit(1)
		}
	}

	if err := git.Init(name); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("created project %q at %s\n", name, dir)
}
