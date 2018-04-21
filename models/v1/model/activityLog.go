package model

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DanielRenne/GoCore/core"
	"github.com/DanielRenne/GoCore/core/dbServices"
	"github.com/DanielRenne/GoCore/core/fileCache"
	"github.com/DanielRenne/GoCore/core/logger"
	"github.com/DanielRenne/GoCore/core/serverSettings"
	"github.com/DanielRenne/GoCore/core/store"
	"github.com/DanielRenne/GoCore/core/utils"
	"github.com/asdine/storm"
	"github.com/globalsign/mgo/bson"
	"log"
	"reflect"
	"sync"
	"time"
)

var ActivityLogs modelActivityLogs

type modelActivityLogs struct{}

var collectionActivityLogsMutex *sync.RWMutex

type ActivityLogJoinItems struct {
	Count int            `json:"Count"`
	Items *[]ActivityLog `json:"Items"`
}

var GoCoreActivityLogsHasBootStrapped bool

func init() {
	collectionActivityLogsMutex = &sync.RWMutex{}

	//ActivityLogs.Index()
	go func() {
		time.Sleep(time.Second * 5)
		ActivityLogs.Bootstrap()
	}()
	store.RegisterStore(ActivityLogs)
}

func (self *ActivityLog) GetId() string {
	return self.Id.Hex()
}

type ActivityLog struct {
	Id            bson.ObjectId  `json:"Id" storm:"id"`
	AccountId     string         `json:"AccountId" storm:"index"`
	UserId        string         `json:"UserId" storm:"index"`
	Entity        string         `json:"Entity"`
	EntityId      string         `json:"EntityId"`
	Action        string         `json:"Action"`
	Value         string         `json:"Value"`
	CreateDate    time.Time      `json:"CreateDate" bson:"CreateDate"`
	UpdateDate    time.Time      `json:"UpdateDate" bson:"UpdateDate"`
	LastUpdateId  string         `json:"LastUpdateId" bson:"LastUpdateId"`
	BootstrapMeta *BootstrapMeta `json:"BootstrapMeta" bson:"-"`

	Errors struct {
		Id        string `json:"Id"`
		AccountId string `json:"AccountId"`
		UserId    string `json:"UserId"`
		Entity    string `json:"Entity"`
		EntityId  string `json:"EntityId"`
		Action    string `json:"Action"`
		Value     string `json:"Value"`
	} `json:"Errors" bson:"-"`

	Joins struct {
		Account *Account `json:"Account,omitempty" join:"Accounts,Account,AccountId,false,"`
	} `json:"Joins" bson:"-"`
}

func (self modelActivityLogs) Single(field string, value interface{}) (retObj ActivityLog, e error) {
	e = dbServices.BoltDB.One(field, value, &retObj)
	return
}

func (obj modelActivityLogs) Search(field string, value interface{}) (retObj []ActivityLog, e error) {
	e = dbServices.BoltDB.Find(field, value, &retObj)
	if len(retObj) == 0 {
		retObj = []ActivityLog{}
	}
	return
}

func (obj modelActivityLogs) SearchAdvanced(field string, value interface{}, limit int, skip int) (retObj []ActivityLog, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj)
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Find(field, value, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	return
}

func (obj modelActivityLogs) All() (retObj []ActivityLog, e error) {
	e = dbServices.BoltDB.All(&retObj)
	if len(retObj) == 0 {
		retObj = []ActivityLog{}
	}
	return
}

func (obj modelActivityLogs) AllAdvanced(limit int, skip int) (retObj []ActivityLog, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.All(&retObj)
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.All(&retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	return
}

func (obj modelActivityLogs) AllByIndex(index string) (retObj []ActivityLog, e error) {
	e = dbServices.BoltDB.AllByIndex(index, &retObj)
	if len(retObj) == 0 {
		retObj = []ActivityLog{}
	}
	return
}

func (obj modelActivityLogs) AllByIndexAdvanced(index string, limit int, skip int) (retObj []ActivityLog, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj)
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.AllByIndex(index, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	return
}

func (obj modelActivityLogs) Range(min, max, field string) (retObj []ActivityLog, e error) {
	e = dbServices.BoltDB.Range(field, min, max, &retObj)
	if len(retObj) == 0 {
		retObj = []ActivityLog{}
	}
	return
}

func (obj modelActivityLogs) RangeAdvanced(min, max, field string, limit int, skip int) (retObj []ActivityLog, e error) {
	if limit == 0 && skip == 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj)
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if limit > 0 && skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit), storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if limit > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Limit(limit))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	if skip > 0 {
		e = dbServices.BoltDB.Range(field, min, max, &retObj, storm.Skip(skip))
		if len(retObj) == 0 {
			retObj = []ActivityLog{}
		}
		return
	}
	return
}

