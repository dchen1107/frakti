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

package manager

import (
	"fmt"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	kubeapi "k8s.io/kubernetes/pkg/kubelet/api/v1alpha1/runtime"
)

// CreateContainer creates a new container in specified PodSandbox
func (s *KubeHyperManager) CreateContainer(ctx context.Context, req *kubeapi.CreateContainerRequest) (*kubeapi.CreateContainerResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// StartContainer starts the container.
func (s *KubeHyperManager) StartContainer(ctx context.Context, req *kubeapi.StartContainerRequest) (*kubeapi.StartContainerResponse, error) {
	glog.V(3).Infof("StartContainer with request %s", req.String())

	err := s.hyperClient.StartContainer(*req.ContainerId)
	if err != nil {
		glog.Errorf("Start container %s failed: %v", *req.ContainerId, err)
		return nil, err
	}

	return &kubeapi.StartContainerResponse{}, nil
}

// StopContainer stops a running container with a grace period (i.e., timeout).
func (s *KubeHyperManager) StopContainer(ctx context.Context, req *kubeapi.StopContainerRequest) (*kubeapi.StopContainerResponse, error) {
	glog.V(3).Infof("StopContainer with request %s", req.String())

	err := s.hyperClient.StopContainer(*req.ContainerId, *req.Timeout)
	if err != nil {
		glog.Errorf("Stop container %s failed: %v", *req.ContainerId, err)
		return nil, err
	}

	return &kubeapi.StopContainerResponse{}, nil
}

// RemoveContainer removes the container.
func (s *KubeHyperManager) RemoveContainer(ctx context.Context, req *kubeapi.RemoveContainerRequest) (*kubeapi.RemoveContainerResponse, error) {
	glog.V(3).Infof("RemoveContainer with request %s", req.String())

	err := s.hyperClient.RemoveContainer(*req.ContainerId)
	if err != nil {
		glog.Errorf("Remove container %s failed: %v", *req.ContainerId, err)
		return nil, err
	}

	return &kubeapi.RemoveContainerResponse{}, nil
}

// ListContainers lists all containers by filters.
func (s *KubeHyperManager) ListContainers(ctx context.Context, req *kubeapi.ListContainersRequest) (*kubeapi.ListContainersResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// ContainerStatus returns the container status.
func (s *KubeHyperManager) ContainerStatus(ctx context.Context, req *kubeapi.ContainerStatusRequest) (*kubeapi.ContainerStatusResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// Exec execute a command in the container.
func (s *KubeHyperManager) Exec(stream kubeapi.RuntimeService_ExecServer) error {
	return fmt.Errorf("Not implemented")
}
