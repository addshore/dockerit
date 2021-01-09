package cmd

import (
	"os/signal"
	"strings"
	"github.com/docker/docker/api/types/strslice"
	"bufio"
	"os"
	"io"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func init() {
	rootCmd.AddCommand(nowCmd)
}

// TODO allow port as an easy runtime option as ports may need to be exposed?
type RunNowOptions struct {
	Image		   string
	Pull			bool
	Cmd			 strslice.StrSlice
}

var nowCmd = &cobra.Command{
	Use:   "now",
	Run: func(cmd *cobra.Command, args []string) {
		RunNow(RunNowOptions{
			Image: args[0],
			Pull: false,
			Cmd: args[1:],
		})
		},
	}

func RunNow(options RunNowOptions) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println("Unable to create docker client")
		panic(err)
	}

	var inout chan []byte

	// TODO optionally Pull
	//pull(cli,options);

	// TODO working directory
	// TODO ports
	// TODO volumes
	// TODO labels?
	cont, err := containerCreate(cli, options)

	// Handle Ctrl + C and exit (removing the container)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func(){
		for range c {
			// sig is a ^C, handle it
			cli.ContainerRemove( context.Background(), cont.ID, types.ContainerRemoveOptions{
				Force: true,
				} )
			os.Exit(1)
		}
	}()

	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Println("Error Starting")
		panic(err)
	}

	waiter, err := cli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{
		Stderr:	   true,
		Stdout:	   true,
		Stdin:		true,
		Stream:	   true,
	})

	go  io.Copy(os.Stdout, waiter.Reader)
	go  io.Copy(os.Stderr, waiter.Reader)
	//go io.Copy(aResp.Conn, os.Stdin)

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
			pull(cli,options);
			return containerCreateNoPullFallback(cli, options)
		}
	return cont, err;
}

func containerCreateNoPullFallback(cli *client.Client, options RunNowOptions) (container.ContainerCreateCreatedBody, error) {
	return cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: options.Image,
			Cmd: options.Cmd,
			AttachStderr:true,
			AttachStdin: true,
			Tty:		 true,
			AttachStdout:true,
			OpenStdin:   true,
		},
		&container.HostConfig{
		}, nil, nil, "");
}

func pull(cli *client.Client, options RunNowOptions) {
	r, err := cli.ImagePull(
		context.Background(),
		options.Image,
		types.ImagePullOptions{},
	)
	if err != nil {
		panic(err)
	}
	// TODO fixme this is super verbose...
	fmt.Println("Error Pulling")
	io.Copy(os.Stdout, r)
}
