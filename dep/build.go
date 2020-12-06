package dep

import (
	"fmt"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/postal"
	"github.com/paketo-buildpacks/packit/scribe"
	"path/filepath"
)

//DependencyService for packit
type DependencyService interface {
	Resolve(path, id, version, stack string) (postal.Dependency, error)
	Install(dependency postal.Dependency, cnbPath, layerPath string) error
}

//Build for packit
func Build(logger scribe.Logger, service DependencyService) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		buildpackID := context.BuildpackInfo.ID
		buildpackName := context.BuildpackInfo.Name

		dependencyLayer, err := context.Layers.Get("my-layer", packit.BuildLayer, packit.CacheLayer, packit.LaunchLayer)
		if err != nil {
			return packit.BuildResult{}, err
		}

		logger.Process(fmt.Sprintf("Providing %s", buildpackName))
		if len(filepath.Join(dependencyLayer.Path, "*")) != 0 {
			dep, err := service.Resolve(filepath.Join(context.CNBPath, "buildpack.toml"), buildpackID, "", context.Stack)
			if err != nil {
				return packit.BuildResult{}, err
			}

			logger.Subprocess(fmt.Sprintf("Downloading %s", buildpackName))
			err = service.Install(dep, context.CNBPath, dependencyLayer.Path)
			if err != nil {
				return packit.BuildResult{}, err
			}
		} else {
			logger.Subprocess(fmt.Sprintf("Reusing cached %s", buildpackName))
		}

		return packit.BuildResult{
			Plan: context.Plan,
			Layers: []packit.Layer{
				dependencyLayer,
			},
		}, nil
	}
}

