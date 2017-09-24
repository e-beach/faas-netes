// Copyright (c) Alex Ellis 2017. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"os"

	"github.com/openfaas/faas-netes/handlers"
	"github.com/openfaas/faas-provider"
	bootTypes "github.com/openfaas/faas-provider/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

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

	readConfig := ReadConfig{}
	osEnv := OsEnv{}
	cfg := readConfig.Read(osEnv)

	var port int
	port = 8080
	bootstrapConfig := bootTypes.FaaSConfig{
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		TCPPort:      &port,
	}

	bootstrap.Serve(&bootstrapHandlers, &bootstrapConfig)
}
