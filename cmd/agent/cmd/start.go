// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yunify/qingcloud-cni/pkg/nicmanagr"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"github.com/yunify/qingcloud-cni/pkg/qingactor"
	actorlog "github.com/AsynkronIT/protoactor-go/log"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start QingCloud container networking agent",
	Long: `QingCloud container networking agent is a daemon process which allocates and recollects nics resources.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLoglevel()
		runtime.GOMAXPROCS(runtime.NumCPU())
		runtime.GC()

		remote.Start("0.0.0.0:31080", remote.WithEndpointWriterBatchSize(10000))


		msg, err := qingactor.NewQingcloudInitializeMessage(viper.GetString("QYAccessFilePath"), viper.GetString("zone"))
		if err != nil {
			log.Errorf("Invalid QingCloud configuration: %v", err)
			return
		}
		pid := actor.NewLocalPID(qingactor.QingCloudActorName)
		pid.Tell(*msg)

		poolInitMsg, err := nicmanagr.NewResourcePoolInitMessage(viper.GetStringSlice("vxnet"),viper.GetString("policy"))

		pid = actor.NewLocalPID(nicmanagr.NicManagerActorName)
		pid.Tell(*poolInitMsg)

		gatewaymsg:=nicmanagr.InitGatewayMessage{Nsname:viper.GetString("gatewayns")}
		pid = actor.NewLocalPID(nicmanagr.GatewayManagerActorName)
		pid.Tell(gatewaymsg)

		//event loop
		systemCh := make(chan os.Signal, 4)
		signal.Notify(systemCh, os.Interrupt, syscall.SIGTERM)
		for {
			select {
			case <-systemCh:
				log.Infof("Got interrupt event, shutdown agent..")
				return
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().String("QYAccessFilePath", "/etc/qingcloud/client.yaml", "QingCloud Access File Path")
	startCmd.Flags().String("zone", "pek3a", "QingCloud zone")
	startCmd.Flags().StringSlice("vxnet", []string{"vxnet-xxxxxxx"}, "vxnet id list")
	startCmd.Flags().String("iface", "eth0", "Default nic which is used by host and will not be deleted")
	startCmd.Flags().String("policy","FailRotate", "policy of Selecting which vxnet to create nic from.(FailRotate,RoundRotate,Random)")
	startCmd.Flags().String("gatewayns","hostnic","gateway nic name")
	viper.BindPFlags(startCmd.Flags())
	log.SetLevel(log.DebugLevel)
	actor.SetLogLevel(actorlog.DebugLevel)
}