func (obj modelActivityLogs) ById(objectID interface{}, joins []string) (value reflect.Value, err error) {
	var retObj ActivityLog
	q := obj.Query()
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.ById(objectID, &retObj)
	value = reflect.ValueOf(&retObj)
	return
}
func (obj modelActivityLogs) NewByReflection() (value reflect.Value) {
	retObj := ActivityLog{}
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelActivityLogs) ByFilter(filter map[string]interface{}, inFilter map[string]interface{}, excludeFilter map[string]interface{}, joins []string) (value reflect.Value, err error) {
	var retObj []ActivityLog
	q := obj.Query().Filter(filter)
	if len(inFilter) > 0 {
		q = q.In(inFilter)
	}
	if len(excludeFilter) > 0 {
		q = q.Exclude(excludeFilter)
	}
	for i := range joins {
		joinValue := joins[i]
		q = q.Join(joinValue)
	}
	err = q.All(&retObj)
	value = reflect.ValueOf(&retObj)
	return
}

func (obj modelActivityLogs) Query() *Query {
	query := new(Query)
	var elapseMs int
	for {
		collectionActivityLogsMutex.RLock()
		bootstrapped := GoCoreActivityLogsHasBootStrapped
		collectionActivityLogsMutex.RUnlock()

		if bootstrapped {
			break
		}
		elapseMs = elapseMs + 2
		time.Sleep(time.Millisecond * 2)
		if elapseMs%10000 == 0 {
			log.Println("ActivityLogs has not bootstrapped and has yet to get a collection pointer")
		}
	}
	query.collectionName = "ActivityLogs"
	query.entityName = "ActivityLog"
	return query
}
func (obj modelActivityLogs) Index() error {
	return dbServices.BoltDB.Init(&ActivityLog{})
}

