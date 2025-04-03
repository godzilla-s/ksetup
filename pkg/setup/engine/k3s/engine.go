package k3s

import (
	"ksetup/pkg/client/remote"
	"ksetup/pkg/log"
	"ksetup/pkg/setup/engine"
)

type k3sEngine struct {
	isMaster bool
	log      *log.Logger
	client   *remote.Client
}

func New() (engine.Engine, error) {
	return nil, nil
}

func (e *k3sEngine) Install() error {
	cmd := []string{"INSTALL_K3S_SKIP_DOWNLOAD=true", "install.sh"}
	if e.isMaster {
		cmd = []string{"INSTALL_K3S_SKIP_DOWNLOAD=true", "INSTALL_K3S_EXEC='agent'", "install.sh"}
	}
	output, err := e.client.RunCommand(cmd)
	if err != nil {
		return err
	}
	_ = output
	return nil
}

func (e *k3sEngine) Uninstall() error {
	cmd := []string{"uninstall.sh"}
	if !e.isMaster {
		cmd = []string{"uninstall-agent.sh"}
	}
	output, err := e.client.RunCommand(cmd)
	if err != nil {
		return err
	}
	_ = output
	return nil
}

func (e *k3sEngine) Start() error {
	cmd := []string{"systemctl", "start"}
	if e.isMaster {
		cmd = append(cmd, "k3s")
	} else {
		cmd = append(cmd, "k3s-agent")
	}
	output, err := e.client.RunCommand(cmd)
	if err != nil {
		e.log.Errorf("failed to start k3s, output: %s, error: %v", output, err)
		return err
	}
	return nil
}

func (e *k3sEngine) Stop() error {
	cmd := []string{"systemctl", "stop"}
	if e.isMaster {
		cmd = append(cmd, "k3s")
	} else {
		cmd = append(cmd, "k3s-agent")
	}
	output, err := e.client.RunCommand(cmd)
	if err != nil {
		e.log.Errorf("failed to stop k3s, output: %s, error: %v", output, err)
		return err
	}
	return nil
}

func (e *k3sEngine) Restart() error {
	cmd := []string{"systemctl", "restart"}
	if e.isMaster {
		cmd = append(cmd, "k3s")
	} else {
		cmd = append(cmd, "k3s-agent")
	}
	output, err := e.client.RunCommand(cmd)
	if err != nil {
		e.log.Errorf("failed to restart k3s, output: %s, error: %v", output, err)
		return err
	}
	return nil
}

func (e *k3sEngine) Status() (remote.Status, error) {
	svcName := "k3s"
	if e.isMaster {
		svcName = "k3s-agent"
	}
	return e.client.SystemdStatus(svcName)
}
