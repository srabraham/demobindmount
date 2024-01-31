package main

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"log"
	"os"
	"path/filepath"
)

func fate(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// This is a demo for my comment here: https://github.com/testcontainers/testcontainers-go/issues/2163#issuecomment-1917702951
//
// This program works as expected at testcontainers 0.26.0, but doesn't work at 0.27.0
// You can switch between the two by doing
//
//	go get github.com/testcontainers/testcontainers-go@v0.26.0
//	go get github.com/testcontainers/testcontainers-go@v0.27.0
//
// then run the program with
//
//	go run main.go
func main() {
	ctx := context.Background()

	// this assumes that the working dir is the same as where the go program lives,
	// so it'll work if you `go run main.go`
	pwd, err := os.Getwd()
	fate(err)
	// mount this local file into containerDir
	mountFile := filepath.Join(pwd, "mounts", "somefile.txt")
	containerDir := "/my/container/dir/"
	mounts := []testcontainers.ContainerMount{
		testcontainers.BindMount(
			mountFile,
			testcontainers.ContainerMountTarget(containerDir+filepath.Base(mountFile)),
		),
	}
	// start a container and do an "ls" to check for the file's existence
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:      "alpine:latest",
			Mounts:     mounts,
			Cmd:        []string{"ls", "-l", containerDir},
			WaitingFor: wait.ForExit(),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	fate(err)
	logs, err := container.Logs(ctx)
	fate(err)
	val, err := io.ReadAll(logs)
	fate(err)
	log.Printf("Container logs:\n%s", val)

	fate(container.Terminate(ctx))
}
