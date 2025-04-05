package setup

import (
	"ksetup/pkg/config"
	"ksetup/pkg/log"
)

type Install struct {
	log *log.Logger
}

func NewInstall(conf *config.Config, l *log.Logger) (*Install, error) {
	return &Install{
		log: l,
	}, nil
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
