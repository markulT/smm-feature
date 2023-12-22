package analytics

import (
	"context"
	"fmt"
	"smm-features/pkg/repository"
	"smm-features/utils/scrapper"
)

type BridgeAnalyticsRepo struct {
	originalRepo repository.AnalyticsRepo
	channelRepo  repository.ChannelRepository
}

func CreateBridgeRepo(repo repository.AnalyticsRepo, chRepo repository.ChannelRepository) AnalyticsTaskRepo {
	return &BridgeAnalyticsRepo{originalRepo: repo, channelRepo: chRepo}
}

func (br *BridgeAnalyticsRepo) SaveChannelAvgViews(ch string, v float64) error {
	return br.originalRepo.SaveChannelAvgViews(context.Background(), ch, v)
}

func (br *BridgeAnalyticsRepo) GetAllChannels() []*scrapper.Channel {

	var channels []*scrapper.Channel
	originalChannels, err := br.originalRepo.GetAllOriginalChannels(context.Background())
	fmt.Println(originalChannels)
	if err != nil {
		return nil
	}

	for i, channel := range originalChannels {
		if i == 0 {
			ch := &scrapper.Channel{
				Name:        "@privatnamemarnya",
				ChannelHash: "",
			}
			channels = append(channels, ch)
		}
		ch := &scrapper.Channel{
			Name:        channel.Name,
			ChannelHash: channel.ChannelHash,
		}
		channels = append(channels, ch)

	}

	return channels
}

func (br *BridgeAnalyticsRepo) SaveChannelHash(ch string, hash string) error {
	return br.channelRepo.SaveChannelHash(context.Background(), ch, hash)
}
