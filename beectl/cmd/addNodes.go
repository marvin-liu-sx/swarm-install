/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"errors"

	"github.com/marvin9002/swarm-install/beectl/nodes"
	"github.com/spf13/cobra"
)

// addNodesCmd represents the addNodes command
var addNodesCmd = &cobra.Command{
	Use:   "add-nodes",
	Short: "添加bee节点",
	RunE: func(cmd *cobra.Command, args []string) error {
		count, err := cmd.Flags().GetInt("count")
		if err != nil {
			return err
		}
		dir, err := cmd.Flags().GetString("data-dir")
		if err != nil {
			return err
		}

		password, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}

		ports, err := cmd.Flags().GetString("ports")
		if err != nil {
			return err
		}
		swapEndpoint, err := cmd.Flags().GetString("swap-endpoint")
		if err != nil {
			return err
		}

		return install(count, dir, password, ports, swapEndpoint)
	},
}

func init() {
	rootCmd.AddCommand(addNodesCmd)
	addNodesCmd.Flags().IntP("count", "", 1, "要添加几个节点 ")
	addNodesCmd.Flags().StringP("data-dir", "", "/root/", "data-dir (required) bee数据存放目录")
	addNodesCmd.Flags().StringP("password", "", "", "[可选] bee启动密码,不输入将随机生成")
	addNodesCmd.Flags().StringP("ports", "", "", "[可选] 必须用逗号分隔输入如: 1633,1634,1635 分别对应 api-addr/p2p-addr/debug-api-addr的端口(如输入只会创建一个节点) 注意不要端口重复,默认自动生成端口")
	addNodesCmd.Flags().StringP("swap-endpoint", "", "", "swap-endpoint address (required)")
}

func install(count int, dir, pwd, ports, endpoint string) error {
	if count <= 0 {
		return errors.New("节点数错误")
	}
	if dir == "" {
		return errors.New("数据目录错误")
	}
	if endpoint == "" {
		return errors.New("endpoint 错误")
	}

	bee := nodes.New(count, dir, pwd, ports, endpoint)
	err := bee.Install()
	if err != nil {
		return err
	}

	err = bee.Start()
	if err != nil {
		return err
	}

	return nil
}