func (obj modelActivityLogs) BootStrapComplete() {
	collectionActivityLogsMutex.Lock()
	GoCoreActivityLogsHasBootStrapped = true
	collectionActivityLogsMutex.Unlock()
}
func (obj modelActivityLogs) Bootstrap() error {
	start := time.Now()
	defer func() {
		log.Println(logger.TimeTrack(start, "Bootstraping of ActivityLogs Took"))
	}()
	if serverSettings.WebConfig.Application.BootstrapData == false {
		obj.BootStrapComplete()
		return nil
	}

	var isError bool
	var query Query
	var rows []ActivityLog
	cnt, errCount := query.Count(&rows)
	if errCount != nil {
		cnt = 1
	}

	dataString := ""

	var files [][]byte
	var err error
	var distDirectoryFound bool
	err = fileCache.LoadCachedBootStrapFromKeyIntoMemory(serverSettings.WebConfig.Application.ProductName + "ActivityLogs")
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for ActivityLogs due to caching issue: " + err.Error())
		return err
	}

	files, err, distDirectoryFound = BootstrapDirectory("activityLogs", cnt)
	if err != nil {
		obj.BootStrapComplete()
		log.Println("Failed to bootstrap data for ActivityLogs: " + err.Error())
		return err
	}

	if dataString != "" {
		data, err := base64.StdEncoding.DecodeString(dataString)
		if err != nil {
			obj.BootStrapComplete()
			log.Println("Failed to bootstrap data for ActivityLogs: " + err.Error())
			return err
		}
		files = append(files, data)
	}

	var v []ActivityLog
	for _, file := range files {
		var fileBootstrap []ActivityLog
		hash := md5.Sum(file)
		hexString := hex.EncodeToString(hash[:])
		err = json.Unmarshal(file, &fileBootstrap)
		if !fileCache.DoesHashExistInCache(serverSettings.WebConfig.Application.ProductName+"ActivityLogs", hexString) || cnt == 0 {
			if err != nil {

				logger.Message("Failed to bootstrap data for ActivityLogs: "+err.Error(), logger.RED)
				utils.TalkDirtyToMe("Failed to bootstrap data for ActivityLogs: " + err.Error())
				continue
			}

			fileCache.UpdateBootStrapMemoryCache(serverSettings.WebConfig.Application.ProductName+"ActivityLogs", hexString)

			for i, _ := range fileBootstrap {
				fb := fileBootstrap[i]
				v = append(v, fb)
			}
		}
	}
	fileCache.WriteBootStrapCacheFile(serverSettings.WebConfig.Application.ProductName + "ActivityLogs")

	var actualCount int
	originalCount := len(v)
	log.Println("Total count of records attempting ActivityLogs", len(v))

	for _, doc := range v {
		var original ActivityLog
		if doc.Id.Hex() == "" {
			doc.Id = bson.NewObjectId()
		}
		err = query.ById(doc.Id, &original)
		if err != nil || (err == nil && doc.BootstrapMeta != nil && doc.BootstrapMeta.AlwaysUpdate) || "EquipmentCatalog" == "ActivityLogs" {
			if doc.BootstrapMeta != nil && doc.BootstrapMeta.DeleteRow {
				err = doc.Delete()
				if err != nil {
					log.Println("Failed to delete data for ActivityLogs:  " + doc.Id.Hex() + "  " + err.Error())
					isError = true
				}
			} else {
				valid := 0x01
				var reason map[string]bool
				reason = make(map[string]bool, 0)

				if doc.BootstrapMeta != nil && doc.BootstrapMeta.Version > 0 && doc.BootstrapMeta.Version <= serverSettings.WebConfig.Application.VersionNumeric {
					valid &= 0x00
					reason["Version Mismatch"] = true
				}
				if doc.BootstrapMeta != nil && doc.BootstrapMeta.Domain != "" && doc.BootstrapMeta.Domain != serverSettings.WebConfig.Application.ServerFQDN {
					valid &= 0x00
					reason["FQDN Mismatch With Domain"] = true
				}
				if doc.BootstrapMeta != nil && len(doc.BootstrapMeta.Domains) > 0 && !utils.InArray(serverSettings.WebConfig.Application.ServerFQDN, doc.BootstrapMeta.Domains) {
					valid &= 0x00
					reason["FQDN Mismatch With Domains"] = true
				}
				if doc.BootstrapMeta != nil && doc.BootstrapMeta.ProductName != "" && doc.BootstrapMeta.ProductName != serverSettings.WebConfig.Application.ProductName {
					valid &= 0x00
					reason["ProductName does not Match"] = true
				}
				if doc.BootstrapMeta != nil && len(doc.BootstrapMeta.ProductNames) > 0 && !utils.InArray(serverSettings.WebConfig.Application.ProductName, doc.BootstrapMeta.ProductNames) {
					valid &= 0x00
					reason["ProductNames does not Match Product"] = true
				}
				if doc.BootstrapMeta != nil && doc.BootstrapMeta.ReleaseMode != "" && doc.BootstrapMeta.ReleaseMode != serverSettings.WebConfig.Application.ReleaseMode {
					valid &= 0x00
					reason["ReleaseMode does not match"] = true
				}

				if valid == 0x01 {
					actualCount += 1
					err = doc.Save()
					if err != nil {
						log.Println("Failed to bootstrap data for ActivityLogs:  " + doc.Id.Hex() + "  " + err.Error())
						isError = true
					}
				} else if serverSettings.WebConfig.Application.ReleaseMode == "development" {
					log.Println("ActivityLogs skipped a row for some reason on " + doc.Id.Hex() + " because of " + core.Debug.GetDump(reason))
				}
			}
		} else {
			actualCount += 1
		}
	}
	if isError {
		log.Println("FAILED to bootstrap ActivityLogs")
	} else {

		if distDirectoryFound == false {
			err = BootstrapMongoDump("activityLogs", "ActivityLogs")
		}
		if err == nil {
			log.Println("Successfully bootstrapped ActivityLogs")
			if actualCount != originalCount {
				logger.Message("ActivityLogs counts are different than original bootstrap and actual inserts, please inpect data."+core.Debug.GetDump("Actual", actualCount, "OriginalCount", originalCount), logger.RED)
			}
		}
	}
	obj.BootStrapComplete()
	return nil
}

func (obj modelActivityLogs) New() *ActivityLog {
	return &ActivityLog{}
}

func (obj *ActivityLog) NewId() {
	obj.Id = bson.NewObjectId()
}

