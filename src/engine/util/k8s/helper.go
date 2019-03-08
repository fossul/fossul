package k8s

import (
	"flag"
	"path/filepath"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
)

func getKubeConfig(accessWithinCluster string) *rest.Config {
	var kubeConfig *rest.Config
	var err error
	if accessWithinCluster == "true" {
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}		
	} else {
		var kubeconfigFile *string
		if home := homeDir(); home != "" {
			kubeconfigFile = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		} else {
			kubeconfigFile = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
		}
		flag.Parse()
	
		// use the current context in kubeconfig
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", *kubeconfigFile)
		if err != nil {
			panic(err.Error())
		}		
	}

	return kubeConfig
}