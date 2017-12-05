package stopwatchapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	astilectron "github.com/asticode/go-astilectron"
	"github.com/pkg/errors"
)

const astiMainTargetID = "main"

// Provision event names
const (
	EventNameProvisionAstilectronAlreadyProvisioned = "provision.astilectron.already.provisioned"
	EventNameProvisionAstilectronFinished           = "provision.astilectron.finished"
	EventNameProvisionAstilectronMoved              = "provision.astilectron.moved"
	EventNameProvisionAstilectronUnzipped           = "provision.astilectron.unzipped"
	EventNameProvisionElectronAlreadyProvisioned    = "provision.electron.already.provisioned"
	EventNameProvisionElectronFinished              = "provision.electron.finished"
	EventNameProvisionElectronMoved                 = "provision.electron.moved"
	EventNameProvisionElectronUnzipped              = "provision.electron.unzipped"
)

// Provision event names mapping keys
const (
	provisionEventNamesMappingKeyAlreadyProvisioned = "already.provisioned"
	provisionEventNamesMappingKeyFinished           = "finished"
	provisionEventNamesMappingKeyMoved              = "moved"
	provisionEventNamesMappingKeyUnzipped           = "unzipped"
)

var provisionEventNamesMapping = map[string]map[string]string{
	"Astilectron": {
		provisionEventNamesMappingKeyAlreadyProvisioned: EventNameProvisionAstilectronAlreadyProvisioned,
		provisionEventNamesMappingKeyFinished:           EventNameProvisionAstilectronFinished,
		provisionEventNamesMappingKeyMoved:              EventNameProvisionAstilectronMoved,
		provisionEventNamesMappingKeyUnzipped:           EventNameProvisionAstilectronUnzipped,
	},
	"Electron": {
		provisionEventNamesMappingKeyAlreadyProvisioned: EventNameProvisionElectronAlreadyProvisioned,
		provisionEventNamesMappingKeyFinished:           EventNameProvisionElectronFinished,
		provisionEventNamesMappingKeyMoved:              EventNameProvisionElectronMoved,
		provisionEventNamesMappingKeyUnzipped:           EventNameProvisionElectronUnzipped,
	},
}

// VendorVersions represents the provision status
type VendorVersions struct {
	Astilectron *VendorVersionsPackage            `json:"astilectron,omitempty"`
	Electron    map[string]*VendorVersionsPackage `json:"electron,omitempty"`
}

// VendorVersionsPackage represents the provision status of a package
type VendorVersionsPackage struct {
	Version string `json:"version"`
}

// provisionElectronKey returns the electron's provision status key
func provisionElectronKey() string {
	return fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
}

// provisioner implements the astilectron.Provisioner interface.
type provisioner struct {
	vendorDir    string
	versionsPath string
	restoreFn    AssetRestoreFn
}

// NewProvisioner returns an initialized provisioner that will unpack the
// embedded vendor runtime resources to their correct locations.
func NewProvisioner(dest string, restoreFn AssetRestoreFn) astilectron.Provisioner {
	p := &provisioner{
		vendorDir: path.Join(dest, "vendor"),
		restoreFn: restoreFn,
	}

	p.versionsPath = path.Join(p.vendorDir, "versions.json")
	return p
}

// Provision implements the provisioner interface
func (p *provisioner) Provision(ctx context.Context, appName, os, arch string, paths astilectron.Paths) (err error) {
	// Retrieve provision status
	var curVer VendorVersions
	if curVer, err = p.VendorVersions(); err != nil {
		err = errors.Wrap(err, "retrieving provisioning status failed")
		return
	}
	defer p.updateVendorVersions(&curVer)

	// Get versions that should be used.
	var reqVer = BundledVendorVersions()

	// provision astilectron package
	if err = p.provisionAstilectron(ctx, curVer.Astilectron, reqVer.Astilectron); err != nil {
		err = errors.Wrap(err, "provisioning astilectron failed")
		return
	}
	curVer.Astilectron = reqVer.Astilectron

	// Provision electron package
	if err = p.provisionElectron(ctx, curVer.Electron[provisionElectronKey()], reqVer.Electron[provisionElectronKey()]); err != nil {
		err = errors.Wrap(err, "provisioning electron failed")
		return
	}
	curVer.Electron[provisionElectronKey()] = reqVer.Electron[provisionElectronKey()]

	return
}

// VendorVersions reads current provision versions from vendor destination path.
func (p *provisioner) VendorVersions() (s VendorVersions, err error) {
	var f *os.File

	s.Electron = make(map[string]*VendorVersionsPackage)
	if f, err = os.Open(p.versionsPath); err != nil {
		if !os.IsNotExist(err) {
			err = errors.Wrapf(err, "opening file %s failed", f.Name())
		} else {
			err = nil
		}
		return
	}

	// If decoding fails, destroy file and force regeneration of vendor assets
	if _err := json.NewDecoder(f).Decode(&s); _err != nil {
		log.Println(errors.Wrapf(_err, "failure decoding status file %s", f.Name()))
		log.Printf("Removing broken status file...")
		if _err = os.Remove(f.Name()); _err != nil {
			err = errors.Wrapf(err, "unable to remove broken status file %s", f.Name())
		}
	}
	return
}

