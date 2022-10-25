package htfs

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pathlib"
	"github.com/robocorp/rcc/pretty"
)

type unmanaged struct {
	delegate  MutableLibrary
	path      string
	resolved  bool
	protected bool
}

func Unmanaged(core MutableLibrary) MutableLibrary {
	return &unmanaged{
		delegate:  core,
		path:      "",
		resolved:  false,
		protected: false,
	}
}

func (it *unmanaged) Identity() string {
	return it.delegate.Identity()
}

func (it *unmanaged) Stage() string {
	return it.delegate.Stage()
}

func (it *unmanaged) CatalogPath(key string) string {
	return "Unmanaged Does Not Support Catalog Path Request"
}

func (it *unmanaged) Remove([]string) error {
	return fmt.Errorf("Not supported yet on unmanaged holotree.")
}

func (it *unmanaged) Export([]string, string) error {
	return fmt.Errorf("Not supported yet on unmanaged holotree.")
}

func (it *unmanaged) resolve(blueprint []byte) error {
	if it.resolved {
		return nil
	}
	controller := []byte(common.ControllerIdentity())
	space := []byte(common.HolotreeSpace)
	path, err := it.TargetDir(blueprint, controller, space)
	if err != nil {
		return nil
	}
	if !pathlib.Exists(path) {
		it.path = path
		it.resolved = true
		return nil
	}
	identityfile := filepath.Join(path, "identity.yaml")
	_, identity, err := ComposeFinalBlueprint([]string{identityfile}, "")
	if err != nil {
		return nil
	}
	expected := BlueprintHash(blueprint)
	actual := BlueprintHash(identity)
	if actual != expected {
		return fmt.Errorf("Unmanaged fingerprint %q does not match requested one %q! Quitting!", actual, expected)
	}
	it.path = path
	it.protected = true
	it.resolved = true
	return nil
}

func (it *unmanaged) ValidateBlueprint(blueprint []byte) error {
	err := it.resolve(blueprint)
	if err != nil {
		return err
	}
	if it.protected {
		return nil
	}
	return it.delegate.ValidateBlueprint(blueprint)
}

func (it *unmanaged) Record(blueprint []byte) error {
	it.resolve(blueprint)
	if it.protected {
		common.Timeline("holotree unmanaged record prevention")
		return nil
	}
	return it.delegate.Record(blueprint)
}

func (it *unmanaged) TargetDir(blueprint, client, tag []byte) (string, error) {
	return it.delegate.TargetDir(blueprint, client, tag)
}

func (it *unmanaged) Restore(blueprint, client, tag []byte) (string, error) {
	defer common.Log("%sThis is unmanaged holotree space for blueprint: %v%s", pretty.Magenta, BlueprintHash(blueprint), pretty.Reset)
	it.resolve(blueprint)
	if !it.protected {
		return it.delegate.Restore(blueprint, client, tag)
	}
	common.Timeline("holotree unmanaged restore prevention")
	if len(it.path) > 0 {
		return it.path, nil
	}
	return "", fmt.Errorf("Unmanaged path resolution failed!")
}

func (it *unmanaged) Open(digest string) (readable io.Reader, closer Closer, err error) {
	return it.delegate.Open(digest)
}

func (it *unmanaged) ExactLocation(key string) string {
	return it.delegate.ExactLocation(key)
}

func (it *unmanaged) Location(key string) string {
	return it.delegate.Location(key)
}

func (it *unmanaged) HasBlueprint(blueprint []byte) bool {
	return it.delegate.HasBlueprint(blueprint)
}
