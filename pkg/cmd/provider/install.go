// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package provider

import (
	"github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/internal/util/apiclient/server"
	"github.com/daytonaio/daytona/pkg/os"
	"github.com/daytonaio/daytona/pkg/provider/manager"
	"github.com/daytonaio/daytona/pkg/serverapiclient"
	"github.com/daytonaio/daytona/pkg/views/provider"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var providerInstallCmd = &cobra.Command{
	Use:     "install",
	Short:   "Install provider",
	Args:    cobra.NoArgs,
	Aliases: []string{"i"},
	Run: func(cmd *cobra.Command, args []string) {
		apiClient, err := server.GetApiClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		serverConfig, res, err := apiClient.ServerAPI.GetConfigExecute(serverapiclient.ApiGetConfigRequest{})
		if err != nil {
			log.Fatal(apiclient.HandleErrorResponse(res, err))
		}

		providersManifest, err := manager.GetProvidersManifest(*serverConfig.RegistryUrl)
		if err != nil {
			log.Fatal(err)
		}

		if providersManifest == nil {
			log.Fatal("Could not get providers manifest")
		}

		pluginList := convertToDTO(providersManifest)

		providerToInstall := provider.GetProviderFromPrompt(pluginList, "Choose a Provider to install")

		if providerToInstall == nil {
			return
		}

		downloadUrls := convertToStringMap((*providersManifest)[*providerToInstall.Name].Versions[*providerToInstall.Version].DownloadUrls)
		res, err = apiClient.ProviderAPI.InstallProviderExecute(serverapiclient.ApiInstallProviderRequest{}.Provider(serverapiclient.InstallProviderRequest{
			Name:         providerToInstall.Name,
			DownloadUrls: &downloadUrls,
		}))
		if err != nil {
			log.Fatal(apiclient.HandleErrorResponse(res, err))
		}

		if err != nil {
			log.Fatal(err)
		}
	},
}

func convertToDTO(manifest *manager.ProvidersManifest) []serverapiclient.Provider {
	pluginList := []serverapiclient.Provider{}
	for pluginName, pluginManifest := range *manifest {
		for version := range pluginManifest.Versions {
			pluginList = append(pluginList, serverapiclient.Provider{
				Name:    &pluginName,
				Version: &version,
			})
		}
	}

	return pluginList
}

func convertToStringMap(downloadUrls map[os.OperatingSystem]string) map[string]string {
	stringMap := map[string]string{}
	for os, url := range downloadUrls {
		stringMap[string(os)] = url
	}

	return stringMap
}
