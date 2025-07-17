package iputil

import (
	"encoding/json"
	"fmt"
	"github.com/oschwald/geoip2-golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"path"
	"reflect"
	IpType "shandianyu-minisdk-monitor/constant"
	"shandianyu-minisdk-monitor/entity"
	"shandianyu-minisdk-monitor/provider/config"
	"shandianyu-minisdk-monitor/provider/logger"
	"shandianyu-minisdk-monitor/provider/mongodb"
	"shandianyu-minisdk-monitor/thirdparty/feishu"
	"shandianyu-minisdk-monitor/util/arrayutil"
	"shandianyu-minisdk-monitor/util/httputil"
	"shandianyu-minisdk-monitor/util/systemutil"
	"strconv"
	"strings"
	"sync/atomic"
)

type IpSbResponse struct {
	Organization    string  `json:"organization"`
	Longitude       float64 `json:"longitude"`
	City            string  `json:"city"`
	Timezone        string  `json:"timezone"`
	Isp             string  `json:"isp"`
	Offset          int     `json:"offset"`
	Region          string  `json:"region"`
	Asn             int     `json:"asn"`
	AsnOrganization string  `json:"asn_organization"`
	Country         string  `json:"country"`
	Ip              string  `json:"ip"`
	Latitude        float64 `json:"latitude"`
	PostalCode      string  `json:"postal_code"`
	ContinentCode   string  `json:"continent_code"`
	CountryCode     string  `json:"country_code"`
	RegionCode      string  `json:"region_code"`
}

