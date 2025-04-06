package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Package struct {
	Name       string   `yaml:"name"`
	Type       string   `yaml:"type"`
	LocalPath  string   `yaml:"local_path"`
	TargetPath string   `yaml:"target_path"`
	PreExec    []string `yaml:"pre_exec"`
	Exec       []string `yaml:"exec"`
	PostExec   []string `yaml:"post_exec"`
}

type Chart struct {
	Name        string                 `yaml:"name"`
	ReleaseName string                 `yaml:"release_name"`
	Namespace   string                 `yaml:"namespace"`
	Path        string                 `yaml:"path"`
	Values      map[string]interface{} `yaml:"values"`
	Timeout     time.Duration          `yaml:"timeout"`
	Repo        *ChartRepo             `yaml:"repo"`
}

type ChartRepo struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Secret   string `yaml:"secret"`
}

type Node struct {
	Host         string   `yaml:"host"`
	Hostname     string   `yaml:"hostname"`
	Port         string   `yaml:"port"`
	RootPassword string   `yaml:"root_password"`
	PreInstall   []string `yaml:"pre_install"`
	PostInstall  []string `yaml:"post_install"`
}

type Cluster struct {
	Provider  string      `yaml:"provider"`
	Master    ClusterRole `yaml:"master"`
	Agent     ClusterRole `yaml:"agent"`
	Charts    []string    `yaml:"charts"`
	Manifests []string    `yaml:"manifests"`
}

type ClusterRole struct {
	Config  map[string]interface{} `yaml:"config"`
	Members []Node                 `yaml:"members"`
}

type GlobalSetting struct {
	VersiobRequired string     `yaml:"version_required"`
	Repo            *ChartRepo `yaml:"repo"`
}

type Config struct {
	GlobalSetting GlobalSetting    `yaml:"global"`
	Charts        map[string]Chart `yaml:"charts"`
	Install       Cluster          `yaml:"install"`
}

func Load(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
