package setup

import "ksetup/pkg/log"

type Install struct {
	log *log.Logger
}

func (i *Install) Run() error {
	i.log.Info("Installing")
	return nil
}

func (i *Install) InstallEngine() error {
	return nil
}

func (i *Install) InstallCharts() error {
	return nil
}

func (i *Install) InstallManifest() error {
	return nil
}
