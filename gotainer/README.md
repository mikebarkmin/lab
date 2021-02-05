# Container From Scratch

## Namespaces

> What you can see

- UNIX Timesharing System
- Process IDs
- File system (mount points)
- Users
- IPC
- Networking

## Control groups

> What you can use

- CPU
- Memory
- Disk I/O
- Network
- Device permissions (/dev)

CGroups can not be used for rootless containers.

Creating a CGroup can be done by creating a directory in `/sys/fs/cgroup/[group]`.

## Images

An image is a file system, which will be used here:
`must(syscall.Chroot("/rootfs")) // need to change to a rootfs`

### Setup a RootFS

```bash
docker export $(docker create ubuntu) --output="latest.tar"
mkdir rootfs
tar xvf latest.tar -C rootfs
rm latest.tar
```

## Containers

> Containers are processes running inside the isolation layer but are visible from the host — https://youtu.be/_TsSmSu57Zo?t=842

## Rootless Containers

> You are not allowed to create the namespaces you need for isolation as an unprivileged user — https://youtu.be/_TsSmSu57Zo?t=959

- syscall.CLONE_NEWN: Allows an unprivileged user to mount

## References

- https://www.youtube.com/watch?v=Utf-A4rODH8
- https://www.youtube.com/watch?v=_TsSmSu57Zo
