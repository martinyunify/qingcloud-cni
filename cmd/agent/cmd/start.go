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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/yunify/qingcloud-cni/pkg/nicmanagr"
	"runtime"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/AsynkronIT/protoactor-go/mailbox"
	log "github.com/sirupsen/logrus"
	"github.com/yunify/qingcloud-cni/pkg/messages"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start QingCloud container networking agent",
	Long: `QingCloud container networking agent is a daemon process which allocates and recollects nics resources.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		runtime.GOMAXPROCS(runtime.NumCPU())
		runtime.GC()
		remote.Start("0.0.0.0:31080", remote.WithEndpointWriterBatchSize(10000))
		props := actor.FromProducer(nicmanagr.NewNicManager).WithMailbox(mailbox.Bounded(1000))
		pid, err := actor.SpawnNamed(props, "nicmanager")
		if err != nil {
			log.Error(err)
		}
		pid.Tell(messages.QingcloudInitializeMessage{
			QingAccessFile: viper.GetString("QYAccessFilePath"),
			Zone:           viper.GetString("Zone"),
		})

		for
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
	startCmd.Flags().String("Zone", "pek3a", "QingCloud zone")
	startCmd.Flags().StringSlice("vxnet", []string{"vxnet-xxxxxxx"}, "vxnet id list")

	viper.BindPFlags(startCmd.Flags())
}