// updates vendor status file with current values
func (p *provisioner) updateVendorVersions(v *VendorVersions) error {
	if v == nil {
		return nil
	}

	var f *os.File
	var err error

	if f, err = os.Create(p.versionsPath); err != nil {
		return errors.Wrapf(err, "failed to create versions file %s", p.versionsPath)
	}

	if err = json.NewEncoder(f).Encode(v); err != nil {
		err = errors.Wrapf(err, "unable to encode new version info to file %s", f.Name())
	}

	return err
}

func (p *provisioner) provisionAstilectron(ctx context.Context, curVer *VendorVersionsPackage, reqVer *VendorVersionsPackage) error {
	var dest = path.Join(p.vendorDir, "astilectron")

	// test if provisioned
	if p.isProvisioned(curVer, reqVer, dest) {
		log.Printf("Astilectron already provisioned, skipping...")
		//d.Dispatch(astilectron.Event{Name: provisionEventNamesMapping[appName][provisionEventNamesMappingKeyAlreadyProvisioned], TargetID: astiMainTargetID})
		return nil
	}

	// Provision.
	if err := p.provisionPackage(ctx, "Astilectron", curVer, reqVer, EmbeddedAstilectronPath(), p.vendorDir); err != nil {
		return err
	}

	// rename the unpacked dir if it exists. The astilectron zip contains a
	// directory named after the version, we want to use the dir without version.
	versionedPath := path.Join(p.vendorDir, "astilectron-"+strings.TrimPrefix(astilectronVersion, "v"))
	if s, err := os.Stat(versionedPath); err == nil && s.IsDir() {
		return os.Rename(versionedPath, path.Join(p.vendorDir, "astilectron"))
	}

	return nil
}

func (p *provisioner) provisionElectron(ctx context.Context, curVer *VendorVersionsPackage, reqVer *VendorVersionsPackage) error {
	var dest = path.Join(p.vendorDir, "electron-"+provisionElectronKey())

	return p.provisionPackage(ctx, "Electron", curVer, reqVer, EmbeddedElectronPath(), dest)
}

func (p *provisioner) isProvisioned(curVer *VendorVersionsPackage, reqVer *VendorVersionsPackage, dest string) bool {
	// check that the destination dir is found, if not, we have to provision,
	// regardless of the versions.
	if _, err := os.Stat(dest); err != nil && os.IsNotExist(err) {
		return false
	}

	return curVer != nil && curVer.Version == reqVer.Version
}

func (p *provisioner) provisionPackage(ctx context.Context, name string, curVer *VendorVersionsPackage, reqVer *VendorVersionsPackage, src, dest string) error {
	// check if provisioned already
	if p.isProvisioned(curVer, reqVer, dest) {
		log.Printf("%s already provisioned, skipping...", name)
		//d.Dispatch(astilectron.Event{Name: provisionEventNamesMapping[name][provisionEventNamesMappingKeyAlreadyProvisioned], TargetID: astiMainTargetID})
		return nil
	}
	log.Printf("Need to prepare %s, please wait...", name)

	// delete previous installation if the destination exists.
	if fi, err := os.Stat(dest); err == nil && fi != nil && fi.IsDir() {
		// log.Printf("Removing previous installation from %s", dest)
		if err = os.RemoveAll(dest); err != nil {
			return errors.Wrapf(err, "unable to remove previous installation: %s", dest)
		}
	}

	// Make destination dir
	// log.Printf("Generating installation directory to %s", dest)
	if err := os.MkdirAll(dest, 0700); err != nil {
		return errors.Wrapf(err, "unable to create destination for installation: %s", dest)
	}

	// Read the package source from memory to disk
	// log.Printf("Unloading package %s to %s", src, dest)
	if err := p.restoreFn(path.Join(dest, ".tmp"), src); err != nil {
		return errors.Wrapf(err, "unable to unpack embedded resource %s into %s", src, dest)
	}
	defer cleanProvisionerTmpData(path.Join(dest, ".tmp"))

	// Generate the path for the unpacked .zip
	zipPath := path.Join(path.Join(dest, ".tmp"), src)

	// log.Printf("Unzipping package content to %s from %s", dest, zipPath)
	// Unpack the zip
	if err := Unpack(zipPath, dest); err != nil {
		return errors.Wrapf(err, "unable to unzip embedded resources to %s", dest)
	}

	// Notify that we're progressing
	//d.Dispatch(astilectron.Event{Name: provisionEventNamesMapping[name][provisionEventNamesMappingKeyUnzipped], TargetID: astiMainTargetID})

	// Dispatch info that provisioning is done
	//d.Dispatch(astilectron.Event{Name: provisionEventNamesMapping[name][provisionEventNamesMappingKeyFinished], TargetID: astiMainTargetID})
	return nil
}

func cleanProvisionerTmpData(tmpdir string) (err error) {
	if tmpdir != "" {
		// log.Printf("Removing temporary dir %s", tmpdir)
		err = os.RemoveAll(tmpdir)
	}
	return
}