func (self *ActivityLog) Save() error {
	if self.Id == "" {
		self.Id = bson.NewObjectId()
	}
	t := time.Now()
	self.CreateDate = t
	self.UpdateDate = t
	dbServices.CollectionCache{}.Remove("ActivityLogs", self.Id.Hex())
	return dbServices.BoltDB.Save(self)
}

func (self *ActivityLog) SaveWithTran(t *Transaction) error {

	return self.CreateWithTran(t, false)
}
func (self *ActivityLog) ForceCreateWithTran(t *Transaction) error {

	return self.CreateWithTran(t, true)
}
func (self *ActivityLog) CreateWithTran(t *Transaction, forceCreate bool) error {

	dbServices.CollectionCache{}.Remove("ActivityLogs", self.Id.Hex())
	return self.Save()
}

func (self *ActivityLog) ValidateAndClean() error {

	return validateFields(ActivityLog{}, self, reflect.ValueOf(self).Elem())
}

func (self *ActivityLog) Reflect() []Field {

	return Reflect(ActivityLog{})
}

func (self *ActivityLog) Delete() error {
	dbServices.CollectionCache{}.Remove("ActivityLogs", self.Id.Hex())
	return dbServices.BoltDB.Delete("ActivityLog", self.Id.Hex())
}

func (self *ActivityLog) DeleteWithTran(t *Transaction) error {
	dbServices.CollectionCache{}.Remove("ActivityLogs", self.Id.Hex())
	return dbServices.BoltDB.Delete("ActivityLogs", self.Id.Hex())
}

func (self *ActivityLog) JoinFields(remainingRecursions string, q *Query, recursionCount int) (err error) {

	source := reflect.ValueOf(self).Elem()

	var joins []join
	joins, err = getJoins(source, remainingRecursions)

	if len(joins) == 0 {
		return
	}

	s := source
	for _, j := range joins {
		id := reflect.ValueOf(q.CheckForObjectId(s.FieldByName(j.joinFieldRefName).Interface())).String()
		joinsField := s.FieldByName("Joins")
		setField := joinsField.FieldByName(j.joinFieldName)

		endRecursion := false
		if serverSettings.WebConfig.Application.LogJoinQueries {
			fmt.Print("Remaining Recursions")
			fmt.Println(fmt.Sprintf("%+v", remainingRecursions))
			fmt.Println(fmt.Sprintf("%+v", j.collectionName))
		}
		if remainingRecursions == j.joinSpecified {
			endRecursion = true
		}
		err = joinField(j, id, setField, j.joinSpecified, q, endRecursion, recursionCount)
		if err != nil {
			return
		}
	}
	return
}

func (self *ActivityLog) Unmarshal(data []byte) error {

	err := bson.Unmarshal(data, &self)
	if err != nil {
		return err
	}
	return nil
}

func (obj *ActivityLog) JSONString() (string, error) {
	bytes, err := json.Marshal(obj)
	return string(bytes), err
}

func (obj *ActivityLog) JSONBytes() ([]byte, error) {
	return json.Marshal(obj)
}

func (obj *ActivityLog) BSONString() (string, error) {
	bytes, err := bson.Marshal(obj)
	return string(bytes), err
}

func (obj *ActivityLog) BSONBytes() (in []byte, err error) {
	err = bson.Unmarshal(in, obj)
	return
}

func (obj *ActivityLog) ParseInterface(x interface{}) (err error) {
	data, err := json.Marshal(x)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, obj)
	return
}
func (obj modelActivityLogs) ReflectByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "Entity":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "EntityId":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Action":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Value":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Id":
		obj, ok := x.(bson.ObjectId)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "AccountId":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "UserId":
		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	}
	return
}

func (obj modelActivityLogs) ReflectBaseTypeByFieldName(fieldName string, x interface{}) (value reflect.Value, err error) {

	switch fieldName {
	case "Id":
		if x == nil {
			var obj bson.ObjectId
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(bson.ObjectId)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "AccountId":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "UserId":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Entity":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "EntityId":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Action":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	case "Value":
		if x == nil {
			var obj string
			value = reflect.ValueOf(obj)
			return
		}

		obj, ok := x.(string)
		if !ok {
			err = errors.New("Failed to typecast interface.")
			return
		}
		value = reflect.ValueOf(obj)
		return
	}
	return
}
