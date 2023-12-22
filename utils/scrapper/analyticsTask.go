package scrapper

import (
	"context"
	"fmt"
	"time"
)

type Channel struct {
	Name        string
	ChannelHash string
}

type ScrapperTask struct {
	//AuthorizeBot(c context.Context, respch chan struct{})
	//CollectAvgViews(channelName string) (string, error)
	//CollectAvgViewsByChannelID(channelID string, chn string) error
}

func RunAnalyticsTask(chRepo ChannelAuthorizerRepository, s *Scrapper) {
	time.Sleep(10 * time.Second)
	RunAnalytics(chRepo, *s)
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:
			RunAnalytics(chRepo, *s)
		}
	}
}

func RunAnalytics(chRepo ChannelAuthorizerRepository, s Scrapper) {
	var channelList []*Channel
	channelList = chRepo.GetAllChannels()
	initTime := time.Now()
	for _, channel := range channelList {
		respch := make(chan struct{}, 1)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		processingCompleted := false

		if channel.ChannelHash != "" {
			go func(ch chan struct{}) {
				s.CollectAvgViewsByChannelID(channel.ChannelHash, channel.Name)
				respch <- struct{}{}
			}(respch)
		} else {
			go func(ch chan struct{}) {
				hash, _ := s.CollectAvgViews(channel.Name)
				chRepo.SaveChannelHash(channel.Name, hash)
			}(respch)
		}

		select {
		case <-respch:
			// Processing completed within the time limit
			processingCompleted = true
		case <-ctx.Done():
			// Processing exceeded the time limit
		}

		cancel()

		// Skip to the next iteration if processing took more than 5 seconds
		if !processingCompleted {
			continue
		}
	}
	elapsed := time.Since(initTime)
	fmt.Printf("The whole counting process took %s to run.\n", elapsed)
}
