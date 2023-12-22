package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"smm-features/pkg/models"
	"smm-features/utils"
	"time"
)

type AnalyticsRepo interface {
	GetManyByChannelIDAndDate(c context.Context, chID string, dateStart time.Time, dateEnd time.Time) ([]*models.AnalyticsUnit, error)
	SaveChannelAvgViews(c context.Context, channelName string, v float64) error
	GetAllOriginalChannels(c context.Context) ([]*models.Channel, error)
}

type defaultAnalyticsRepo struct{}

func NewAnalyticsRepo() AnalyticsRepo {
	return &defaultAnalyticsRepo{}
}

func (ar *defaultAnalyticsRepo) GetAllOriginalChannels(c context.Context) ([]*models.Channel, error) {

	channelsCollection := utils.DB.Collection("channels")

	//distinct, err := analyticsCollection.Distinct(context.Background(), "name", bson.M{})
	pipeline := bson.A{
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"title": "$title",
				},
				"doc": bson.M{
					"$first": "$$ROOT",
				},
			},
		},
		bson.M{
			"$replaceRoot": bson.M{
				"newRoot": "$doc",
			},
		},
	}

	channelsCursor, err := channelsCollection.Aggregate(c, pipeline)
	if err != nil {
		return nil, err
	}
	defer channelsCursor.Close(context.Background())

	var result []*models.Channel
	for channelsCursor.Next(c) {
		var channel models.Channel
		if err := channelsCursor.Decode(&channel); err != nil {
			return nil, err
		}
		result = append(result, &channel)
	}
	fmt.Println(result)
	return result, nil
}

func (ar *defaultAnalyticsRepo) SaveChannelAvgViews(c context.Context, channelName string, v float64) error {
	analyticsCollection := utils.DB.Collection("analytics")

	analyticsUnit := models.AnalyticsUnit{
		ID:          uuid.New(),
		ChannelName: channelName,
		Value:       v,
		Date:        time.Now(),
		Type:        "avgViews",
	}

	_, err := analyticsCollection.InsertOne(c, analyticsUnit)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (ar *defaultAnalyticsRepo) GetManyByChannelIDAndDate(c context.Context, channelName string, dateStart time.Time, dateEnd time.Time) ([]*models.AnalyticsUnit, error) {

	var analyticsUnitArray []*models.AnalyticsUnit
	analyticsCollection := utils.DB.Collection("analytics")
	curs, err := analyticsCollection.Find(c, bson.M{"channelName": channelName, "date": bson.M{
		"$gte": time.Date(dateStart.Year(), dateStart.Month(), dateStart.Day(), dateStart.Hour(), 0, 0, 0, time.UTC),
		"$lt":  dateEnd,
	}})
	defer curs.Close(c)
	if err != nil {
		return nil, err
	}
	if curs.Err() != nil {
		return nil, curs.Err()
	}

	for curs.Next(c) {
		var analyticsUnit models.AnalyticsUnit
		if err := curs.Decode(&analyticsUnit); err != nil {
			return nil, err
		}
		analyticsUnitArray = append(analyticsUnitArray, &analyticsUnit)
	}

	return analyticsUnitArray, nil
}
