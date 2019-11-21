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
	"bytes"
	"fmt"
	"fossul/src/engine/util"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"regexp"
	"strings"
)

func ExecuteCommand(podName, containerName, namespace, accessWithinCluster string, args ...string) util.Result {
	baseCmd := args[0]
	cmdArgs := args[1:]

	var result util.Result
	var messages []util.Message

	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		message := util.SetMessage("ERROR", err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result
	}

	s0 := fmt.Sprintf("Executing command [%s %s] in namespace [%s] on pod [%s] container [%s]", baseCmd, strings.Join(cmdArgs, " "), namespace, podName, containerName)
	message := util.SetMessage("CMD", s0)
	messages = append(messages, message)

	var (
		execOut bytes.Buffer
		execErr bytes.Buffer
	)

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't create kube config: "+err.Error())
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
		Stdin:     false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(kubeConfig, "POST", req.URL())
	if err != nil {
		message := util.SetMessage("ERROR", "Failed to init executor: "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &execOut,
		Stderr: &execErr,
		Tty:    false,
	})

	message = util.SetMessage("DEBUG", "STDOUT: "+execOut.String())
	messages = append(messages, message)

	if err != nil {
		message := util.SetMessage("ERROR", "Could not execute command: "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result
	}

	if execErr.Len() > 0 {
		message := util.SetMessage("WARN", "STDERR: "+execErr.String())
		messages = append(messages, message)
	}

	s1 := fmt.Sprintf("Command [%s %s] on pod [%s] container [%s] completed successfully", baseCmd, strings.Join(cmdArgs, " "), podName, containerName)
	message = util.SetMessage("INFO", s1)
	messages = append(messages, message)

	result = util.SetResult(0, messages)
	return result
}

func ExecuteCommandWithStdout(podName, containerName, namespace, accessWithinCluster string, args ...string) (util.Result, string) {
	baseCmd := args[0]
	cmdArgs := args[1:]

	var result util.Result
	var messages []util.Message

	err, kubeConfig := getKubeConfig(accessWithinCluster)
	if err != nil {
		message := util.SetMessage("ERROR", err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result, ""
	}

	s0 := fmt.Sprintf("Executing command [%s %s] in namespace [%s] on pod [%s] container [%s]", baseCmd, strings.Join(cmdArgs, " "), namespace, podName, containerName)
	message := util.SetMessage("CMD", s0)
	messages = append(messages, message)

	var (
		execOut bytes.Buffer
		execErr bytes.Buffer
	)

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		message := util.SetMessage("ERROR", "Couldn't create kube config: "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result, ""
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
		Stdin:     false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(kubeConfig, "POST", req.URL())
	if err != nil {
		message := util.SetMessage("ERROR", "Failed to init executor: "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result, ""
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &execOut,
		Stderr: &execErr,
		Tty:    false,
	})

	message = util.SetMessage("DEBUG", "Command stdout: "+execOut.String())
	messages = append(messages, message)

	if err != nil {
		message := util.SetMessage("ERROR", "Could not execute command: "+err.Error())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result, ""
	}

	if execErr.Len() > 0 {
		message := util.SetMessage("ERROR", "Stderr: "+execErr.String())
		messages = append(messages, message)

		result = util.SetResult(1, messages)
		return result, ""
	}

	s1 := fmt.Sprintf("Command [%s %s] on pod [%s] container [%s] completed successfully", baseCmd, strings.Join(cmdArgs, " "), podName, containerName)
	message = util.SetMessage("INFO", s1)
	messages = append(messages, message)

	result = util.SetResult(0, messages)
	return result, execOut.String()
}

func IsRemoteCommand(arg string) bool {
	re := regexp.MustCompile(`:\S+`)
	match := re.FindStringSubmatch(arg)
	if match != nil {
		return true
	} else {
		return false
	}
}
