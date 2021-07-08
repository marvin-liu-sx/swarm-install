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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marvin9002/swarm-install/beectl/nodes"
	"github.com/marvin9002/swarm-install/beectl/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// addConfigCmd represents the addConfig command
var addConfigCmd = &cobra.Command{
	Use:   "add-config",
	Short: "添加已经存在的bee节点配置到此监控",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cmd.Flags().GetString("bee-config")
		if err != nil {
			return err
		}
		bee := nodes.NewSwarm()
		bee = bee
		if utils.IsFile(cfg) {
			config := viper.New()
			config.SetConfigType("yaml")
			config.SetConfigFile(cfg)
			//尝试进行配置读取
			if err := config.ReadInConfig(); err != nil {
				panic(err)
			}

			//打印文件读取出来的内容:
			fmt.Println(config.Get("debug-api-addr"))
			fmt.Println(config.Get("config"))
			fmt.Println(config.Get("data-dir"))
			err := bee.UpdateCfg(fmt.Sprintf("bee%s", strings.Split(config.Get("debug-api-addr").(string), ":")[1]), strings.Split(config.Get("debug-api-addr").(string), ":")[1])
			if err != nil {
				return err
			}
		} else {
			var files []string

			err := filepath.Walk(cfg, func(path string, f os.FileInfo, err error) error {
				if f == nil {
					return err
				}
				if f.IsDir() {
					return nil
				}
				files = append(files, path)
				return nil
			})
			if err != nil {
				fmt.Printf("filepath.Walk() returned %v\n", err)
			}
			for _, file := range files {
				config := viper.New()
				config.SetConfigType("yaml")
				config.SetConfigFile(file)
				//尝试进行配置读取
				if err := config.ReadInConfig(); err != nil {
					panic(err)
				}

				//打印文件读取出来的内容:
				fmt.Println(config.Get("debug-api-addr"))
				fmt.Println(config.Get("config"))
				fmt.Println(config.Get("data-dir"))
				err := bee.UpdateCfg(fmt.Sprintf("bee%s", strings.Split(config.Get("debug-api-addr").(string), ":")[1]), strings.Split(config.Get("debug-api-addr").(string), ":")[1])
				if err != nil {
					return err
				}
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addConfigCmd)
	addConfigCmd.Flags().StringP("bee-config", "", "", "导入已经存在的bee配置文件|目录 如: /etc/bee 或者 /etc/bee/bee.yaml , 如果是目录将自动加载目录下的所有yaml配置.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addConfigCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addConfigCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
