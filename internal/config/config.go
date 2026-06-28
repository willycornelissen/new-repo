package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	ProjectName string
	ProjectDir  string
	Force       bool
	SkillsSrc   string
	SkillsFile  string
	GitBin      string
}

func New(name string, force bool) Config {
	return Config{
		ProjectName: name,
		ProjectDir:  name,
		Force:       force,
		SkillsSrc:   skillsSourceDir(),
		GitBin:      "git",
	}
}

func skillsSourceDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "opencode", "skills")
}

func (c Config) SkillsDir() string {
	return filepath.Join(c.ProjectDir, ".opencode", "skills")
}
