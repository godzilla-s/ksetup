package remote

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Client struct {
	log  *logrus.Entry
	ssh  *ssh.Client
	sftp *sftp.Client
}

func New(c Config, log *logrus.Entry) (*Client, error) {
	cfg := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshCli, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", c.Host, c.Port), cfg)
	if err != nil {
		return nil, err
	}

	sftpCli, err := sftp.NewClient(sshCli)
	if err != nil {
		return nil, err
	}
	return &Client{
		ssh:  sshCli,
		sftp: sftpCli,
		log:  log,
	}, nil
}

func (c *Client) RunCommand(cmds []string) ([]byte, error) {
	session, err := c.ssh.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()
	return session.CombinedOutput(strings.Join(cmds, " "))
}

func (c *Client) Close() {
	if c.sftp != nil {
		c.sftp.Close()
	}
	if c.ssh != nil {
		c.ssh.Close()
	}
}

func (c *Client) Copy(src, dst string, override bool) error {
	stat, err := os.Stat(src)
	if err != nil {
		c.log.Errorf("failed to stat source file %s: %v", src, err)
		return err
	}

	dstStat, err := c.sftp.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			c.log.Errorf("failed to stat target file %s: %v", dst, err)
			return err
		}

		err = c.sftp.MkdirAll(filepath.Dir(dst))
		if err != nil {
			c.log.Errorf("failed to create target directory %s: %v", dst, err)
			return err
		}
	}

	if dstStat != nil && dstStat.IsDir() {
		dst = filepath.Join(dst, filepath.Base(src))
	}

	reader, err := os.Open(src)
	if err != nil {
		c.log.Errorf("failed to open source file %s: %v", src, err)
		return err
	}

	writer, err := c.sftp.Create(dst)
	if err != nil {
		c.log.Errorf("failed to create target file %s: %v", dst, err)
		reader.Close()
		return err
	}

	defer func() {
		reader.Close()
		writer.Close()
	}()

	_, err = io.Copy(writer, reader)
	if err != nil {
		c.log.Errorf("failed to copy source file %s to dst %s: %v", src, dst, err)
		return err
	}

	c.sftp.Chmod(dst, stat.Mode())

	return nil
}

func (c *Client) RemoveFile(file string, force bool) error {
	stat, err := c.sftp.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if stat.IsDir() {
		if !force {
			return fmt.Errorf("cannot remove directory %s", file)
		}
		return c.sftp.RemoveDirectory(file)
	}

	c.sftp.RemoveAll(file)
	return nil
}

func (c *Client) ReadFile(targetFile string) ([]byte, error) {
	stat, err := c.sftp.Stat(targetFile)
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("cannot read directory %s", targetFile)
	}

	f, err := c.sftp.Open(targetFile)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *Client) SystemInfo() {

}

type Status string

var (
	StatusNone       Status = ""
	StatusActive     Status = "active"
	StatusInactive   Status = "inactive"
	StatusUnknown    Status = "unknown"
	StatusActivating Status = "activating"
)

func (c *Client) SystemdStatus(svcName string) (Status, error) {
	output, err := c.RunCommand([]string{"systemctl", "is-active", svcName})
	if err != nil {
		c.log.Errorf("failed to get status of %s: %v", svcName, err)
		return StatusNone, err
	}
	return Status(strings.TrimSpace(string(output))), nil
}
