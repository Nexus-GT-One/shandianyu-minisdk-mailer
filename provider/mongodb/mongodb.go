package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"regexp"
	"shandianyu-minisdk-mailer/provider/config"
	loggerFactory "shandianyu-minisdk-mailer/provider/logger"
	"shandianyu-minisdk-mailer/util/arrayutil"
	"shandianyu-minisdk-mailer/util/jsonutil"
	"shandianyu-minisdk-mailer/util/maputil"
	"shandianyu-minisdk-mailer/util/stringutil"
	"strconv"
	"strings"
	"sync"
	"time"
)

var connectionMap sync.Map
var logger = loggerFactory.GetLogger()

type mongodb struct {
	DB *mongo.Database
}

func init() {
	for _, conf := range config.GetList("data.mongodb") {
		name := fmt.Sprintf("%v", maputil.GetValueFromMap(conf.(map[string]any), "name", "default"))
		uri := fmt.Sprintf("%v", maputil.GetValueFromMap(conf.(map[string]any), "uri", ""))
		timeout, _ := strconv.Atoi(fmt.Sprintf("%v", maputil.GetValueFromMap(conf.(map[string]any), "timeout", 0)))

		if _, ok := connectionMap.Load(name); ok {
			logger.Fatal("only one default connection is allowed")
		}

		connectionMap.Store(name, &mongodb{DB: ConnectMongodb(name, uri, timeout)})
	}
}

func ConnectMongodb(name string, uri string, timeout int) *mongo.Database {
	conn, ok := connectionMap.Load(name)
	if ok {
		return conn.(*mongo.Database)
	}

	// 设置客户端连接配置
	_timeout := time.Duration(timeout) * time.Second
	ctx := context.Background()
	co := options.Client().ApplyURI(uri).SetTimeout(_timeout).SetSocketTimeout(_timeout).SetConnectTimeout(_timeout)
	// 连接到MongoDB
	client, err := mongo.Connect(ctx, co)
	if err != nil {
		logger.Fatal("%v", err)
	}
	// 检查连接
	err = client.Ping(ctx, nil)
	if err != nil {
		logger.Fatal("%v", err)
	}
	if co.Auth == nil {
		logger.Info("mongodb connect [%s] %s success ...", name, uri)
	} else {
		logger.Info("mongodb connect [%s] mongodb://*****:*****@%s/%s success ...", name, co.Hosts[0], dbName(uri))
	}
	return client.Database(dbName(uri))
}

func dbName(uri string) string {
	lastIndexSlash := strings.LastIndex(uri, "/") + 1
	lastIndexFactor := strings.LastIndex(uri, "?")
	if lastIndexFactor > 0 {
		return strings.TrimSpace(uri[lastIndexSlash:lastIndexFactor])
	}
	return strings.TrimSpace(uri[lastIndexSlash:])
}

func throwErrorIfNotNil(err error) {
	if err != nil {
		//panic(err.Error())
	}
}

func GetInstance() *mongodb {
	connection, ok := connectionMap.Load("default")
	if !ok {
		panic("unable to get [default] connection")
	}
	return connection.(*mongodb)
}

func GetLoggingInstance() *mongodb {
	connection, ok := connectionMap.Load("logging")
	if !ok {
		panic("unable to get [logging] connection")
	}
	return connection.(*mongodb)
}

