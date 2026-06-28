package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	ProjectName string
	ProjectDir  string
	Force       bool
	SkillsSrcs  []string
	SkillsFile  string
	GitBin      string
}

func New(name string, force bool) Config {
	return Config{
		ProjectName: name,
		ProjectDir:  name,
		Force:       force,
		SkillsSrcs:  skillsSourceDirs(),
		GitBin:      "git",
	}
}

func skillsSourceDirs() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{".config/opencode/skills", ".opencode/skills"}
	}
	return []string{
		filepath.Join(home, ".config", "opencode", "skills"),
		filepath.Join(home, ".opencode", "skills"),
	}
}

func (c Config) SkillsDir() string {
	return filepath.Join(c.ProjectDir, ".opencode", "skills")
}
