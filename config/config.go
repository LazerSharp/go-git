package config

import (
	"fmt"
	"io"
)

type GitConfig map[string]map[string]string

func NewGitConfig() GitConfig {
	return GitConfig{}
}

func DeaultConfig() GitConfig {
	dc := NewGitConfig()
	dc.AddSection("core")
	dc.Add("core", "repositoryformatversion", "0")
	dc.Add("core", "filemode", "true")
	dc.Add("core", "bare", "false")
	dc.Add("core", "logallrefupdates", "true")
	return dc
}

func (c GitConfig) Write(w io.Writer) {
	for secname, section := range c {
		fmt.Fprintf(w, "[%s]\n", secname)
		for k, v := range section {
			fmt.Fprintf(w, "\t%s = %s\n", k, v)
		}
	}
}

func (c GitConfig) AddSection(section string) error {
	if c[section] != nil {
		return fmt.Errorf("section '%s' already present", section)
	}
	c[section] = map[string]string{}
	return nil
}

func (c GitConfig) Add(section, key, val string) error {
	if c[section] == nil {
		return fmt.Errorf("section '%s' not present", section)
	}
	c[section][key] = val
	return nil
}
