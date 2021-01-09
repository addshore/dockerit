package cmd

import (
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
	Image           string
	Pull            bool
	Cmd             strslice.StrSlice
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

	if(options.Pull){
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

	// TODO working directory
	// TODO ports
	// TODO volumes
	// TODO labels?
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: options.Image,
			Cmd: options.Cmd,
			AttachStderr:true,
			AttachStdin: true,
			Tty:         true,
			AttachStdout:true,
			OpenStdin:   true,
		},
		&container.HostConfig{
		}, nil, nil, "")
	if err != nil {
        fmt.Println("Error Creating")
		panic(err)
	}

	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
        fmt.Println("Error Starting")
        panic(err)
	}

	waiter, err := cli.ContainerAttach(context.Background(), cont.ID, types.ContainerAttachOptions{
        Stderr:       true,
        Stdout:       true,
        Stdin:        true,
        Stream:       true,
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