// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"os"
	"time"
	"strconv"

	"github.com/openfaas/faas-netes/handlers"
	"github.com/openfaas/faas-provider"
	bootTypes "github.com/openfaas/faas-provider/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func parseIntValue(val string, fallback int) int {
	if len(val) > 0 {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal >= 0 {
			return parsedVal
		}
	}
	return fallback
}


func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	functionNamespace := "default"

	if namespace, exists := os.LookupEnv("function_namespace"); exists {
		functionNamespace = namespace
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	bootstrapHandlers := bootTypes.FaaSHandlers{
		FunctionProxy:  handlers.MakeProxy(functionNamespace),
		DeleteHandler:  handlers.MakeDeleteHandler(functionNamespace, clientset),
		DeployHandler:  handlers.MakeDeployHandler(functionNamespace, clientset),
		FunctionReader: handlers.MakeFunctionReader(functionNamespace, clientset),
		ReplicaReader:  handlers.MakeReplicaReader(functionNamespace, clientset),
		ReplicaUpdater: handlers.MakeReplicaUpdater(functionNamespace, clientset),
		UpdateHandler:  handlers.MakeUpdateHandler(functionNamespace, clientset),
	}

	readTimeout := parseIntValue(os.Getenv("read_timeout"), 8)
	writeTimeout := parseIntValue(os.Getenv("write_timeout"), 8)

	var port int
	port = 8080
	bootstrapConfig := bootTypes.FaaSConfig{
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		TCPPort:      &port,
	}

	bootstrap.Serve(&bootstrapHandlers, &bootstrapConfig)
}
