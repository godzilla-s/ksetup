package setup

import (
	"ksetup/pkg/client/remote"
	"ksetup/pkg/log"
)

type installation interface {
	install() error
	uninstall() error
}

type file struct {
	localPath  string
	targetPath string
	override   bool
	preExec    []string
	exec       []string
	postExec   []string
	log        *log.Logger
	client     *remote.Client
}

func (f *file) install() error {
	if len(f.preExec) > 0 {
		for _, exec := range f.preExec {
			_, err := f.client.RunCommand([]string{exec})
			if err != nil {
				return err
			}
		}
	}
	f.client.Copy(f.localPath, f.targetPath, f.override)
	return nil
}

func (f *file) uninstall() error {
	return nil
}
