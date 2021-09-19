package main

import (
    "fmt"
    "os"
    "os/exec"
    "syscall"
    "path/filepath"
    "strconv"
    "io/ioutil"
)

func main() {
    switch os.Args[1] {
    case "run":
        run()
    case "ns":
        ns()
    default:
        panic("help:\nrun /bin/bash\nrun echo hello world")
    }
}

func run() {
    fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

    cmd := exec.Command("/proc/self/exe", append([]string{"ns"}, os.Args[2:]...)...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    cmd.SysProcAttr = &syscall.SysProcAttr {
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
        Unshareflags: syscall.CLONE_NEWNS,
    }

    must(cmd.Run())
}

func ns() {
    fmt.Printf("Running %v in a new UTS namespace as PID %d\n", os.Args[2:], os.Getpid())

    cg()

    must(syscall.Sethostname([]byte("inside-container")))
    must(syscall.Chroot("/rootfs"))
    must(os.Chdir("/"))
    must(syscall.Mount("proc", "proc", "proc", 0, ""))

    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    must(cmd.Run())

    syscall.Unmount("/proc", 0)
}

func cg() {
    cgroups := "/sys/fs/cgroup/"
    pids := filepath.Join(cgroups, "pids")
    os.Mkdir(filepath.Join(pids, "ourContainer"), 0755)
    ioutil.WriteFile(filepath.Join(pids, "ourContainer/pids.max"), []byte("10"), 0700)
    //up here we limit the number of child processes to 10

    ioutil.WriteFile(filepath.Join(pids, "ourContainer/notify_on_release"), []byte("1"), 0700)

    ioutil.WriteFile(filepath.Join(pids, "ourContainer/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)
    // up here we write container PIDs to cgroup.procs
}

func must(err error) {
    if err != nil {
        panic(err)
    }
}