type Ip2LocationResponse struct {
	Ip          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionName  string  `json:"region_name"`
	CityName    string  `json:"city_name"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Asn         string  `json:"asn"`
	As          string  `json:"as"`
	IsProxy     bool    `json:"is_proxy"`
}

var counter int64
var db = mongodb.GetLoggingInstance()
var asnDB *geoip2.Reader
var cityDB *geoip2.Reader
var asnErr error
var cityErr error
var isProd = reflect.DeepEqual("prod", config.GetString("env"))
var ip2locationKeys = config.GetList("ip2location")

func init() {
	bashPath := config.GetString("data.geoip2.path")
	if len(bashPath) <= 0 {
		return
	}

	cityDbPath := path.Join(bashPath, "GeoLite2-City.mmdb")
	cityDB, cityErr = geoip2.Open(cityDbPath)
	if cityErr != nil {
		logger.GetLogger().Fatal("load city db error: %v", cityErr)
		return
	}
	logger.GetLogger().Info("load city db success: %v", cityDbPath)

	asnDbPath := path.Join(bashPath, "GeoLite2-ASN.mmdb")
	asnDB, asnErr = geoip2.Open(asnDbPath)
	if asnErr != nil {
		logger.GetLogger().Fatal("load asn db error: %v", asnErr)
		return
	}
	logger.GetLogger().Info("load asn db success: %v", asnDbPath)
}

func Location(ip string) *entity.Ip {
	ctx, cursor := db.FindOne(bson.D{{"ip", ip}}, entity.Ip{})
	ipInfo := mongodb.DecodeOne(ctx, cursor, entity.Ip{})
	if ipInfo != nil {
		return ipInfo
	}

	if len(ip) <= 0 || asnDB == nil || cityDB == nil {
		return &entity.Ip{}
	}

	parsedIp := net.ParseIP(ip)
	city, readCityErr := cityDB.City(parsedIp)
	if city == nil || (len(city.Country.Names) <= 0 && len(city.City.Names) <= 0) || readCityErr != nil {
		logger.GetLogger().Error("ip: %s, read city error: %v", ip, readCityErr)
		return GetByIpSb(ip)
	}

	asn, readAsnErr := asnDB.ASN(parsedIp)
	if asn == nil || readAsnErr != nil {
		logger.GetLogger().Error("ip: %s, read asn error: %v", ip, readAsnErr)
		return GetByIpSb(ip)
	}

	// 谷歌的ip认为是高风险ip
	isGoogleIp := strings.Contains(strings.ToLower(strings.TrimSpace(asn.AutonomousSystemOrganization)), "google")
	isAppleIp := strings.Contains(strings.ToLower(strings.TrimSpace(asn.AutonomousSystemOrganization)), "apple")
	isDangerIp := isGoogleIp || isAppleIp
	return record(&entity.Ip{
		Id:              primitive.ObjectID{},
		Ip:              ip,
		Type:            systemutil.If(isDangerIp, IpType.BLOCK, IpType.NORMAL).(IpType.IpType),
		Longitude:       city.Location.Longitude,
		City:            city.City.Names["en"],
		Timezone:        city.Location.TimeZone,
		Region:          arrayutil.First(city.Subdivisions).Names["en"],
		RegionCode:      arrayutil.First(city.Subdivisions).IsoCode,
		Asn:             asn.AutonomousSystemNumber,
		AsnOrganization: asn.AutonomousSystemOrganization,
		Country:         city.Country.Names["en"],
		Latitude:        city.Location.Latitude,
		PostalCode:      city.Postal.Code,
		ContinentCode:   city.Continent.Code,
		CountryCode:     city.Country.IsoCode,
	})
}

func GetByIp2location(ip string) *entity.Ip {
	// 非正式环境不浪费api额度
	if !isProd {
		return Location(ip)
	}

	ip2LocationResponse := &Ip2LocationResponse{}
	for i := 0; i < len(ip2locationKeys); i++ {
		index := atomic.AddInt64(&counter, 1) % int64(len(ip2locationKeys))
		key := ip2locationKeys[index]
		url := fmt.Sprintf("https://api.ip2location.io/?key=%s&ip=%s", key, ip)
		headers := map[string]string{"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"}
		json.Unmarshal(httputil.Get(url, headers), &ip2LocationResponse)
		if len(ip2LocationResponse.Ip) > 0 {
			asn, _ := strconv.Atoi(ip2LocationResponse.Asn)
			ipInfo := record(&entity.Ip{
				Ip:              ip,
				Longitude:       ip2LocationResponse.Longitude,
				City:            ip2LocationResponse.CityName,
				Timezone:        ip2LocationResponse.TimeZone,
				Region:          ip2LocationResponse.RegionName,
				Asn:             uint(asn),
				AsnOrganization: ip2LocationResponse.As,
				Country:         ip2LocationResponse.CountryName,
				Latitude:        ip2LocationResponse.Latitude,
				PostalCode:      ip2LocationResponse.ZipCode,
				CountryCode:     ip2LocationResponse.CountryCode,
			})

			return ipInfo
		}
	}

	feishu.AdminRobot().SendRobotMessage("请求ip2location错误: api额度已耗尽")

	// 如果都尝试不成功，就调用maxmind的
	return Location(ip)
}

func GetByIpSb(ip string) *entity.Ip {
	url := fmt.Sprintf("https://api.ip.sb/geoip/%s", ip)
	headers := map[string]string{"user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"}
	ipSbResponse := &IpSbResponse{}
	json.Unmarshal(httputil.Get(url, headers), &ipSbResponse)

	// 谷歌的ip认为是高风险ip
	isGoogleIp := strings.Contains(strings.ToLower(strings.TrimSpace(ipSbResponse.AsnOrganization)), "google")
	isAppleIp := strings.Contains(strings.ToLower(strings.TrimSpace(ipSbResponse.AsnOrganization)), "apple")
	isDangerIp := isGoogleIp || isAppleIp
	return record(&entity.Ip{
		Ip:              ip,
		Type:            systemutil.If(isDangerIp, IpType.BLOCK, IpType.NORMAL).(IpType.IpType),
		Longitude:       ipSbResponse.Longitude,
		City:            ipSbResponse.City,
		Timezone:        ipSbResponse.Timezone,
		Region:          ipSbResponse.Region,
		RegionCode:      ipSbResponse.RegionCode,
		Asn:             uint(ipSbResponse.Asn),
		AsnOrganization: ipSbResponse.AsnOrganization,
		Country:         ipSbResponse.Country,
		Latitude:        ipSbResponse.Latitude,
		PostalCode:      ipSbResponse.PostalCode,
		ContinentCode:   ipSbResponse.ContinentCode,
		CountryCode:     ipSbResponse.CountryCode,
	})
}

func getSystemConfig(module, key string) []*entity.SystemConfig {
	query := bson.D{{"module", module}, {"key", key}}
	ctx, cursor := mongodb.GetInstance().Find(query, entity.SystemConfig{}, &options.FindOptions{Sort: bson.D{{"_id", 1}}})
	conf := mongodb.DecodeList(ctx, cursor, entity.SystemConfig{})
	return conf
}

func record(ipInfo *entity.Ip) *entity.Ip {
	if len(ipInfo.Ip) <= 0 {
		return nil
	}

	update := bson.D{}
	data := make(map[string]any)
	b, _ := json.Marshal(ipInfo)
	json.Unmarshal(b, &data)
	delete(data, "id")
	delete(data, "_id")
	delete(data, "type")
	for key, value := range data {
		update = append(update, bson.E{Key: key, Value: value})
	}
	db.UpdateMany(entity.Ip{}, bson.D{{"ip", ipInfo.Ip}}, update, options.Update().SetUpsert(true))
	return ipInfo
}

func SendClearCache(ip []string) {
	if len(ip) <= 0 {
		return
	}
	systemutil.Goroutine(func() {
		for _, item := range arrayutil.First(getSystemConfig("system", "serverNodes")).Value.(primitive.A) {
			url := fmt.Sprintf("http://%v/api/v1/system/system/clearCache", item)
			body, _ := json.Marshal(map[string]any{"ip": ip})
			httputil.Post(url, make(map[string]string), body, make(map[string]any))
		}
	})
}
