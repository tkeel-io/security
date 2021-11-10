/*
Copyright 2021 The tKeel Authors.
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
package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/tkeel-io/security/pkg/apiserver/options"
	"github.com/tkeel-io/security/pkg/version"

	"github.com/spf13/cobra"
)

func NewAuthServerCmd() *cobra.Command {
	opts := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "auth: an http server of tkeel plugin of auth.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(SetupSignalHandler(), opts)
		},
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "print the version of tkeel auth server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(version.Get())
		},
	}

	cmd.AddCommand(versionCmd)

	return cmd
}

func Run(ctx context.Context, opt *options.ServerRunOptions) (err error) {
	apiServer, err := opt.NewAPIServer(ctx.Done())
	if err != nil {
		return
	}
	err = apiServer.PrepareRun(ctx.Done())
	if err != nil {
		return
	}
	err = apiServer.Run(ctx)
	return
}

func SetupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}
