package util

import (
	"fmt"
	"net/http"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

func ListPods() {
	        // Build a rest.Config from configuration injected into the Pod by
        // Kubernetes. Clients will use the Pod's ServiceAccount principal.
        restconfig, err := rest.InClusterConfig()
        if err != nil {
                panic(err)
        }

        // If you need to know the Pod's Namespace, adjust the Pod's spec to pass
        // the information into an environment variable in advance via the downward
        // API.
        namespace := os.Getenv("NAMESPACE")
        if namespace == "" {
                panic("NAMESPACE was not set")
        }

        // Create a Kubernetes core/v1 client.
        coreclient, err := corev1client.NewForConfig(restconfig)
        if err != nil {
                panic(err)
        }

        mux := http.NewServeMux()
        mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
                rw.Header().Set("Cache-Control", "no-store, must-revalidate")
                rw.Header().Set("Content-Type", "text/plain")

                // List all Pods in our current Namespace.
                pods, err := coreclient.Pods(namespace).List(metav1.ListOptions{})
                if err != nil {
                        panic(err)
                }

                fmt.Fprintf(rw, "Pods in namespace %s:\n", namespace)
                for _, pod := range pods.Items {
                        fmt.Fprintf(rw, "  %s\n", pod.Name)
                }
        })

        // Run an HTTP server on port 8080 which will serve the pod and build list.
        err = http.ListenAndServe(":8080", mux)
        if err != nil {
                panic(err)
        }
}