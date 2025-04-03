package kube

import (
	"bytes"
	"io"
	"ksetup/pkg/client/chart"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/disk"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/delete"
	"k8s.io/kubectl/pkg/cmd/util"
)

type Client struct {
	dynamic    dynamic.Interface
	restConfig *rest.Config
	log        *logrus.Logger
}

func New(masterUrl string, kubeConfig []byte, log *logrus.Logger) (*Client, error) {

	apiConfig, err := clientcmd.Load(kubeConfig)
	if err != nil {
		return nil, err
	}

	kubeConfigGetter := func() (*clientcmdapi.Config, error) {
		return apiConfig, nil
	}
	config, err := clientcmd.BuildConfigFromKubeconfigGetter(masterUrl, kubeConfigGetter)
	if err != nil {
		return nil, err
	}

	dynamicCli, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &Client{
		dynamic: dynamicCli,
		log:     log,
	}, nil
}

func (c *Client) NewChartClient() (*chart.Client, error) {
	return nil, nil
}

func (c *Client) ToRESTConfig() (*rest.Config, error) {
	return c.restConfig, nil
}

func (c *Client) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	disk.NewCachedDiscoveryClientForConfig(nil, "", "", 10*time.Second)
	return nil, nil
}

func (c *Client) ToRESTMapper() (meta.RESTMapper, error) {
	return nil, nil
}

func (c *Client) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return nil
}

func (c *Client) Apply(manifest string) {
	data, err := os.ReadFile(manifest)
	if err != nil {
		return
	}
	buffer := bytes.NewBuffer(data)
	factory := util.NewFactory(c)
	stream := genericiooptions.IOStreams{
		In:     buffer,
		Out:    os.Stdout,
		ErrOut: io.Discard,
	}
	applyCmd := apply.NewCmdApply("", factory, stream)
	applyCmd.Execute()
}

func (c *Client) Delete(manifest string) {
	data, err := os.ReadFile(manifest)
	if err != nil {
		return
	}
	buffer := bytes.NewBuffer(data)
	stream := genericiooptions.IOStreams{
		In:     buffer,
		Out:    os.Stdout,
		ErrOut: io.Discard,
	}
	deleteCmd := delete.NewCmdDelete(util.NewFactory(c), stream)
	deleteCmd.Execute()
}
