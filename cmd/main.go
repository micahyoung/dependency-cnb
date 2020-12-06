package main

import (
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/cargo"
	"github.com/paketo-buildpacks/packit/postal"
	"github.com/paketo-buildpacks/packit/scribe"
	"os"

	"github.com/micahyoung/dependency-cnb/dep"
)

func main() {
	packit.Run(
		dep.Detect(),
		dep.Build(
			scribe.NewLogger(os.Stdout),
			postal.NewService(cargo.NewTransport()),
		),
	)
}
