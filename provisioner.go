package astilectron

import (
	"net/http"
	"os"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// Provisioner represents an object capable of provisioning Astilectron
type Provisioner interface {
	Provision(p *Paths) error
}

// Default provisioner
var DefaultProvisioner = &defaultProvisioner{
	httpClient: &http.Client{},
}

// defaultProvisioner represents the default provisioner
type defaultProvisioner struct {
	httpClient *http.Client // We need to set up our own client in case we need to tweak some options such as timeout or proxy
}

// Provision implements the provisioner interface
func (p *defaultProvisioner) Provision(paths *Paths) (err error) {
	// Provision astilectron
	if err = p.provisionAstilectron(paths); err != nil {
		err = errors.Wrap(err, "provisioning astilectron failed")
		return
	}

	// Provision electron
	if err = p.provisionElectron(paths); err != nil {
		err = errors.Wrap(err, "provisioning electron failed")
		return
	}
	return
}

// provisionAstilectron provisions astilectron
func (p *defaultProvisioner) provisionAstilectron(paths *Paths) error {
	return p.provisionDownloadableZipFile("Astilectron", paths.AstilectronApplication(), paths.AstilectronDownloadSrc(), paths.AstilectronDownloadDst(), paths.AstilectronUnzipSrc(), paths.AstilectronDirectory())
}

// provisionElectron provisions electron
func (p *defaultProvisioner) provisionElectron(paths *Paths) error {
	return p.provisionDownloadableZipFile("Electron", paths.ElectronExecutable(), paths.ElectronDownloadSrc(), paths.ElectronDownloadDst(), paths.ElectronUnzipSrc(), paths.ElectronDirectory())
}

// provisionDownloadableZipFile provisions a downloadable .zip file
func (p *defaultProvisioner) provisionDownloadableZipFile(name, pathExists, pathDownloadSrc, pathDownloadDst, pathUnzipSrc, pathDirectory string) (err error) {
	// Log
	astilog.Debugf("Provisioning %s...", name)

	// We need to provision
	if _, err = os.Stat(pathExists); os.IsNotExist(err) {
		// Download the .zip file
		if err = Download(p.httpClient, pathDownloadDst, pathDownloadSrc); err != nil {
			return errors.Wrapf(err, "downloading %s into %s failed", pathDownloadSrc, pathDownloadDst)
		}

		// Remove previous install
		astilog.Debugf("Removing %s", pathDirectory)
		if err = os.RemoveAll(pathDirectory); err != nil && !os.IsNotExist(err) {
			return errors.Wrapf(err, "removing %s failed", pathDirectory)
		}

		// Unzip
		if err = Unzip(pathDirectory, pathUnzipSrc); err != nil {
			return errors.Wrapf(err, "unzipping %s into %s failed", pathUnzipSrc, pathDirectory)
		}
	} else if err != nil {
		return errors.Wrapf(err, "stating %s failed", pathExists)
	} else {
		astilog.Debugf("%s already exists, skipping %s provision...", pathExists, name)
	}
	return
}