// 创建索引
//
// entity-数据实体
func EnsureIndex[T any](entity T) {
	ctx := context.Background()
	typ := reflect.TypeOf(entity)
	connection := GetInstance()
	collection := typ.Name()
	for i := 0; i < typ.NumField(); i++ {
		indexJSON := strings.TrimSpace(typ.Field(i).Tag.Get("index"))
		if len(indexJSON) <= 0 {
			continue
		}

		indexMap := make(map[string]any)
		indexJSON = strings.ReplaceAll(indexJSON, "'", "\"")
		json.Unmarshal([]byte(indexJSON), &indexMap)
		if len(indexMap) <= 0 {
			// 忽略解析失败
			continue
		}

		keys := make(bson.D, 0)
		regexpPattern, _ := regexp.Compile("\\s+")
		indexJSON = regexpPattern.ReplaceAllString(indexJSON, "")
		keyString := indexJSON[strings.Index(indexJSON, `"keys":{`)+8 : strings.Index(indexJSON, `}`)]
		for _, key := range strings.Split(keyString, ",") {
			arr := strings.Split(key, ":")
			k := strings.ReplaceAll(strings.TrimSpace(arr[0]), `"`, "")
			v, _ := strconv.Atoi(strings.TrimSpace(arr[1]))
			keys = append(keys, bson.E{Key: k, Value: v})
		}

		indexName := fmt.Sprintf("%v", maputil.GetValueFromMap(indexMap, "name", ""))
		sparse, _ := strconv.ParseBool(fmt.Sprintf("%v", maputil.GetValueFromMap(indexMap, "sparse", false)))
		unique, _ := strconv.ParseBool(fmt.Sprintf("%v", maputil.GetValueFromMap(indexMap, "unique", false)))
		expireAfterSeconds, _ := strconv.Atoi(fmt.Sprintf("%v", maputil.GetValueFromMap(indexMap, "expireAfterSeconds", "0")))
		indexModel := mongo.IndexModel{Keys: keys}
		indexModel.Options = options.Index().SetName(indexName).SetBackground(true)
		if unique {
			indexModel.Options.SetUnique(unique)
		}
		if sparse {
			indexModel.Options.SetSparse(sparse)
		}
		if expireAfterSeconds > 0 {
			indexModel.Options.SetExpireAfterSeconds(int32(expireAfterSeconds))
		}

		_, err := connection.DB.Collection(collection).Indexes().CreateOne(ctx, indexModel)
		throwErrorIfNotNil(err)

		key_s, _ := json.Marshal(indexModel.Keys)
		option_s, _ := json.Marshal(indexModel.Options)
		logger.Info("collection [%s] create index: %v, index options: %v", collection, string(key_s), jsonutil.RemoveNullValues(string(option_s)))
	}
}

// 反序列化mongodb的数据到指定是数据实体
//
// ctx-上下文; cursor-mongodb返回的游标; entity-数据实体
func DecodeList[T any](ctx context.Context, cursor *mongo.Cursor, entity T) []*T {
	defer cursor.Close(ctx)
	result := make([]*T, 0)
	for cursor.Next(ctx) {
		var element T
		copier.Copy(element, entity)
		cursor.Decode(&element)
		result = append(result, &element)
	}
	return result
}

// 反序列化mongodb的数据到指定是数据实体
//
// ctx-上下文; cursor-mongodb返回的游标; entity-数据实体
func DecodeOne[T any](ctx context.Context, cursor *mongo.Cursor, entity T) *T {
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var element T
		copier.Copy(element, entity)
		cursor.Decode(&element)
		return &element
	}
	return nil
}

// 插入单条数据
//
// entity-数据实体
func (m *mongodb) InsertOne(entity any) primitive.ObjectID {
	return arrayutil.First(m.InsertMany(entity))
}

