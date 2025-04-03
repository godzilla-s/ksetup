package setup

import "ksetup/pkg/log"

type Uninstall struct {
	log *log.Logger
}

func (u *Uninstall) Run() error {
	return nil
}
