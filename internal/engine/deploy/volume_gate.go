package deploy

import (
	"fmt"
	"sync"
)

type VolumeOperations interface{ AcquireVolume(string) (func(), error) }

type VolumeGate struct {
	mu     sync.Mutex
	active map[string]bool
}

func NewVolumeGate() *VolumeGate { return &VolumeGate{active: make(map[string]bool)} }
func (g *VolumeGate) AcquireVolume(name string) (func(), error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.active[name] {
		return nil, fmt.Errorf("volume %s is busy; wait for its current operation to finish", name)
	}
	g.active[name] = true
	var once sync.Once
	return func() { once.Do(func() { g.mu.Lock(); delete(g.active, name); g.mu.Unlock() }) }, nil
}
func (d *DatabaseDeployer) acquireVolume(id string) (func(), error) {
	if d.volumes == nil {
		return func() {}, nil
	}
	return d.volumes.AcquireVolume("codedock-db-data-" + id)
}
func (d *DatabaseDeployer) SetVolumeOperations(volumes VolumeOperations) { d.volumes = volumes }
