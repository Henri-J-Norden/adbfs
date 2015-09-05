// TODO: Write more tests.
// TODO: Implement better file read.
package fs

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/stretchr/testify/assert"
	"github.com/zach-klippenstein/goadb"
)

func TestGetAttrRoot(t *testing.T) {
	dev := &MockDeviceClient{
		&MockDirEntry{
			DirEntry: &goadb.DirEntry{
				Name: "/",
				Size: 0,
				Mode: os.ModeDir | 0755,
			},
		},
	}
	fs, err := NewAdbFileSystem(Config{
		Mountpoint:    "",
		ClientFactory: func() DeviceClient { return dev },
	})
	assert.NoError(t, err)

	attr, status := fs.GetAttr("", NewContext(1, 2, 3))
	assert.True(t, status.Ok(), "Expected status to be Ok, but was %s", status)
	assert.NotNil(t, attr)

	assert.Equal(t, uint64(0), attr.Size)
	assert.False(t, attr.IsRegular())
	assert.True(t, attr.IsDir())
	assert.False(t, attr.IsBlock())
	assert.False(t, attr.IsChar())
	assert.False(t, attr.IsFifo())
	assert.False(t, attr.IsSocket())
	assert.False(t, attr.IsSymlink())
	assert.Equal(t, uint32(0755), attr.Mode&uint32(os.ModePerm))
}

func TestGetAttrRegularFile(t *testing.T) {
	dev := &MockDeviceClient{
		&MockDirEntry{
			DirEntry: &goadb.DirEntry{
				Name: "/version.txt",
				Size: 42,
				Mode: 0444,
			},
		},
	}
	fs, err := NewAdbFileSystem(Config{
		Mountpoint:    "",
		ClientFactory: func() DeviceClient { return dev },
	})
	assert.NoError(t, err)

	attr, status := fs.GetAttr("version.txt", NewContext(1, 2, 3))
	assert.True(t, status.Ok(), "Expected status to be Ok, was %s", status)
	assert.NotNil(t, attr)

	assert.Equal(t, uint64(42), attr.Size)
	assert.True(t, attr.IsRegular())
	assert.False(t, attr.IsDir())
	assert.False(t, attr.IsBlock())
	assert.False(t, attr.IsChar())
	assert.False(t, attr.IsFifo())
	assert.False(t, attr.IsSocket())
	assert.False(t, attr.IsSymlink())
	assert.Equal(t, uint32(0444), attr.Mode&uint32(os.ModePerm))
}

func NewContext(uid int, gid int, pid int) *fuse.Context {
	return &fuse.Context{
		Owner: fuse.Owner{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
		Pid: uint32(pid),
	}
}

type MockDeviceClient struct {
	Root *MockDirEntry
}

type MockDirEntry struct {
	*goadb.DirEntry
}

func (d *MockDeviceClient) OpenRead(path string) (io.ReadCloser, error) {
	return nil, nil
}

func (d *MockDeviceClient) Stat(path string) (*goadb.DirEntry, error) {
	if path == d.Root.Name {
		return d.Root.DirEntry, nil
	}
	return nil, fmt.Errorf("Path does not exist: %s", path)
}

func (d *MockDeviceClient) ListDirEntries(path string) (*goadb.DirEntries, error) {
	return nil, nil
}

func (d *MockDeviceClient) RunCommand(cmd string, args ...string) (string, error) {
	return "", nil
}