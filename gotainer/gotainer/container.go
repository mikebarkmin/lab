package gotainer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/google/uuid"
)

type Container struct {
	id    uuid.UUID
	image *Image
}

func (c *Container) Exec(commands []string) error {
	childCMD := []string{"child", c.GetID()}
	cmd := exec.Command("/proc/self/exe", append(childCMD, commands...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// setup a new namespace which will be used for this cmd
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER,
		Credential: &syscall.Credential{Uid: 0, Gid: 0},
		UidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getuid(), Size: 1},
		},
		GidMappings: []syscall.SysProcIDMap{
			{ContainerID: 0, HostID: os.Getgid(), Size: 1},
		},
	}

	return cmd.Run()
}

func (c *Container) Remove() error {
	return os.RemoveAll(c.getPath())
}

func (c *Container) getPath() string {
	return filepath.Join(containerPath, c.GetID())
}

func (c *Container) Child(commands []string) {
	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot(c.getPath())) // need to change to a rootfs
	cmd := exec.Command(commands[0], commands[1:]...)
	fmt.Println(cmd.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(cmd.Run())
}

func (c *Container) GetID() string {
	return c.id.String()
}

func List() []*Container {
	return nil
}

func GetContainer(id string) (*Container, error) {
	fileName := id
	containerPath := filepath.Join(containerPath, fileName)
	if f, err := os.Stat(containerPath); err == nil && f.IsDir() {
		i, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		return NewContainer(i, nil), nil
	}
	return nil, ErrImageNotExist
}

func NewContainer(id uuid.UUID, image *Image) *Container {
	c := &Container{id, image}
	return c
}