// 插入多条数据
//
// entity-数据实体(可以传数组或者传无限个单条)
func (m *mongodb) InsertMany(entity ...any) []primitive.ObjectID {
	ctx := context.Background()
	classType := reflect.TypeOf(arrayutil.First(entity))
	ignoreColumn := []string{"_id", "id", "createTime", "updateTime"}
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	documents := make([]any, 0)
	for _, item := range entity {
		document := make(bson.D, 0)
		pValue := reflect.ValueOf(item)
		for i := 0; i < classType.NumField(); i++ {
			omitempty := strings.Contains(classType.Field(i).Tag.Get("bson"), "omitempty")
			fieldName := stringutil.FirstLowerCase(classType.Field(i).Tag.Get("bson"))
			fieldValue := pValue.Field(i).Interface()
			isZero := reflect.ValueOf(fieldValue).IsZero()
			if arrayutil.Contains(ignoreColumn, fieldName) || (omitempty && isZero) {
				continue
			}
			document = append(document, bson.E{Key: fieldName, Value: fieldValue})
		}

		now := time.Now().UnixMilli()
		tmp := make([]map[string]any, 0)
		b, _ := json.Marshal(entity)
		json.Unmarshal(b, &tmp)
		if !maputil.ContainsKey(tmp[0], "CreateTime") {
			tmp[0]["CreateTime"] = now
		}
		document = append(document, bson.E{Key: "createTime", Value: maputil.GetValueFromMap(tmp[0], "CreateTime", nil)})
		document = append(document, bson.E{Key: "updateTime", Value: now})
		documents = append(documents, document)
	}

	result, err := c.InsertMany(ctx, documents)
	throwErrorIfNotNil(err)
	insertedIDs := make([]primitive.ObjectID, 0)
	for _, insertedID := range result.InsertedIDs {
		insertedIDs = append(insertedIDs, insertedID.(primitive.ObjectID))
	}
	return insertedIDs
}

// 更新单条数据
//
// entity-数据实体; id-主键id; update-需要更新的数据; opts-数据更新选项
func (m *mongodb) UpdateOne(entity any, id string, update bson.D, opts ...*options.UpdateOptions) *mongo.UpdateResult {
	objectID, err := primitive.ObjectIDFromHex(id)
	throwErrorIfNotNil(err)
	return m.UpdateMany(entity, bson.D{{"_id", objectID}}, update, opts...)
}

// 原子更新单条数据
//
// entity-数据实体; query-查询条件; update-需要更新的数据; operator-操作符; opts-数据更新选项
func (m *mongodb) FindOneAndUpdate(entity any, query, update bson.D, operator string, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	return c.FindOneAndUpdate(ctx, query, bson.D{
		{operator, update},
		{"$set", bson.D{{Key: "updateTime", Value: time.Now().UnixMilli()}}},
	}, opts...)
}

// 更新多条数据
//
// entity-数据实体; query-查询条件; update-需要更新的数据; opts-数据更新选项
func (m *mongodb) UpdateMany(entity any, query, update bson.D, opts ...*options.UpdateOptions) *mongo.UpdateResult {
	return m.doUpdateMany(entity, query, update, "$set", opts...)
}

// 删除字段
//
// entity-数据实体; query-查询条件; update-需要删除的字段; opts-数据更新选项
func (m *mongodb) Unset(entity any, query, update bson.D, opts ...*options.UpdateOptions) *mongo.UpdateResult {
	return m.doUpdateMany(entity, query, update, "$unset", opts...)
}

// 指定表名查询多条数据
//
// query-查询条件; collectionName-指定表名; opts-查询选项
func (m *mongodb) FindByCollectionName(query bson.D, collectionName string, opts ...*options.FindOptions) (context.Context, *mongo.Cursor) {
	ctx := context.Background()
	c := m.DB.Collection(collectionName)
	cursor, err := c.Find(ctx, query, opts...)
	throwErrorIfNotNil(err)
	return ctx, cursor
}

// 单条数据的字段自增
//
// entity-数据实体; id-主键id; update-需要自增的数据; opts-数据更新选项
func (m *mongodb) IncOne(entity any, id string, update bson.D, opts ...*options.UpdateOptions) *mongo.UpdateResult {
	objectID, err := primitive.ObjectIDFromHex(id)
	throwErrorIfNotNil(err)
	return m.doUpdateMany(entity, bson.D{{"_id", objectID}}, update, "$inc", opts...)
}

// 获取某个字段的唯一值
//
// entity-数据实体; fieldName-字段名; query-查询条件
func (m *mongodb) Distinct(entity any, fieldName string, query bson.D) []any {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	result, err := c.Distinct(ctx, fieldName, query)
	throwErrorIfNotNil(err)
	return result
}

