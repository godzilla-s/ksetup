package setup

import (
	"fmt"
	"ksetup/pkg/client/kube"
	"ksetup/pkg/log"
)

type cluster struct {
	clusterIP string
	members   []*ClusterNode
	log       *log.Logger
}

func (c *cluster) addMember(node *ClusterNode) {
	c.members = append(c.members, node)
}

func (c *cluster) getKubeClient() (*kube.Client, error) {
	for _, m := range c.members {
		if m.isMaster {
			kubeConfig, err := m.GetKubeConfig()
			if err != nil {
				c.log.Errorf("failed to get kubecofnig: %v", err)
				return nil, err
			}

			return kube.New(c.clusterIP, kubeConfig, c.log)
		}
	}
	return nil, fmt.Errorf("no master found in cluster")
}
