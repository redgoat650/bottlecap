/*Package cmd defines command cli

Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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

	"github.com/spf13/cobra"
)

var (
	srcImgPath string
	resX       int
	resY       int
)

// mosaicCmd represents the mosaic command
var mosaicCmd = &cobra.Command{
	Use:   "mosaic",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mosaic called", args)
	},
}

func init() {
	rootCmd.AddCommand(mosaicCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mosaicCmd.PersistentFlags().String("foo", "", "A help for foo")
	mosaicCmd.PersistentFlags().StringVarP(&srcImgPath, "img", "i", "", "Image to be tiled")
	err := mosaicCmd.MarkPersistentFlagRequired("img")
	if err != nil {
		fmt.Println("sdf", err)
		os.Exit(1)
	}
	err = mosaicCmd.MarkPersistentFlagFilename("img")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mosaicCmd.PersistentFlags().IntVar(&resX, "resX", 100, "Resolution of the resulting mosaic X-axis")
	mosaicCmd.PersistentFlags().IntVar(&resY, "resY", 0, "Resolution of the resulting mosaic Y-axis")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mosaicCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
