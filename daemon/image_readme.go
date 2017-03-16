package daemon

import (
	"io"

	"github.com/docker/docker/image/tarexport"
)

// ExportImage exports a list of images to the given output stream. The
// exported images are archived into a tar when written to the output
// stream. All images with the given tag and all versions containing
// the same tag are exported. names is the set of tags to export, and
// outStream is the writer which the images are written to.
func (daemon *Daemon) ReadmeImage(names []string, outStream io.Writer) error {
	imageExporter := tarexport.NewTarExporter(daemon.imageStore, daemon.layerStore, daemon.referenceStore, daemon)
	return imageExporter.Readme(names, outStream)
}

