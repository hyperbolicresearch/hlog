package v1

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/hyperbolicresearch/hlog/config"
	"github.com/hyperbolicresearch/hlog/internal/mongodb"
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
			found := false
			for _, countedSender := range senders {
				if countedSender == sender_ {
					found = true
					break
				}
			}
			if !found {
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
			found := false
			for _, countedLevel := range levels {
				if countedLevel == level_ {
					found = true
					break
				}
			}
			if !found {
				levels = append(levels, level_)
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
		LevelsCount:       int64(len(levels)),
		TotalIngestedLogs: totalIngestedLogs,
	}
	return genObs
}

// ObservablesTail is the data structure that pushes the different
// types of Observables to the Observer.
type ObservablesTail struct {
	config *config.Config
}

// NewObservablesTail creates a new ObservablesTail instance
func NewObservablesTail(cfg *config.Config) *ObservablesTail {
	// NOTE: The system is now running only with MongoDB, and therefore,
	// there is no rush to make this work with different DBs yet.
	// But in the future, we should define a way to determine
	// how it will choose the observables.

	return &ObservablesTail{
		config: cfg,
	}
}

// Start will start periodically computing the observables and
// consequently push the observables through their respective loggers.
func (o *ObservablesTail) Start(sig chan os.Signal) error {
	ticker := time.NewTicker(o.config.APIv1.PushInterval)
	run := true
	for run {
		select {
		case <-sig:
			// Handle closing here
			run = false
			return nil
		case <-ticker.C:
			// At each push interval, we leverage the logger's ability
			// to write to io.Writer to actually broadcast the observables
			// to all the listeners aka websocket connections.
			// NOTE: These connections are manage by the API's endpoints.
			if o.config.SendGeneralObservables {
				genObs := GetMongoDBGeneralObservables(o.config.MongoDB)
				err := o.config.APIv1.GeneralObservablesLogger.Log(genObs)
				if err != nil {
					return err
				}
			}
		}

	}
	return nil
}
