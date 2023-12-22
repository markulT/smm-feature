package main

import (
	"smm-features/pkg/repository"
	"smm-features/utils"
	"smm-features/utils/analytics"
)

func init() {
	utils.LoadEnvVariables()
	utils.ConnectToDb()
}

func main() {
	analyticsRepo := repository.NewAnalyticsRepo()
	channelRepo := repository.NewChannelRepo()

	analyticsTask := &analytics.AnalyticsTask{}
	bridgeRepo := analytics.CreateBridgeRepo(analyticsRepo, channelRepo)
	analyticsTask.SetRepo(bridgeRepo)
	go analyticsTask.RunAnalyticsModule()
	for {
	}
}
