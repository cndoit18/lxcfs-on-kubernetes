/*
 *     Copyright 2021 Yinan Li <cndoit18@outlook.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"flag"
	"os"

	"github.com/spf13/pflag"

	lxcfsadmission "github.com/cndoit18/lxcfs-on-kubernetes/pkg/admission"
	"github.com/cndoit18/lxcfs-on-kubernetes/version"
	klog "k8s.io/klog/v2"
	klogr "k8s.io/klog/v2/klogr"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	logr "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	k8sadmission "sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var log = logr.Log.WithName("main")

func init() {
	logr.SetLogger(klogr.New())
}

// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=*
// +kubebuilder:rbac:groups="",resources=configmaps;events,verbs=*

func main() {
	var (
		leaderElection          = flag.Bool("leader-election", false, "LeaderElection determines whether or not to use leader election when starting the manager.")
		leaderElectionNamespace = flag.String("leader-election-namespace", "default", "The leader election namespace.")
		leaderElectionID        = flag.String("leader-election-id", "lxcfs-on-kubernetes-leader-election", "The leader election id.")
		lxcfsPath               = flag.String("lxcfs-path", "/var/lib/lxc/lxcfs", "Path for lxcfs mounts.")
	)

	// set logging
	fs := pflag.NewFlagSet(os.Args[0], pflag.ExitOnError)
	fs.AddGoFlagSet(flag.CommandLine)

	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	fs.AddGoFlagSet(klogFlags)

	if err := fs.Parse(os.Args); err != nil {
		log.Error(err, "Failed to parse command line args")
		os.Exit(1)
	}

	log.Info("manager Starting", "Version", version.Version())

	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Failed to get configuration.")
		os.Exit(1)
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		LeaderElection:          *leaderElection,
		LeaderElectionNamespace: *leaderElectionNamespace,
		LeaderElectionID:        *leaderElectionID,
		HealthProbeBindAddress:  ":8081",
	})
	if err != nil {
		log.Error(err, "Failed to create manager.")
		os.Exit(1)
	}

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		log.Error(err, "Failed to add healthzCheck to manager")
		os.Exit(1)
	}

	if err := lxcfsadmission.AddToManager(mgr,
		lxcfsadmission.WithMutatePath(*lxcfsPath),
		lxcfsadmission.WithMutateDecoder(k8sadmission.NewDecoder(mgr.GetScheme())),
	); err != nil {
		log.Error(err, "Failed to add admission to manager")
		os.Exit(1)
	}

	// Start the Cmd
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Error(err, "unable to start the manager")
		os.Exit(1)
	}
}
