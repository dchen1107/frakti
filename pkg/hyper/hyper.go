/*
Copyright 2016 The Kubernetes Authors All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hyper

import (
	"fmt"
	"io"
	"time"

	"github.com/hyperhq/hyperd/lib/promise"
	"github.com/hyperhq/hyperd/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// HyperClient is the gRPC client for hyperd
type HyperClient struct {
	client  types.PublicAPIClient
	ctx     context.Context
	timeout time.Duration
}

// NewHyperClient creates a new *HyperClient
func NewHyperClient(server string, timeout time.Duration) (*HyperClient, error) {
	conn, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &HyperClient{
		client:  types.NewPublicAPIClient(conn),
		ctx:     context.Background(),
		timeout: timeout,
	}, nil
}

// GetPodInfo gets pod info by podID
func (c *HyperClient) GetPodInfo(podID string) (*types.PodInfo, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	request := types.PodInfoRequest{
		PodID: podID,
	}
	pod, err := c.client.PodInfo(ctx, &request)
	if err != nil {
		return nil, err
	}

	return pod.PodInfo, nil
}

// GetPodList get a list of Pods
func (c *HyperClient) GetPodList() ([]*types.PodListResult, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	request := types.PodListRequest{}
	podList, err := c.client.PodList(
		ctx,
		&request,
	)
	if err != nil {
		return nil, err
	}

	return podList.PodList, nil
}

// GetContainerList gets a list of containers
func (c *HyperClient) GetContainerList(auxiliary bool) ([]*types.ContainerListResult, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	req := types.ContainerListRequest{
		Auxiliary: auxiliary,
	}
	containerList, err := c.client.ContainerList(
		ctx,
		&req,
	)
	if err != nil {
		return nil, err
	}

	return containerList.ContainerList, nil
}

// GetContainerInfo gets container info by container name or id
func (c *HyperClient) GetContainerInfo(container string) (*types.ContainerInfo, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	req := types.ContainerInfoRequest{
		Container: container,
	}
	cinfo, err := c.client.ContainerInfo(
		ctx,
		&req,
	)
	if err != nil {
		return nil, err
	}

	return cinfo.ContainerInfo, nil
}

// GetContainerLogs gets container log by container name or id
func (c *HyperClient) GetContainerLogs(container string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	req := types.ContainerLogsRequest{
		Container:  container,
		Follow:     false,
		Timestamps: false,
		Tail:       "",
		Since:      "",
		Stdout:     true,
		Stderr:     true,
	}
	stream, err := c.client.ContainerLogs(
		ctx,
		&req,
	)
	if err != nil {
		return nil, err
	}

	ret := []byte{}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			if req.Follow == true {
				continue
			}
			break
		}
		if err != nil {
			return nil, err
		}
		ret = append(ret, res.Log...)
	}

	return ret, nil
}

// GetImageList gets a list of images
func (c *HyperClient) GetImageList() ([]*types.ImageInfo, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	req := types.ImageListRequest{}
	imageList, err := c.client.ImageList(
		ctx,
		&req,
	)
	if err != nil {
		return nil, err
	}

	return imageList.ImageList, nil
}

// CreatePod creates a pod
func (c *HyperClient) CreatePod(spec *types.UserPod) (string, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	req := types.PodCreateRequest{
		PodSpec: spec,
	}
	resp, err := c.client.PodCreate(
		ctx,
		&req,
	)
	if err != nil {
		return "", err
	}

	return resp.PodID, nil
}

// StartContainer starts a hyper container
func (c *HyperClient) StartContainer(containerID string) error {
	return fmt.Errorf("Not implemented")
}

// StopContainer stops a hyper container
func (c *HyperClient) StopContainer(containerID string, timeout int64) error {
	return fmt.Errorf("Not implemented")
}

// RemoveContainer stops a hyper container
func (c *HyperClient) RemoveContainer(containerID string) error {
	return fmt.Errorf("Not implemented")
}

// CreateContainer creates a container
func (c *HyperClient) CreateContainer(podID string, spec *types.UserContainer) (string, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	req := types.ContainerCreateRequest{
		PodID:         podID,
		ContainerSpec: spec,
	}

	resp, err := c.client.ContainerCreate(ctx, &req)
	if err != nil {
		return "", err
	}

	return resp.ContainerID, nil
}

// RemovePod removes a pod by podID
func (c *HyperClient) RemovePod(podID string) error {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	_, err := c.client.PodRemove(
		ctx,
		&types.PodRemoveRequest{PodID: podID},
	)

	if err != nil {
		return err
	}

	return nil
}

// ContainerExec exec a command in a container with input stream in and output stream out
func (c *HyperClient) ContainerExec(container, tag string, command []string, tty bool, stdin io.ReadCloser, stdout, stderr io.Writer) error {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	request := types.ContainerExecRequest{
		ContainerID: container,
		Command:     command,
		Tag:         tag,
		Tty:         tty,
	}
	stream, err := c.client.ContainerExec(ctx)
	if err != nil {
		return err
	}
	if err := stream.Send(&request); err != nil {
		return err
	}
	var recvStdoutError chan error
	if stdout != nil || stderr != nil {
		recvStdoutError = promise.Go(func() (err error) {
			for {
				in, err := stream.Recv()
				if err != nil && err != io.EOF {
					return err
				}
				if in != nil && in.Stdout != nil {
					nw, ew := stdout.Write(in.Stdout)
					if ew != nil {
						return ew
					}
					if nw != len(in.Stdout) {
						return io.ErrShortWrite
					}
				}
				if err == io.EOF {
					break
				}
			}
			return nil
		})
	}
	if stdin != nil {
		go func() error {
			defer stream.CloseSend()
			buf := make([]byte, 32)
			for {
				nr, err := stdin.Read(buf)
				if nr > 0 {
					if err := stream.Send(&types.ContainerExecRequest{Stdin: buf[:nr]}); err != nil {
						return err
					}
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
			}
			return nil
		}()
	}
	if stdout != nil || stderr != nil {
		if err := <-recvStdoutError; err != nil {
			return err
		}
	}

	return nil
}

// StartPod starts a pod by podID
func (c *HyperClient) StartPod(podID string) error {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	stream, err := c.client.PodStart(ctx)
	if err != nil {
		return err
	}

	req := types.PodStartMessage{
		PodID: podID,
	}
	if err := stream.Send(&req); err != nil {
		return err
	}

	if _, err := stream.Recv(); err != nil {
		return err
	}

	return nil
}

// StopPod stops a pod
func (c *HyperClient) StopPod(podID string) (int, string, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	resp, err := c.client.PodStop(ctx, &types.PodStopRequest{
		PodID: podID,
	})
	if err != nil {
		return -1, "", err
	}

	return int(resp.Code), resp.Cause, nil
}

// Wait gets exitcode by container and processID
func (c *HyperClient) Wait(container, processID string, noHang bool) (int32, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	request := types.WaitRequest{
		Container: container,
		ProcessId: processID,
		NoHang:    noHang,
	}
	response, err := c.client.Wait(ctx, &request)
	if err != nil {
		return -1, err
	}

	return response.ExitCode, nil
}

// PullImage pulls a image from registry
func (c *HyperClient) PullImage(image, tag string, auth *types.AuthConfig, out io.Writer) error {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	request := types.ImagePullRequest{
		Image: image,
		Tag:   tag,
		Auth:  auth,
	}
	stream, err := c.client.ImagePull(ctx, &request)
	if err != nil {
		return err
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if out != nil {
			n, err := out.Write(res.Data)
			if err != nil {
				return err
			}
			if n != len(res.Data) {
				return io.ErrShortWrite
			}
		}
	}

	return nil
}

// RemoveImage removes a image from hyperd
func (c *HyperClient) RemoveImage(image string) error {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	_, err := c.client.ImageRemove(ctx, &types.ImageRemoveRequest{Image: image})
	return err
}

// Version gets hyperd version
func (c *HyperClient) Version() (string, string, error) {
	ctx, cancel := context.WithTimeout(c.ctx, c.timeout)
	defer cancel()

	resp, err := c.client.Version(ctx, &types.VersionRequest{})
	if err != nil {
		return "", "", err
	}

	return resp.Version, resp.ApiVersion, nil
}
