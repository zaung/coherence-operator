/*
 * Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
 * Licensed under the Universal Permissive License v 1.0 as shown at
 * http://oss.oracle.com/licenses/upl.
 */

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-logr/logr"
	hoflags "github.com/operator-framework/operator-sdk/pkg/helm/flags"
	"github.com/operator-framework/operator-sdk/pkg/helm/release"
	"github.com/oracle/coherence-operator/pkg/flags"
	"github.com/oracle/coherence-operator/pkg/operator"
	cohrest "github.com/oracle/coherence-operator/pkg/rest"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"

	"github.com/oracle/coherence-operator/pkg/apis"
	"github.com/oracle/coherence-operator/pkg/controller"

	helmctl "github.com/operator-framework/operator-sdk/pkg/helm/controller"
	"github.com/operator-framework/operator-sdk/pkg/helm/watches"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	kubemetrics "github.com/operator-framework/operator-sdk/pkg/kube-metrics"
	"github.com/operator-framework/operator-sdk/pkg/leader"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/operator-framework/operator-sdk/pkg/restmapper"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	"github.com/spf13/pflag"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// NOTE: This file was generated by the Operator-SDK and then edited to add extra functionality for the Coherence Operator.
// Upgrading from one version of the Operator-SDK to another may require this file to be re-genrated (or rather the code
// in this file will need to be edited. To make this easier comment blocks have been added to make it obvious where tius
// file was changed from the original generated file.

// >>>>>>>> Coherence Operator code added to Operator SDK the generated file ---------------------------

// BuildInfo is a pipe delimited string of build information injected by the Go linker at build time.
var BuildInfo string

// <<<<<<<< Coherence Operator code added to Operator SDK the generated file ---------------------------

// Change below variables to serve metrics on different host or port.
var (
	metricsHost               = "0.0.0.0"
	metricsPort         int32 = 8383
	operatorMetricsPort int32 = 8686
)

var log = logf.Log.WithName("cmd")

func printVersion() {
	log.Info(fmt.Sprintf("Go Version: %s", runtime.Version()))
	log.Info(fmt.Sprintf("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH))
	log.Info(fmt.Sprintf("Version of operator-sdk: %v", sdkVersion.Version))
}

func main() {
	// Add the zap logger flag set to the CLI. The flag set must
	// be added before calling pflag.Parse().
	pflag.CommandLine.AddFlagSet(zap.FlagSet())

	// Add flags registered by imported packages (e.g. glog and
	// controller-runtime)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	// >>>>>>>> Coherence Operator code added to Operator SDK the generated file ---------------------------
	// create Coherence Operator flags
	cohf := flags.AddTo(pflag.CommandLine)

	// create Helm Operator flags
	hflags := hoflags.AddTo(pflag.CommandLine)
	fmt.Println(hflags)
	// <<<<<<<< Coherence Operator code added to Operator SDK the generated file ---------------------------

	pflag.Parse()

	// Use a zap logr.Logger implementation. If none of the zap
	// flags are configured (or if the zap flag set is not being
	// used), this defaults to a production zap logger.
	//
	// The logger instantiated here can be changed to any logger
	// implementing the logr.Logger interface. This logger will
	// be propagated through the whole operator, generating
	// uniform and structured logs.
	logf.SetLogger(zap.Logger())

	// >>>>>>>> Coherence Operator code added to Operator SDK the generated file ---------------------------
	printBuildInfo(log)
	// <<<<<<<< Coherence Operator code added to Operator SDK the generated file ---------------------------
	printVersion()

	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		log.Error(err, "failed to get watch namespace")
		os.Exit(1)
	}

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	ctx := context.TODO()
	// Become the leader before proceeding
	err = leader.Become(ctx, "coherence-operator-lock")
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		Namespace:          namespace,
		MapperProvider:     restmapper.NewDynamicRESTMapper,
		MetricsBindAddress: fmt.Sprintf("%s:%d", metricsHost, metricsPort),
	})
	if err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// >>>>>>>> Coherence Operator code added to Operator SDK the generated file ---------------------------

	// we must start the Operator ReST endpoint before any controllers start
	restServer, err := cohrest.StartRestServer(mgr, cohf)
	if err != nil {
		log.Error(err, "Error starting ReST server")
		os.Exit(1)
	}

	// wait until we can hit the server to ensure that it is up
	log.Info("Waiting for rest server to start")
	for i := 0; i < 10; i++ {
		_, err = http.Get(fmt.Sprintf("http://localhost:%d/ready", restServer.GetPort()))
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// ensure that the CRDs exist
	if err := operator.EnsureCRDs(mgr, cohf, log); err != nil {
		log.Error(err, "Error ensuring that CRDs exist")
		os.Exit(1)
	}

	// ensure that the configuration secret exists
	if err := operator.EnsureOperatorSecret(namespace, mgr, restServer, cohf, log); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// <<<<<<<< Coherence Operator code added to Operator SDK the generated file ---------------------------

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	// >>>>>>>> Coherence Operator code added to Operator SDK the generated file ---------------------------

	// Configure the Helm operator
	if err := setupHelm(mgr, namespace, hflags); err != nil {
		log.Error(err, "Manager exited non-zero")
		os.Exit(1)
	}

	// <<<<<<<< Coherence Operator code added to Operator SDK the generated file ---------------------------

	if err = serveCRMetrics(cfg); err != nil {
		log.Info("Could not generate and serve custom resource metrics", "error", err.Error())
	}

	// >>>>>>>> Coherence Operator code added to Operator SDK the generated file ---------------------------
	// We do not want to add a service here so the gernated code is commented out.

	//// Add to the below struct any other metrics ports you want to expose.
	//servicePorts := []v1.ServicePort{
	//	{Port: metricsPort, Name: metrics.OperatorPortName, Protocol: v1.ProtocolTCP, TargetPort: intstr.IntOrString{TypeIs: intstr.Int, IntVal: metricsPort}},
	//	{Port: operatorMetricsPort, Name: metrics.CRPortName, Protocol: v1.ProtocolTCP, TargetPort: intstr.IntOrString{TypeIs: intstr.Int, IntVal: operatorMetricsPort}},
	//}
	//// Create Service object to expose the metrics port(s).
	//_, err = metrics.CreateMetricsService(ctx, cfg, servicePorts)
	//if err != nil {
	//	log.Info(err.Error())
	//}
	// <<<<<<<< Coherence Operator code added to Operator SDK the generated file ---------------------------

	log.Info("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "Manager exited non-zero")
		os.Exit(1)
	}
}

