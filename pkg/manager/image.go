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
	"strings"

	"github.com/golang/glog"
	"github.com/hyperhq/hyperd/types"
	"golang.org/x/net/context"
	"k8s.io/frakti/pkg/utils"
	kubeapi "k8s.io/kubernetes/pkg/kubelet/api/v1alpha1/runtime"
)

// ListImages lists existing images.
func (s *KubeHyperManager) ListImages(ctx context.Context, req *kubeapi.ListImagesRequest) (*kubeapi.ListImagesResponse, error) {
	glog.V(3).Infof("ListImages with request %s", req.String())

	filter := *req.Filter.Image.Image
	images, err := s.hyperClient.GetImageList()
	if err != nil {
		glog.Errorf("Get image list failed: %v", err)
		return nil, err
	}

	var results []*kubeapi.Image
	for _, img := range images {
		if img.Id == filter {
			continue
		}
		if utils.InList(filter, img.RepoTags) {
			continue
		}

		imageSize := uint64(img.VirtualSize)
		results = append(results, &kubeapi.Image{
			Id:          &img.Id,
			RepoTags:    img.RepoTags,
			RepoDigests: img.RepoDigests,
			Size_:       &imageSize,
		})
	}

	return &kubeapi.ListImagesResponse{
		Images: results,
	}, nil
}

// ImageStatus returns the status of the image.
func (s *KubeHyperManager) ImageStatus(ctx context.Context, req *kubeapi.ImageStatusRequest) (*kubeapi.ImageStatusResponse, error) {
	glog.V(3).Infof("ImageStatus with request %s", req.String())

	return nil, fmt.Errorf("Not implemented")
}

func getHyperAuthConfig(auth *kubeapi.AuthConfig) *types.AuthConfig {
	if auth == nil {
		return nil
	}

	config := &types.AuthConfig{}
	if auth.Username != nil {
		config.Username = *auth.Username
	}
	if auth.Password != nil {
		config.Password = *auth.Password
	}
	if auth.Auth != nil {
		config.Auth = *auth.Auth
	}
	if auth.RegistryToken != nil {
		config.Registrytoken = *auth.RegistryToken
	}
	if auth.ServerAddress != nil {
		config.Serveraddress = *auth.ServerAddress
	}

	return config
}

// Get a repos name and returns the right reposName + tag|digest
// The tag can be confusing because of a port in a repository name.
//     Ex: localhost.localdomain:5000/samalba/hipache:latest
//     Digest ex: localhost:5000/foo/bar@sha256:bc8813ea7b3603864987522f02a76101c17ad122e1c46d790efc0fca78ca7bfb
func parseRepositoryTag(repos string) (string, string) {
	n := strings.Index(repos, "@")
	if n >= 0 {
		parts := strings.Split(repos, "@")
		return parts[0], parts[1]
	}
	n = strings.LastIndex(repos, ":")
	if n < 0 {
		return repos, ""
	}
	if tag := repos[n+1:]; !strings.Contains(tag, "/") {
		return repos[:n], tag
	}
	return repos, ""
}

// PullImage pulls a image with authentication config.
func (s *KubeHyperManager) PullImage(ctx context.Context, req *kubeapi.PullImageRequest) (*kubeapi.PullImageResponse, error) {
	glog.V(3).Infof("PullImage with request %s", req.String())

	image := *req.Image.Image
	repo, tag := parseRepositoryTag(image)
	auth := getHyperAuthConfig(req.Auth)
	err := s.hyperClient.PullImage(repo, tag, auth, nil)
	if err != nil {
		glog.Errorf("Pull image %s failed: %v", image, err)
		return nil, err
	}

	return &kubeapi.PullImageResponse{}, nil
}

// RemoveImage removes the image.
func (s *KubeHyperManager) RemoveImage(ctx context.Context, req *kubeapi.RemoveImageRequest) (*kubeapi.RemoveImageResponse, error) {
	glog.V(3).Infof("RemoveImage with request %s", req.String())

	image := *req.Image.Image
	err := s.hyperClient.RemoveImage(image)
	if err != nil {
		glog.Errorf("Remove image %s failed: %v", image, err)
		return nil, err
	}

	return &kubeapi.RemoveImageResponse{}, nil
}
