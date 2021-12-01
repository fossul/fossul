/*
Copyright 2019 The Fossul Authors.
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
package k8s

import (
	"context"

	"github.com/fossul/fossul/src/client/k8s/virtctrl/client/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func listVirtualMachines() {
	cfg, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err)
	}

	client := versioned.NewForConfigOrDie(cfg)
	client.KubevirtV1().VirtualMachineInstances(v1.NamespaceAll).List(context.Background(), v1.ListOptions{})
}
