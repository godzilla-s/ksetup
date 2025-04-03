package engine

import "ksetup/pkg/client/remote"

type Engine interface {
	Start() error
	Stop() error
	Restart() error
	Install() error
	Uninstall() error
	Status() (remote.Status, error)
}