// 更新多条数据
//
// entity-数据实体; query-查询条件; update-需要更新的数据; operator-操作符; opts-数据更新选项
func (m *mongodb) doUpdateMany(entity any, query, update bson.D, operator string, opts ...*options.UpdateOptions) *mongo.UpdateResult {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	updateResult, err := c.UpdateMany(ctx, query, bson.D{
		{operator, update},
		{"$set", bson.D{{Key: "updateTime", Value: time.Now().UnixMilli()}}},
	}, opts...)
	throwErrorIfNotNil(err)
	return updateResult
}

// 删除单条数据
//
// entity-数据实体; query-查询条件; opts-数据删除选项
func (m *mongodb) DeleteOne(entity any, query bson.D, opts ...*options.DeleteOptions) *mongo.DeleteResult {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	updateResult, err := c.DeleteOne(ctx, query, opts...)
	throwErrorIfNotNil(err)
	return updateResult
}

// 原子删除单条数据
//
// entity-数据实体; query-查询条件; opts-数据删除选项
func (m *mongodb) FindOneAndDelete(entity any, query bson.D, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	return c.FindOneAndDelete(ctx, query, opts...)
}

// 删除多条数据
//
// entity-数据实体; query-查询条件; opts-数据删除选项
func (m *mongodb) DeleteMany(entity any, query bson.D, opts ...*options.DeleteOptions) *mongo.DeleteResult {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	updateResult, err := c.DeleteMany(ctx, query, opts...)
	throwErrorIfNotNil(err)
	return updateResult
}

// 查询多条数据
//
// query-查询条件; entity-数据实体; opts-查询选项
func (m *mongodb) Find(query bson.D, entity any, opts ...*options.FindOptions) (context.Context, *mongo.Cursor) {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	cursor, err := c.Find(ctx, query, opts...)
	throwErrorIfNotNil(err)
	return ctx, cursor
}

// 查询单条数据
//
// query-查询条件; entity-数据实体
func (m *mongodb) FindOne(query bson.D, entity any) (context.Context, *mongo.Cursor) {
	opts := &options.FindOptions{}
	opts.SetLimit(1)
	return m.Find(query, entity, opts)
}

// 游标翻页
//
// query-查询条件; cursor-游标; pageSize-分页大小; sort-排序方式; entity-数据实体
func (m *mongodb) PaginationByCursor(query bson.D, cursor *string, pageSize int64, sort bson.D, entity any) (context.Context, *mongo.Cursor) {
	opts := &options.FindOptions{}
	opts.SetLimit(pageSize)
	if len(*cursor) > 0 {
		objectID, _ := primitive.ObjectIDFromHex(*cursor)
		query = bson.D{{"_id", bson.D{{"$lt", objectID}}}}
	}

	finalSort := bson.D{{"_id", -1}}
	if len(sort) > 0 {
		for _, s := range sort {
			finalSort = append(finalSort, s)
		}
	}
	opts.SetSort(finalSort)
	return m.Find(query, entity, opts)
}

// 页码翻页
//
// query-查询条件; page-当前页码; pageSize-分页大小; sort-排序方式; entity-数据实体
func (m *mongodb) PaginationByPage(query bson.D, page, pageSize int64, sort bson.D, entity any) (context.Context, *mongo.Cursor) {
	opts := &options.FindOptions{}
	opts.SetLimit(pageSize)
	opts.SetSkip((page - 1) * pageSize)
	if len(sort) <= 0 {
		sort = bson.D{{"_id", -1}}
	}
	opts.SetSort(sort)
	return m.Find(query, entity, opts)
}

// 统计总条数
//
// query-查询条件; opts-统计选项
func (m *mongodb) Count(entity any, query bson.D, opts ...*options.CountOptions) int64 {
	ctx := context.Background()
	classType := reflect.TypeOf(entity)
	c := m.DB.Collection(strings.Split(fmt.Sprintf("%v", classType), ".")[1])
	count, err := c.CountDocuments(ctx, query, opts...)
	throwErrorIfNotNil(err)
	return count
}
