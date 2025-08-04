package service

import (
	_ "embed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"shandianyu-minisdk-mailer/entity"
	"shandianyu-minisdk-mailer/provider/mongodb"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"shandianyu-minisdk-mailer/util/systemutil"
	"time"
)

type gameService struct{}

var GameService = newGameService()

func newGameService() *gameService {
	return &gameService{}
}

func (a *gameService) GetByName(name string) *entity.Game {
	query := bson.D{{"name", name}}
	ctx, cursor := db.FindOne(query, entity.Game{})
	return mongodb.DecodeOne(ctx, cursor, entity.Game{})
}

func (a *gameService) GetBySymbol(symbol string) *entity.Game {
	query := bson.D{{"symbol", symbol}}
	ctx, cursor := db.FindOne(query, entity.Game{})
	return mongodb.DecodeOne(ctx, cursor, entity.Game{})
}

func (a *gameService) GetOneGameById(id string) *entity.Game {
	objectID, _ := primitive.ObjectIDFromHex(id)
	ctx, cursor := db.FindOne(bson.D{{"_id", objectID}}, entity.Game{})
	return mongodb.DecodeOne(ctx, cursor, entity.Game{})
}

func (a *gameService) UpdateGame(gameId string, updates bson.D) *entity.Game {
	db.UpdateOne(entity.Game{}, gameId, updates)
	return a.GetOneGameById(gameId)
}

func (a *gameService) GetByDeveloperEmail(developerEmail string) []*entity.Game {
	query := bson.D{{"channel", "iOS"}, {"developerEmail", developerEmail}}
	ctx, cursor := db.Find(query, entity.Game{}, &options.FindOptions{Sort: bson.D{{"_id", -1}}})
	return mongodb.DecodeList(ctx, cursor, entity.Game{})
}

func (s *gameService) RecordGameOperateHistory(gameId, operator, _type, remark string) {
	if len(remark) <= 0 {
		return
	}
	db.InsertOne(entity.GameOperateHistory{GameId: gameId, Operator: operator, Type: _type, Remark: remark, CreateTime: time.Now().UnixMilli()})
}

func (a *gameService) GetAuditingVersion(oneGame *entity.Game) string {
	for appVersion, isAudit := range oneGame.Audit {
		if !isAudit {
			continue
		}
		return appVersion
	}

	return a.GetLastSubmitVersion(oneGame)
}

func (p *gameService) GetLastSubmitAutoReleaseDays(bundleId, appVersion string) int {
	ctx, cursor := mongodb.GetLoggingInstance().Find(bson.D{{"bundleId", bundleId}, {"appVersion", appVersion}, {"actionType", "_submit_"}}, entity.AuditTrack{}, &options.FindOptions{Sort: bson.D{{"createTime", -1}}})
	item := arrayutil.First(mongodb.DecodeList(ctx, cursor, entity.AuditTrack{}))
	if item == nil {
		return 0
	}
	return item.AutoReleaseDays
}

func (p *gameService) GetLastSubmitVersion(game *entity.Game) string {
	query := bson.D{
		{"bundleId", game.BundleId},
		{"actionType", "_submit_"},
	}
	ctx, cursor := mongodb.GetLoggingInstance().Find(query, entity.AuditTrack{}, &options.FindOptions{Sort: bson.D{{"_id", -1}}})
	auditTrack := mongodb.DecodeOne(ctx, cursor, entity.AuditTrack{})
	if auditTrack == nil {
		return ""
	}

	return auditTrack.AppVersion
}

func (p *gameService) GetLastSubmitBuildNum(game *entity.Game) int {
	query := bson.D{
		{"bundleId", game.BundleId},
		{"actionType", "_submit_"},
		{"buildNum", bson.D{{"$exists", true}}},
	}
	ctx, cursor := mongodb.GetLoggingInstance().Find(query, entity.AuditTrack{}, &options.FindOptions{Sort: bson.D{{"_id", -1}}})
	auditTrack := mongodb.DecodeOne(ctx, cursor, entity.AuditTrack{})
	if auditTrack == nil {
		return 0
	}

	return auditTrack.BuildNum
}

// 记录过审
func (p *gameService) RecordApproved(bundleId, appVersion string) {
	systemutil.Goroutine(func() {
		// 查出build号
		query := bson.D{
			{"bundleId", bundleId},
			{"appVersion", appVersion},
			{"actionType", "_submit_"},
		}
		ctx, cursor := mongodb.GetLoggingInstance().Find(query, entity.AuditTrack{}, &options.FindOptions{Sort: bson.D{{"_id", -1}}})
		submitTrack := mongodb.DecodeOne(ctx, cursor, entity.AuditTrack{})
		if submitTrack == nil {
			return
		}

		// 查询是否已经标记过审
		query = bson.D{
			{"bundleId", bundleId},
			{"appVersion", appVersion},
			{"actionType", "_approved_"},
			{"buildNum", submitTrack.BuildNum},
		}
		ctx, cursor = mongodb.GetLoggingInstance().Find(query, entity.AuditTrack{}, &options.FindOptions{Sort: bson.D{{"_id", -1}}})
		approvedTrack := mongodb.DecodeOne(ctx, cursor, entity.AuditTrack{})
		if approvedTrack != nil {
			// 已标记的返回
			return
		}

		mongodb.GetLoggingInstance().InsertOne(entity.AuditTrack{
			BundleId:   bundleId,
			ActionType: "_approved_",
			AppVersion: appVersion,
			BuildNum:   submitTrack.BuildNum,
			Remark:     "邮件确认过审",
			CreateTime: time.Now().UnixMilli(),
		})
	})
}

// 记录审核不通过
func (p *gameService) RecordRejected(bundleId, appVersion string) {
	systemutil.Goroutine(func() {
		// 如果已经生成过，就不再生成了
		query := bson.D{
			{"bundleId", bundleId},
			{"appVersion", appVersion},
			{"actionType", "_rejected_"},
		}
		ctx, cursor := mongodb.GetLoggingInstance().Find(query, entity.AuditTrack{}, &options.FindOptions{Sort: bson.D{{"_id", -1}}})
		rejectedTrack := mongodb.DecodeOne(ctx, cursor, entity.AuditTrack{})
		if rejectedTrack != nil {
			return
		}

		mongodb.GetLoggingInstance().InsertOne(entity.AuditTrack{
			BundleId:   bundleId,
			ActionType: "_reject_",
			AppVersion: appVersion,
			BuildNum:   rejectedTrack.BuildNum,
			Remark:     "邮件确认被拒",
			CreateTime: time.Now().UnixMilli(),
		})
	})
}
