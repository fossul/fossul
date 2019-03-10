package k8s

import (
//	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"fmt"
	"strings"
)

func GetPod(namespace, serviceName, accessWithinCluster string) string {
	var kubeConfig *rest.Config = getKubeConfig(accessWithinCluster)
	
	// create the clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Println(err.Error())
	}

    pods, err := clientset.CoreV1().Pods(namespace).List(metav1.ListOptions{})
    if err != nil {
    	panic(err)
    }

	fmt.Println("DEBUG: Pods in namespace", namespace)
	var ourPod string
    for _, pod := range pods.Items {
		fmt.Println("DEBUG: Pod", pod.Name, pod.Status.Phase) 
		if strings.Contains(pod.Name,serviceName) && pod.Status.Phase == "Running" {
			fmt.Println("INFO: Running Pod Found:", pod.Name)
			ourPod = pod.Name
		}
	}

	/*            
	pod := "fossil-app-2-zpdgr"
	_, err = clientset.CoreV1().Pods(namespace).Get(pod, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			pod, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	}
	*/

	return ourPod
}