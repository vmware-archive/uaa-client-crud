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
	"code.cloudfoundry.org/lager"
	uaacmd "github.com/cf-platform-eng/uaa-client-crud/pkg/cmd"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	cmd := newRootCmd(os.Args[1:])
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newRootCmd(args []string) *cobra.Command {
	command := &cobra.Command{
		Use:   "uaaclient",
		Short: "uaa-client-crud",
	}

	out := command.OutOrStdout()
	logger := lager.NewLogger("test")
	logger.RegisterSink(lager.NewWriterSink(out, lager.INFO))
	flags := command.PersistentFlags()
	command.AddCommand(
		uaacmd.NewCreateClientCmd(logger),
		uaacmd.NewDeleteClientCmd(logger),
	)

	flags.Parse(args)

	return command
}
