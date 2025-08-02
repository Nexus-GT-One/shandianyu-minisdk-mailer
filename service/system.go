package service

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/provider/mongodb"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"shandianyu-minisdk-mailer/util/httputil"
	"shandianyu-minisdk-mailer/util/systemutil"
)

type systemService struct{}

var SystemService = newSystemService()

func newSystemService() *systemService {
	return &systemService{}
}

func (s *systemService) SendClearCache() {
	systemutil.Goroutine(func() {
		for _, item := range arrayutil.First(SystemService.GetSystemConfig("system", "serverNodes")).Value.(primitive.A) {
			url := fmt.Sprintf("http://%v/api/v1/system/system/clearCache", item)
			httputil.Post(url, make(map[string]string), make([]byte, 0), make(map[string]any))
		}
	})
}

func (s *systemService) SaveLastMailIndex(lastMailIndex int64) {
	query := bson.D{{"module", "system"}, {"key", "lastMailIndex"}}
	db.UpdateMany(entity.SystemConfig{}, query, bson.D{{"value", lastMailIndex}}, options.Update().SetUpsert(true))
}

func (s *systemService) GetLastMailIndex() int64 {
	item := arrayutil.First(s.GetSystemConfig("system", "lastMailIndex"))
	if item == nil {
		return 0
	}
	return item.Value.(int64)
}

func (s *systemService) GetSystemConfig(module, key string) []*entity.SystemConfig {
	query := bson.D{{"module", module}, {"key", key}}
	ctx, cursor := db.Find(query, entity.SystemConfig{}, &options.FindOptions{Sort: bson.D{{"_id", 1}}})
	conf := mongodb.DecodeList(ctx, cursor, entity.SystemConfig{})
	return conf
}
