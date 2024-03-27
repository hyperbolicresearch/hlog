package v1

import (
	"context"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

// GeneralObservables are the metrics that are generic for logs
// to be monitored. They are mostly about metadata that accompany
// the logs and are by default monitored by the system.
type GeneralObservables struct {
	ChannelsCount     int64
	LogsPerChannel    map[string]int64
	LogsPerSender     map[string]int64
	LogsPerLevel      map[string]int64
	SendersCount      int64
	LevelsCount       int64
	TotalIngestedLogs int64
	ThroughputPerTime []int64
}

// GeneralSystemObservables are metrics that are related to the
// infrastructure running the system.
type GeneralSystemObservables struct {
	DatabaseSizeOnDisk int64 // MongoDB:: https://pkg.go.dev/go.mongodb.org/mongo-driver@v1.14.0/mongo#DatabaseSpecification
}

// GetMongoDBGeneralObservables collects GeneralObservables metrics
// from MongoDB and systems built around it (if there are any).
func GetMongoDBGeneralObservables(cfg *config.MongoDB) *GeneralObservables {
	mongoClient := mongodb.Client(cfg.Server)
	db := mongoClient.Database(cfg.Database)

	// ChannelsCount is obtained by taking number of collections in the
	// database. As a reminder, we store each channel data in a separate
	// collection
	channelsNames, err := db.ListCollectionNames(
		context.TODO(),
		bson.D{},
	)
	if err != nil {
		panic(err)
	}
	// NOTE: We substract 1 from the number of channels because now
	// we have a dedicated collection for storing ClickHouse information
	// We should keep in mind to remove that whenever we change this.
	channelsCount := int64(len(channelsNames)) - 1

	// LogsPerChannel is obtained by counting the number of documents
	// that are in each collection.
	logsPerChannel := make(map[string]int64)
	for _, channel := range channelsNames {
		col := db.Collection(channel)
		count, err := col.CountDocuments(
			context.TODO(),
			bson.D{},
		)
		if err != nil {
			panic(err)
		}
		logsPerChannel[channel] = count
	}

	// SendersCount is computed when computing logsPerSender, see below.
	var senders []string

	// LogsPerSender is calculated by first enumerating all the different
	// distinct values of sender_id in the channels and for each of them,
	// we compute the number of documents in each channel.
	logsPerSender := make(map[string]int64)
	for _, channel := range channelsNames {
		col := db.Collection(channel)
		distinctSenderIdValues, err := col.Distinct(
			context.TODO(),
			"sender_id",
			bson.D{},
		)
		if err != nil {
			panic(err)
		}
		// We add (if not exists) new senders to the list of senders
		for _, sender := range distinctSenderIdValues {
			sender_ := sender.(string)
			for _, countedSender := range senders {
				if countedSender == sender_ {
					continue
				}
				senders = append(senders, sender_)
			}
		}
		for _, senderId := range distinctSenderIdValues {
			_senderId := senderId.(string)
			if _, ok := logsPerSender[_senderId]; !ok {
				logsPerSender[_senderId] = 0
			}
			count, err := col.CountDocuments(
				context.TODO(),
				bson.D{{"sender_id", _senderId}},
			)
			if err != nil {
				panic(err)
			}
			logsPerSender[_senderId] += count
		}
	}

	// LevelsCount is computed when computing logsPerLevel, see below
	var levels []string

	// LogsPerLevel is computed analogically to LogsPerSender
	logsPerLevel := make(map[string]int64)
	for _, channel := range channelsNames {
		col := db.Collection(channel)
		distinctLevelValues, err := col.Distinct(
			context.TODO(),
			"level",
			bson.D{},
		)
		if err != nil {
			panic(err)
		}
		// We add (if not exists) new senders to the list of levels
		for _, level := range distinctLevelValues {
			level_ := level.(string)
			for _, countedLevel := range senders {
				if countedLevel == level_ {
					continue
				}
				levels = append(senders, level_)
			}
		}
		for _, level := range distinctLevelValues {
			_level := level.(string)
			if _, ok := logsPerLevel[_level]; !ok {
				logsPerLevel[_level] = 0
			}
			count, err := col.CountDocuments(
				context.TODO(),
				bson.D{{"level", _level}},
			)
			if err != nil {
				panic(err)
			}
			logsPerLevel[_level] += count
		}
	}

	// TotalIngestedLogs is trivially computed
	var totalIngestedLogs int64 = 0
	for _, channel := range channelsNames {
		col := db.Collection(channel)
		count, err := col.CountDocuments(
			context.TODO(),
			bson.D{},
		)
		if err != nil {
			panic(err)
		}
		totalIngestedLogs += count
	}

	// TODO : TroughputPerTime / a little bit tricky to implement

	genObs := &GeneralObservables{
		ChannelsCount:     channelsCount,
		LogsPerChannel:    logsPerChannel,
		LogsPerSender:     logsPerSender,
		LogsPerLevel:      logsPerLevel,
		SendersCount:      int64(len(senders)),
		TotalIngestedLogs: totalIngestedLogs,
	}
	return genObs
}
