// Copyright (c) 2019-Present Pivotal Software, Inc. All Rights Reserved.
//
// This program and the accompanying materials are made available under the terms of the under the Apache License,
// Version 2.0 (the "License”); you may not use this file except in compliance with the License. You may
// obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the
// License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"os"

	"github.com/cf-platform-eng/uaa-client-crud/pkg/cmd"
	"github.com/spf13/cobra"
)

func newRootCmd(args []string) *cobra.Command {
	root := &cobra.Command{
		Use:     "uaaclient",
		Short:   "uaa-client-crud",
		Version: "0.1.0",
	}

	flags := root.PersistentFlags()
	out := root.OutOrStdout()
	root.AddCommand(
		cmd.NewCreateClientCmd(out),
		cmd.NewDeleteClientCmd(out),
	)

	flags.Parse(args)

	return root
}

func main() {
	rootCmd := newRootCmd(os.Args[1:])
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