// serveCRMetrics gets the Operator/CustomResource GVKs and generates metrics based on those types.
// It serves those metrics on "http://metricsHost:operatorMetricsPort".
func serveCRMetrics(cfg *rest.Config) error {
	// Below function returns filtered operator/CustomResource specific GVKs.
	// For more control override the below GVK list with your own custom logic.
	filteredGVK, err := k8sutil.GetGVKsFromAddToScheme(apis.AddToScheme)
	if err != nil {
		return err
	}
	// Get the namespace the operator is currently deployed in.
	operatorNs, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		return err
	}
	// To generate metrics in other namespaces, add the values below.
	ns := []string{operatorNs}
	// Generate and serve custom resource specific metrics.
	err = kubemetrics.GenerateAndServeCRMetrics(cfg, ns, filteredGVK, metricsHost, operatorMetricsPort)
	if err != nil {
		return err
	}
	return nil
}

// >>>>>>>> Coherence Operator code added to Operator SDK the generated file ---------------------------

func setupHelm(mgr manager.Manager, namespace string, hflags *hoflags.HelmOperatorFlags) error {
	// Setup Helm controller
	watchList, err := watches.Load(hflags.WatchesFile)
	if err != nil {
		log.Error(err, "failed to load Helm watches")
		return err
	}

	fmt.Println(watchList)
	for _, w := range watchList {
		fmt.Println(w)
		err := helmctl.Add(mgr, helmctl.WatchOptions{
			Namespace:               namespace,
			GVK:                     w.GroupVersionKind,
			ManagerFactory:          release.NewManagerFactory(mgr, w.ChartDir),
			ReconcilePeriod:         hflags.ReconcilePeriod,
			WatchDependentResources: w.WatchDependentResources,
		})
		if err != nil {
			log.Error(err, "failed to add Helm watche")
			return err
		}
	}

	return nil
}

// PrintBuildInfo prints the Coherence Operator build information to the log.
func printBuildInfo(log logr.Logger) {
	var (
		version string
		commit  string
		date    string
	)

	if BuildInfo != "" {
		parts := strings.Split(BuildInfo, "|")

		if len(parts) > 0 {
			version = parts[0]
		}

		if len(parts) > 1 {
			commit = parts[1]
		}

		if len(parts) > 2 {
			date = strings.Replace(parts[2], ".", " ", -1)
		}
	}

	log.Info(fmt.Sprintf("Coherence Operator Version: %s", version))
	log.Info(fmt.Sprintf("Coherence Operator Git commit: %s", commit))
	log.Info(fmt.Sprintf("Coherence Operator Build Time: %s", date))
}
