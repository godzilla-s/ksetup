package config

import (
	"encoding/json"
	"os"
	"time"
)

type Chart struct {
	Name        string                 `json:"name"`
	ReleaseName string                 `json:"release_name"`
	Namespace   string                 `json:"namespace"`
	Path        string                 `json:"path"`
	Values      map[string]interface{} `json:"values"`
	Timeout     time.Duration          `json:"timeout"`
	Repo        *ChartRepo             `json:"repo"`
}

type ChartRepo struct {
	Url      string `json:"url"`
	Username string `json:"username"`
	Secret   string `json:"secret"`
}

type Node struct {
	Host         string   `json:"host"`
	Hostname     string   `json:"hostname"`
	Port         string   `json:"port"`
	RootPassword string   `json:"root_password"`
	PreInstall   []string `json:"pre_install"`
	PostInstall  []string `json:"post_install"`
}

type Cluster struct {
	Master    ClusterRole `json:"master"`
	Agent     ClusterRole `json:"agent"`
	Charts    []Chart     `json:"charts"`
	Manifests []string    `json:"manifests"`
}

type ClusterRole struct {
	Config  map[string]interface{} `json:"config"`
	Members []Node                 `json:"members"`
}

type GlobalSetting struct {
	VersiobRequired string     `json:"version_required"`
	Repo            *ChartRepo `json:"repo"`
}

type Config struct {
	GlobalSetting GlobalSetting    `json:"global"`
	Charts        map[string]Chart `json:"charts"`
	Install       Cluster          `json:"install"`
}

func Load(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
