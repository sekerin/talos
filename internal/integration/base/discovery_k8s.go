// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// +build integration_k8s

package base

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/talos-systems/talos/cmd/osctl/pkg/client"
)

func discoverNodesK8s(client *client.Client) ([]string, error) {
	kubeconfig, err := client.Kubeconfig(context.Background())
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", func() (*clientcmdapi.Config, error) {
		return clientcmd.Load(kubeconfig)
	})
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []string

	for _, node := range nodes.Items {
		for _, nodeAddress := range node.Status.Addresses {
			if nodeAddress.Type == v1.NodeInternalIP {
				result = append(result, nodeAddress.Address)
			}
		}
	}

	return result, nil
}
