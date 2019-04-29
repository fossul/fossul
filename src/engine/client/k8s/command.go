package k8s

import (
	"bytes"
	"fmt"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/remotecommand"
	"engine/util"
)

func ExecuteCommand(podName, containerName, namespace, accessWithinCluster string, args ...string) util.Result {
	baseCmd := args[0]
	cmdArgs := args[1:]

	var result util.Result
	var messages []util.Message

	err,kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		message := util.SetMessage("ERROR", err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result
	}

	s0 := fmt.Sprintf("Executing command [%s %s] on pod [%s] container [%s]",baseCmd, strings.Join(cmdArgs, " "),podName,containerName)
	message := util.SetMessage("CMD", s0)
	messages = append(messages, message)

	var (
		execOut bytes.Buffer
		execErr bytes.Buffer
	)

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't create kube config: " + err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")
	req.VersionedParams(&v1.PodExecOptions{
		Container: containerName,
		Command:   args,
		Stdout:    true,
		Stderr:    true,
		Stdin: false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(kubeConfig, "POST", req.URL())
	if err != nil {
		message := util.SetMessage("ERROR", "Failed to init executor: " + err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &execOut,
		Stderr: &execErr,
		Tty:    false,
	})

	message = util.SetMessage("DEBUG", "Command stdout: " + execOut.String())
	messages = append(messages, message)

	if err != nil {
		message := util.SetMessage("ERROR", "Could not execute command: " + err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result
	}

	if execErr.Len() > 0 {
		message := util.SetMessage("ERROR", "Stderr: " + execErr.String())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result
	}

	s1 := fmt.Sprintf("Command [%s %s] on pod [%s] container [%s] completed successfully",baseCmd, strings.Join(cmdArgs, " "),podName,containerName)
	message = util.SetMessage("INFO", s1)
	messages = append(messages, message)

	result = util.SetResult(0, messages)
	return result
}