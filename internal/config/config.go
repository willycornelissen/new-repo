package config

import (
	"path/filepath"
)

type Config struct {
	ProjectName string
	ProjectDir  string
	Force       bool
	SkillsFile  string
	GitBin      string
}

func New(name string, force bool) Config {
	return Config{
		ProjectName: name,
		ProjectDir:  name,
		Force:       force,
		GitBin:      "git",
	}
}

func (c Config) SkillsDir() string {
	return filepath.Join(c.ProjectDir, ".opencode", "skills")
}
