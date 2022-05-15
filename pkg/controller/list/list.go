package list

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aquaproj/aqua/pkg/config"
	finder "github.com/aquaproj/aqua/pkg/config-finder"
	reader "github.com/aquaproj/aqua/pkg/config-reader"
	registry "github.com/aquaproj/aqua/pkg/install-registry"
	"github.com/aquaproj/aqua/pkg/validate"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	stdout            io.Writer
	configFinder      finder.ConfigFinder
	configReader      reader.ConfigReader
	registryInstaller registry.Installer
}

func NewController(configFinder finder.ConfigFinder, configReader reader.ConfigReader, registInstaller registry.Installer) *Controller {
	return &Controller{
		stdout:            os.Stdout,
		configFinder:      configFinder,
		configReader:      configReader,
		registryInstaller: registInstaller,
	}
}

func (ctrl *Controller) List(ctx context.Context, param *config.Param, logE *logrus.Entry) error {
	cfg := &config.Config{}
	cfgFilePath, err := ctrl.configFinder.Find(param.PWD, param.ConfigFilePath)
	if err != nil {
		return err //nolint:wrapcheck
	}

	if err := ctrl.configReader.Read(cfgFilePath, cfg); err != nil {
		return err //nolint:wrapcheck
	}

	if err := validate.Config(cfg); err != nil {
		return fmt.Errorf("configuration is invalid: %w", err)
	}

	registryContents, err := ctrl.registryInstaller.InstallRegistries(ctx, cfg, cfgFilePath, logE)
	if err != nil {
		return err //nolint:wrapcheck
	}
	for registryName, registryContent := range registryContents {
		for _, pkgInfo := range registryContent.PackageInfos {
			fmt.Fprintln(ctrl.stdout, registryName+","+pkgInfo.GetName())
		}
	}

	return nil
}
