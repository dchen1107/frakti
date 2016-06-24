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
	"net"
	"time"

	"github.com/coreos/go-semver/semver"
	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"k8s.io/frakti/pkg/hyper"
	kubeapi "k8s.io/kubernetes/pkg/kubelet/api/v1alpha1/runtime"
)

const (
	hyperRuntimeName    = "hyper"
	minimumHyperVersion = "0.6.0"
	runtimeAPIVersion   = "0.1.0"

	hyperConnectionTimeout = 15 * time.Second
)

// KubeHyperManager serves the kubelet runtime gRPC api which will be
// consumed by kubelet
type KubeHyperManager struct {
	server      *grpc.Server
	hyperClient *hyper.HyperClient
}

// NewKubeHyperManager creates a new KubeHyperManager
func NewKubeHyperManager(hyperEndpoint string) (*KubeHyperManager, error) {
	hyperClient, err := hyper.NewHyperClient(hyperEndpoint, hyperConnectionTimeout)
	if err != nil {
		glog.Fatalf("Initialize hyper client failed: %v", err)
		return nil, err
	}

	version, _, err := hyperClient.Version()
	if err != nil {
		glog.Fatalf("Get hyperd version failed: %v", err)
		return nil, err
	}

	glog.V(3).Infof("Got hyperd version: %s", version)
	if check, err := checkVersion(version); !check {
		return nil, err
	}

	s := &KubeHyperManager{
		hyperClient: hyperClient,
		server:      grpc.NewServer(),
	}
	s.registerServer()

	return s, nil
}

// checkVersion checks whether hyperd's version is >=minimumHyperVersion
func checkVersion(version string) (bool, error) {
	hyperVersion, err := semver.NewVersion(version)
	if err != nil {
		glog.Errorf("Make semver failed: %v", version)
		return false, err
	}
	minVersion, err := semver.NewVersion(minimumHyperVersion)
	if err != nil {
		glog.Errorf("Make semver failed: %v", minimumHyperVersion)
		return false, err
	}
	if hyperVersion.LessThan(*minVersion) {
		return false, fmt.Errorf("Hyperd version is older than %s", minimumHyperVersion)
	}

	return true, nil
}

// Serve starts gRPC server at tcp://addr
func (s *KubeHyperManager) Serve(addr string) error {
	glog.V(1).Infof("Start frakti at %s", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		glog.Fatalf("Failed to listen %s: %v", addr, err)
		return err
	}

	return s.server.Serve(lis)
}

func (s *KubeHyperManager) registerServer() {
	kubeapi.RegisterRuntimeServiceServer(s.server, s)
	kubeapi.RegisterImageServiceServer(s.server, s)
}

// Version returns the runtime name, runtime version and runtime API version
func (s *KubeHyperManager) Version(ctx context.Context, req *kubeapi.VersionRequest) (*kubeapi.VersionResponse, error) {
	version, apiVersion, err := s.hyperClient.Version()
	if err != nil {
		glog.Errorf("Get hyper version failed: %v", err)
		return nil, err
	}

	runtimeName := hyperRuntimeName
	kubeletAPIVersion := runtimeAPIVersion
	return &kubeapi.VersionResponse{
		Version:           &kubeletAPIVersion,
		RuntimeName:       &runtimeName,
		RuntimeVersion:    &version,
		RuntimeApiVersion: &apiVersion,
	}, nil
}

// CreatePodSandbox creates a hyper Pod
func (s *KubeHyperManager) CreatePodSandbox(ctx context.Context, req *kubeapi.CreatePodSandboxRequest) (*kubeapi.CreatePodSandboxResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// DeletePodSandbox deletes the sandbox.
func (s *KubeHyperManager) DeletePodSandbox(ctx context.Context, req *kubeapi.DeletePodSandboxRequest) (*kubeapi.DeletePodSandboxResponse, error) {
	glog.V(3).Infof("DeletePodSandbox with request %s", req.String())

	err := s.hyperClient.RemovePod(*req.PodSandboxId)
	if err != nil {
		glog.Errorf("Remove pod %s failed: %v", *req.PodSandboxId, err)
		return nil, err
	}

	return &kubeapi.DeletePodSandboxResponse{}, nil
}

// PodSandboxStatus returns the Status of the PodSandbox.
func (s *KubeHyperManager) PodSandboxStatus(ctx context.Context, req *kubeapi.PodSandboxStatusRequest) (*kubeapi.PodSandboxStatusResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}

// ListPodSandbox returns a list of SandBox.
func (s *KubeHyperManager) ListPodSandbox(ctx context.Context, req *kubeapi.ListPodSandboxRequest) (*kubeapi.ListPodSandboxResponse, error) {
	return nil, fmt.Errorf("Not implemented")
}
