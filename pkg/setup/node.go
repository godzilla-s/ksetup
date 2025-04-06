package setup

import (
	"fmt"
	"ksetup/pkg/client/remote"
	"ksetup/pkg/config"
	"ksetup/pkg/log"
	"ksetup/pkg/setup/engine"
	"ksetup/pkg/setup/engine/k3s"
)

type ClusterNode struct {
	isMaster      bool
	isInitialized bool
	client        *remote.Client
	engine        engine.Engine
	preInstalls   []installation
	postInstalls  []installation
	log           *log.Entry
}

func newClusterNode(provider string, conf *config.Node, log *log.Logger) (*ClusterNode, error) {
	logEntry := log.NewEntry(map[string]interface{}{
		"node": conf.Hostname,
		"host": conf.Host,
	})
	client, err := remote.New(remote.Config{}, logEntry)
	if err != nil {
		return nil, err
	}

	var e engine.Engine
	switch provider {
	case "k3s":
		e, err = k3s.New()
		if err != nil {
			logEntry.Errorf("failed to get k3s engine: %v", err)
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	cn := &ClusterNode{
		client: client,
		engine: e,
	}
	return cn, nil
}

func (c *ClusterNode) SetInitial() {
	c.isInitialized = true
}

func (c *ClusterNode) SetMaster() {
	c.isMaster = true
}

func (c *ClusterNode) Status() error {
	return nil
}

func (c *ClusterNode) StartEngine() error {
	return c.engine.Start()
}

func (c *ClusterNode) GetKubeConfig() ([]byte, error) {
	return c.client.ReadFile("/etc/rancher/k3s/k3s.yaml")
}
