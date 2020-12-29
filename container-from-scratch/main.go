package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

var rootless = true

func main() {
	fmt.Println(os.Args[1])
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("what??")
	}
}

// go run main.go run <cmd> <args>
func run() {
	// executes this program with argument child and the rest arguments
	childCMD := []string{"child"}
	cmd := exec.Command("/proc/self/exe", append(childCMD, os.Args[2:]...)...)
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

	must(cmd.Run())
}

func child() {
	fmt.Printf("running %v as PID %d\n", os.Args[2:], os.Getpid())

	if !rootless {
		setupCGroup("mike-test")
	}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot("rootfs")) // need to change to a rootfs
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	err := cmd.Run()

	if err == nil {
		fmt.Printf("stopping %v as PID %d\n", os.Args[2:], os.Getpid())
	}
}

func setupCGroup(name string) {
	cgroups := "/sys/fs/cgroup/"
	// limit of number of process which can be run in this cgroup
	pids := filepath.Join(cgroups, "pids")
	cgroup := filepath.Join(pids, name)
	releaseAgent := "rm -rf " + cgroup
	must(os.Mkdir(cgroup, 0755))
	must(ioutil.WriteFile(filepath.Join(cgroup, "notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(cgroup, "release_agent"), []byte(releaseAgent), 0700))
	must(ioutil.WriteFile(filepath.Join(cgroup, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))

	must(ioutil.WriteFile(filepath.Join(cgroup, "pids.max"), []byte("20"), 0700))

	// limit memory allocation
	memory := filepath.Join(cgroups, "memory")
	cgroup = filepath.Join(memory, name)
	releaseAgent = "rm -rf " + cgroup
	must(os.Mkdir(filepath.Join(cgroup), 0755))
	must(ioutil.WriteFile(filepath.Join(cgroup, "notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(cgroup, "release_agent"), []byte(releaseAgent), 0700))
	must(ioutil.WriteFile(filepath.Join(cgroup, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))

	memoryLimitInBytes := 40 * 1024 * 1024 // 40 MB
	must(ioutil.WriteFile(filepath.Join(cgroup, "memory.limit_in_bytes"), []byte(strconv.Itoa(memoryLimitInBytes)), 0700))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
