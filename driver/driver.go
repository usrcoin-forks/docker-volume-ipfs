package driver

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/go-plugins-helpers/volume"
)

type IPFS struct {
	mountPoint string
	volumes    map[string]string
	m          *sync.Mutex
}

func New(ipfsMountPoint string) IPFS {
	d := IPFS{
		mountPoint: ipfsMountPoint,
		volumes:    make(map[string]string),
		m:          &sync.Mutex{},
	}
	return d
}

func (d IPFS) Create(r volume.Request) volume.Response {
	fmt.Printf("Create %v\n", r)

	d.m.Lock()
	defer d.m.Unlock()

	volumeName := r.Name

	if _, ok := d.volumes[volumeName]; ok {
		return volume.Response{}
	}

	volumePath := filepath.Join(d.mountPoint, volumeName)

	_, err := os.Lstat(volumePath)
	if err != nil {
		fmt.Println("Error", volumePath, err.Error())
		return volume.Response{Err: fmt.Sprintf("Error while looking for volumePath %s: %s", volumePath, err.Error())}
	}

	d.volumes[volumeName] = volumePath

	return volume.Response{}
}

func (d IPFS) Path(r volume.Request) volume.Response {
	fmt.Printf("Path %v\n", r)
	fmt.Printf("%v", d.volumes)
	volumeName := r.Name

	if volumePath, ok := d.volumes[volumeName]; ok {
		return volume.Response{Mountpoint: volumePath}
	}

	return volume.Response{}
}

func (d IPFS) Remove(r volume.Request) volume.Response {
	fmt.Printf("Remove %v", r)

	d.m.Lock()
	defer d.m.Unlock()

	volumeName := r.Name

	if _, ok := d.volumes[volumeName]; ok {
		delete(d.volumes, volumeName)
	}

	return volume.Response{}
}

func (d IPFS) Mount(r volume.MountRequest) volume.Response {
	fmt.Printf("Mount %v\n", r)
	volumeName := r.Name

	if volumePath, ok := d.volumes[volumeName]; ok {
		return volume.Response{Mountpoint: volumePath}
	}

	return volume.Response{}
}

func (d IPFS) Unmount(r volume.UnmountRequest) volume.Response {
	fmt.Printf("Unmount %v: nothing to do\n", r)
	return volume.Response{}
}

func (d IPFS) Get(r volume.Request) volume.Response {
	fmt.Printf("Get %v\n", r)
	volumeName := r.Name

	if volumePath, ok := d.volumes[volumeName]; ok {
		return volume.Response{Volume: &volume.Volume{Name: volumeName, Mountpoint: volumePath}}
	}

	return volume.Response{Err: fmt.Sprintf("volume %s does not exists", volumeName)}
}

func (d IPFS) List(r volume.Request) volume.Response {
	fmt.Printf("List %v\n", r)

	volumes := []*volume.Volume{}

	for name, path := range d.volumes {
		volumes = append(volumes, &volume.Volume{Name: name, Mountpoint: path})
	}

	return volume.Response{Volumes: volumes}
}

func (d IPFS) Capabilities(r volume.Request) volume.Response {
	// FIXME(vdemeester) handle capabilities better
	return volume.Response{
		Capabilities: volume.Capability{
			Scope: "local",
		},
	}
}
