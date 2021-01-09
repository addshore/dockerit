package cmd

import (
	"github.com/docker/go-connections/nat"
	"io/ioutil"
	"os/user"
	"github.com/docker/docker/api/types/mount"
	"os/signal"
	"strings"
	"github.com/docker/docker/api/types/strslice"
	"bufio"
	"os"
	"io"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var fPort string
var fPull bool
var fUser string
var fUserMe bool
var fNoEntry bool
var fMountPwd bool
var fMountHome bool

func init() {
	// Defaults
	rootCmd.Flags().StringVarP(&fPort, "port", "", "", "Port mapping <host>:<container> eg. 8080:80")
	rootCmd.Flags().BoolVarP(&fNoEntry, "entry", "", true, "Use the default entrypoint. If entry=0 you must provide one")
	rootCmd.Flags().BoolVarP(&fMountPwd, "pwd", "", false, "Mount the PWD into the container (and set as working directory /pwd)")
	rootCmd.Flags().BoolVarP(&fMountHome, "home", "", false, "Mount the home directory of the user")
	rootCmd.Flags().StringVarP(&fUser, "user", "", "", "User override for the command")
	rootCmd.Flags().BoolVarP(&fUserMe, "me", "", false, "User override for the command, runs as current user")

	// Optional
	rootCmd.Flags().BoolVarP(&fPull, "pull", "", false, "Pull the docker image even if present")
}

// TODO allow port as an easy runtime option as ports may need to be exposed?
type RunNowOptions struct {
	Image		   string
	Pull			bool
	Cmd			 strslice.StrSlice
}

func RunNow(options RunNowOptions) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}
	if Verbose {
		fmt.Println("Docker client loaded");
	}

	var inout chan []byte

	if(fPull){
		pull(cli,options)
	}

	// TODO ports
	// TODO more volumes?
	cont, err := containerCreate(cli, options)

	// Handle Ctrl + C and exit (removing the container)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for range c {
			// sig is a ^C, handle it
			if Verbose {
				fmt.Println("Stopping container");
			}
			cli.ContainerStop( context.Background(), cont.ID, nil )
			os.Exit(1)
		}
	}()

	waiter, err := cli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{
		Stderr:	   true,
		Stdout:	   true,
		Stdin:		true,
		Stream:	   true,
	})

	go io.Copy(os.Stdout, waiter.Reader)
	go io.Copy(os.Stderr, waiter.Reader)

	if Verbose {
		fmt.Println("Starting container");
	}
	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("Error Starting")
		panic(err)
	}

	go io.Copy(waiter.Conn, os.Stdin)

	if err != nil {
		panic(err)
	}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			inout <- []byte(scanner.Text())
		}
	}()

	// Write to docker container
	go func(w io.WriteCloser) {
		for {
			data, ok := <-inout
			//log.Println("Received to send to docker", string(data))
			if !ok {
				fmt.Println("!ok")
				w.Close()
				return
			}

			w.Write(append(data, '\n'))
		}
	}(waiter.Conn)

	statusCh, errCh := cli.ContainerWait(context.Background(), cont.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	if Verbose {
		fmt.Println("Removing container");
	}
	cli.ContainerRemove( context.Background(), cont.ID, types.ContainerRemoveOptions{} )

	return cont.ID, nil
}

func containerCreate(cli *client.Client, options RunNowOptions) (container.ContainerCreateCreatedBody, error) {
	cont, err := containerCreateNoPullFallback(cli, options)
		if err != nil {
			if !strings.Contains(err.Error()," No such image") {
				fmt.Println("Error Creating")
				panic(err)
			}
			// Fallback pulling the image once
			if Verbose {
				fmt.Println("No image, so pulling");
			}
			pull(cli,options);
			return containerCreateNoPullFallback(cli, options)
		}
	return cont, err;
}

func containerCreateNoPullFallback(cli *client.Client, options RunNowOptions) (container.ContainerCreateCreatedBody, error) {
	if Verbose {
		fmt.Println("Creating container");
	}
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	labels := make(map[string]string)
	labels["com.github/addshore/docker-thing/created-app"] = "docker-thing"
	labels["com.github/addshore/docker-thing/created-command"] = "now"

	ContainerConfig := &container.Config{
		Image: options.Image,
		Cmd: options.Cmd,
		AttachStderr:true,
		AttachStdin: true,
		Tty:		 true,
		AttachStdout:true,
		OpenStdin:   true,
		Labels: labels,
	}

	var emptyMountsSliceEntry []mount.Mount
	HostConfig := &container.HostConfig{
		Mounts: emptyMountsSliceEntry,
		AutoRemove: true,
	}

	runAs := fUser
	if(fUserMe) {
		currentUser, err := user.Current()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		runAs = currentUser.Username
	}
	if(len(runAs)>0) {
		usr, err := user.Lookup(runAs)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		ContainerConfig.User = usr.Uid + ":" + usr.Gid
		// TODO die if not running as a known user?
		// TODO check if the home dir actually exists?
		if(fMountHome){
			HostConfig.Mounts = append(
				HostConfig.Mounts,
				mount.Mount{
					Type:   mount.TypeBind,
					Source: usr.HomeDir,
					Target: usr.HomeDir,
				},
			)
		}
	}

	if(len(fPort)>0){
		splits := strings.Split(fPort, ":")
		hostPortString, containerPortString := splits[0], splits[1]
		containerPort := nat.Port(containerPortString+"/tcp")
		ContainerConfig.ExposedPorts = nat.PortSet{
			containerPort: {},
		}
		HostConfig.PortBindings = nat.PortMap{
			containerPort: []nat.PortBinding{
				{
					HostIP: "0.0.0.0",
					HostPort: hostPortString,
				},
			},
		}
	}
	if(fMountPwd){
		ContainerConfig.WorkingDir = "/pwd"
		HostConfig.Mounts = append(
			HostConfig.Mounts,
			mount.Mount{
				Type:   mount.TypeBind,
				Source: pwd,
				Target: "/pwd",
			},
		)
	}
	if(fNoEntry){
		var emptyStrSliceEntry []string
		ContainerConfig.Entrypoint = emptyStrSliceEntry
	}

	return cli.ContainerCreate(
		context.Background(),
		ContainerConfig,
		HostConfig,
		nil,
		nil,
		"",
		);
}

func pull(cli *client.Client, options RunNowOptions) {
	fmt.Println("Pulling image");
	r, err := cli.ImagePull(
		context.Background(),
		options.Image,
		types.ImagePullOptions{},
	)
	if err != nil {
		fmt.Println("Error Pulling")
		panic(err)
	}
	// TODO fixme this is super verbose...
	if Verbose {
		io.Copy(os.Stdout, r)
	} else {
		io.Copy(ioutil.Discard, r)
	}
}
