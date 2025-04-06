package chart

import (
	"ksetup/pkg/config"
	"time"

	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
)

type Object struct {
	Name            string
	ReleaseName     string
	Version         string
	Namespace       string
	CreateNamespace bool
	ChartPkg        string
	Timeout         time.Duration
	ValueFile       string
	Values          map[string]interface{}
	RepoURL         string
	RepoUser        string
	RepoPass        string
}

func ToObject(chart config.Chart) *Object {
	obj := &Object{
		Name:      chart.Name,
		Namespace: chart.Namespace,
	}

	return obj
}

type Client struct {
	log *logrus.Entry
}

func New() (*Client, error) {
	return nil, nil
}

func (c *Client) actionConfig() (*action.Configuration, error) {
	return nil, nil
}

func (c *Client) Install(obj *Object) error {
	ac, err := c.actionConfig()
	if err != nil {
		return err
	}

	install := action.NewInstall(ac)
	install.ReleaseName = obj.ReleaseName
	install.Namespace = obj.Namespace
	install.CreateNamespace = obj.CreateNamespace

	if obj.Timeout > 0 {
		install.Wait = true
		install.Timeout = obj.Timeout
	}

	c.log.Infof("chart %s is installing", obj.Name)

	chart, err := loader.Load(obj.ChartPkg)
	if err != nil {
		return err
	}
	if obj.ValueFile != "" {
		values, err := chartutil.ReadValuesFile(obj.ValueFile)
		if err != nil {
			return err
		}

		newValues, err := chartutil.MergeValues(chart, values)
		if err != nil {
			return err
		}

		chart.Values = newValues
	}

	if obj.Values != nil {
		newValues, err := chartutil.MergeValues(chart, obj.Values)
		if err != nil {
			return err
		}

		chart.Values = newValues
	}

	_, err = install.Run(chart, chart.Values)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Uninstall(obj *Object) error {
	ac, err := c.actionConfig()
	if err != nil {
		return err
	}

	uninstall := action.NewUninstall(ac)
	if obj.Timeout > 0 {
		uninstall.Wait = true
		uninstall.Timeout = obj.Timeout
	}
	uninstall.IgnoreNotFound = true

	_, err = uninstall.Run(obj.ReleaseName)
	if err != nil {
		c.log.Errorf("failed to uninstall %s: %v", obj.Name, err)
		return err
	}

	return nil
}

func (c *Client) Upgrade(obj *Object) error {
	ac, err := c.actionConfig()
	if err != nil {
		c.log.Errorf("failed to get actiob config: %v", err)
		return err
	}

	chartPkg, err := loader.Load(obj.ChartPkg)
	if err != nil {
		c.log.Errorf("failed to load chart: %v", err)
		return err
	}

	if obj.ValueFile != "" {
		values, err := chartutil.ReadValuesFile(obj.ValueFile)
		if err != nil {
			c.log.Errorf("failed to read values file: %v", err)
			return err
		}

		newValues, err := chartutil.MergeValues(chartPkg, values)
		if err != nil {
			c.log.Errorf("failed to merge values: %v", err)
			return err
		}

		chartPkg.Values = newValues
	}

	if obj.Values != nil {
		newValues, err := chartutil.MergeValues(chartPkg, obj.Values)
		if err != nil {
			c.log.Errorf("failed to merge values: %v", err)
			return err
		}

		chartPkg.Values = newValues
	}

	// compare chart values
	if true {
		upgrade := action.NewUpgrade(ac)
		upgrade.Namespace = obj.Namespace
		if obj.Timeout > 0 {
			upgrade.Timeout = obj.Timeout
			upgrade.Wait = true
		}
		upgrade.InsecureSkipTLSverify = true

		_, err = upgrade.Run(obj.ReleaseName, chartPkg, chartPkg.Values)
		if err != nil {
			c.log.Errorf("failed to upgrade chart %s: %v", obj.Name, err)
			return err
		}
	}
	return nil
}

func (c *Client) Get(releaseName string) (*Object, error) {
	ac, err := c.actionConfig()
	if err != nil {
		return nil, err
	}

	get := action.NewGet(ac)
	rl, err := get.Run(releaseName)
	if err != nil {
		c.log.Errorf("failed to get release %s: %v", releaseName, err)
		return nil, err
	}

	obj := &Object{
		Name:        rl.Name,
		Namespace:   rl.Namespace,
		ReleaseName: rl.Name,
		Version:     rl.Chart.AppVersion(),
	}
	return obj, nil
}

func (c *Client) Pull(obj *Object) error {
	pull := action.NewPullWithOpts()
	pull.DestDir = ""
	pull.Version = obj.Version
	pull.RepoURL = obj.RepoURL

	pull.Run(obj.Name)

	return nil
}
