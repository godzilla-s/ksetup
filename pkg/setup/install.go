package setup

import (
	"ksetup/pkg/config"
	"ksetup/pkg/log"
)

type Install struct {
	log          *log.Logger
	clusterNodes []*ClusterNode
	cluster      *cluster
}

func NewInstall(conf *config.Config, l *log.Logger) (*Install, error) {
	for _, member := range conf.Install.Master.Members {
		node, err := newClusterNode("k3s", &member, l)
		if err != nil {

		}
		node.SetMaster()
		node.SetInitial()
	}

	return &Install{
		log: l,
	}, nil
}

func (i *Install) Run() error {
	i.log.Info("Installing")
	return nil
}

func (i *Install) InstallEngine() error {
	// pre instakll
	for _, node := range i.clusterNodes {
		// pre install
		node.StartEngine()
	}
	return nil
}

func (i *Install) InstallCharts(charts []config.Chart) error {
	return nil
}

func (i *Install) InstallManifest() error {
	return nil
}

func (i *Install) applyChart(chart config.Chart) error {
	kubeClient, err := i.cluster.getKubeClient()
	if err != nil {
		return err
	}

	chartClient, err := kubeClient.NewChartClient()
	if err != nil {
		return err
	}

	_, err = chartClient.Get(chart.ReleaseName)
	if err != nil {
		return err
	}

	return nil
}
