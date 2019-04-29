package k8s

import (
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/rest"
	"os"
	"log"
	"errors"
)

func getKubeConfig(accessWithinCluster string) (error,*rest.Config) {
	var kubeConfig *rest.Config
	var err error
	if accessWithinCluster == "true" {
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal(err.Error())
		}		
	} else if accessWithinCluster == "false" {
		var kubeconfigFile string
		if home := homeDir(); home != "" {
			kubeconfigFile = home + "/.kube" + "/config"
			if _, err := os.Stat(kubeconfigFile); os.IsNotExist(err) {
				log.Println(err,"\n" + "ERROR: Kube config not found under " + kubeconfigFile)
				return err,nil
			}
		} else {
			log.Println("ERROR: Could not find homedir, check environment!")
			return err,nil
		}

		// use the current context in kubeconfig
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfigFile)
		if err != nil {
			log.Println(err.Error())
			return err,nil
		}		
	} else {
		log.Println("ERROR: Parameter AccessWithinCluster not set to true or false")
		err := errors.New("Parameter AccessWithinCluster not set to true or false")
		return err,nil
	}

	return nil,kubeConfig
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}