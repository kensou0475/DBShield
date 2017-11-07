package db

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
	// mysql orm
	"github.com/astaxie/beego/orm"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/qiwihui/DBShield/dbshield/logger"
	"github.com/qiwihui/DBShield/dbshield/sql"
)

//MySQL local db
type MySQL struct {
	name string
	UUID string
}

//QueryAction 记录所有操作
type QueryAction struct {
	ID        int    `orm:"column(id)"`
	SessionID string `orm:"column(flow_id);null;size(32)"`
	FlowInfo  string `orm:"column(flow_info);null;type(text)"`
	// 实际查询语句
	Query string `orm:"column(query);null;type(text)"`
	// 查询用户
	User string `orm:"column(user);null;size(128)"`
	// 查询客户端信息
	ClientIP      string `orm:"column(client_ip);null;size(39)"`
	ClientProgram string `orm:"column(client_program);null;size(128)"`
	// server info
	ServerIP   string `orm:"column(server_ip);null;size(39)"`
	ServerPort int    `orm:"column(server_port);null"`
	// 执行的数据库和表
	Database string `orm:"column(db);null;size(128)"`
	Tables   string `orm:"column(tables);null;type(text)"`
	// 执行时间和执行耗时(ms)
	Time     time.Time `orm:"column(time);auto_now_add;type(datetime);size(6)"`
	Duration int64     `orm:"column(duration);default(0)"`
	// 执行结果
	QueryResult bool `orm:"column(query_result);default(true)"`
	// 是否违规操作
	IsAbnormal bool `orm:"column(is_abnormal);default(false)"`
	// 违规操作类型：none, pattern, permission
	AbnormalType string `orm:"column(abnormal_type);size(32);default(none)"`
	// 处理结果：none, learning, pass, drop
	Action string `orm:"column(action);size(36);defult(pass)"`
	// 告警
	IsAlarm bool `orm:"column(is_alarm);default(false)"`
	// 是否分析
	Analysed bool `orm:"column(analysed);default(false)"`
	// sql type
	SQLType string `orm:"column(sql_type);null;size(32)"`
	// dbshield or others
	Tool string `orm:"column(tool);null;size(32)"`
	// 模式
	Pattern string `orm:"column(pattern);null;type(text)"`
	// 区分不同
	UUID string `orm:"column(uuid);size(36)"`
}

//Pattern record trainging set
type Pattern struct {
	ID int `orm:"column(id)"`
	// pattent_key
	Key string `orm:"column(key);null;type(text)"`
	//value
	Value string `orm:"column(value);null;type(text)"`
	// Example Value
	ExampleValue string `orm:"column(example_value);null;type(text)"`
	// 启用状态, true, false
	Enable bool   `orm:"column(enable);default(true)"`
	UUID   string `orm:"column(uuid);size(36)"`
}

//State record abnormal set
type State struct {
	ID              int    `orm:"column(id)"`
	Key             string `orm:"column(key);size(5)"`
	QueryCounter    uint64 `orm:"column(QueryCounter);type(bigint unsigned)"`
	AbnormalCounter uint64 `orm:"column(AbnormalCounter);type(bigint unsigned)"`
	UUID            string `orm:"column(uuid);size(36)"`
}

//Permission 权限规则
type Permission struct {
	ID int `orm:"column(id)"`
	// 数据库
	Db string `orm:"column(db);null;size(128)"`
	// 用户
	User string `orm:"column(user);null;size(128)"`
	// 客户端
	Client string `orm:"column(client);null;size(128)"`
	// 表, "*" 表示全部
	Table string `orm:"column(table);null;size(128)"`
	// 权限, SELECT,UPDATE,DELETE,INSERT,GRANT....
	Permission string `orm:"column(permission);type(text)"`
	// 启用状态, true, false
	Enable bool   `orm:"column(enable);default(true)"`
	UUID   string `orm:"column(uuid);size(36)"`
}

// RecordQueryAction record query and action
func (m *MySQL) RecordQueryAction(context sql.QueryAction) error {
	logger.Debugf("action: %s", context.Action)

	// 异步记录
	go func() {

		// ms
		elapsedMs := context.Duration.Nanoseconds() / 1e6
		// table name
		tables, _ := sql.ExtractTableNames(string(context.Query))
		var tableString string
		if len(tables) > 0 {
			tableString = strings.Join(tables, ",")
		} else {
			tableString = ""
		}

		o := orm.NewOrm()
		var queryAction QueryAction
		queryAction.SessionID = context.SessionID
		queryAction.Query = string(context.Query)
		queryAction.User = string(context.User)
		queryAction.ClientIP = fourByteBigEndianToIP(context.Client)
		queryAction.ClientProgram = ""
		queryAction.ServerIP = context.ServerIP
		queryAction.ServerPort = int(context.ServerPort)
		queryAction.Database = string(context.Database)
		queryAction.Tables = tableString
		//TODO result
		queryAction.SQLType = context.QueryType
		queryAction.Tool = "dbshield"
		queryAction.QueryResult = true
		queryAction.Time = context.Time
		queryAction.Action = context.Action
		queryAction.Duration = elapsedMs
		queryAction.Pattern = formatPattern(sql.Pattern(context.Query))
		queryAction.UUID = m.UUID
		id, err := o.Insert(&queryAction)
		if err != nil {
			logger.Warningf("RecordQuery: %s", err.Error())
		} else {
			logger.Debugf("Query saved, ID: %d", id)
		}
	}()
	return nil
}

// RecordAbnormal record abnormal query
func (m *MySQL) RecordAbnormal(context sql.QueryContext, abType string) error {
	atomic.AddUint64(&AbnormalCounter, 1)
	go func() {
		// table name
		tables, _ := sql.ExtractTableNames(string(context.Query))
		queryType := sql.GetStmtType(string(context.Query))
		var tableString string
		if len(tables) > 0 {
			tableString = strings.Join(tables, ",")
		} else {
			tableString = ""
		}

		o := orm.NewOrm()
		var abnormal QueryAction
		// var sx16 = formatPattern(context.Marshal())  // pattern
		abnormal.SessionID = context.SessionID
		abnormal.Query = string(context.Query)
		abnormal.Tool = "dbshield"
		abnormal.SQLType = queryType
		abnormal.User = string(context.User)
		abnormal.ClientIP = fourByteBigEndianToIP(context.Client)
		abnormal.ClientProgram = ""
		abnormal.ServerIP = context.ServerIP
		abnormal.ServerPort = int(context.ServerPort)
		abnormal.Database = string(context.Database)
		abnormal.Tables = tableString
		abnormal.Time = context.Time
		abnormal.Duration = 0
		abnormal.QueryResult = false
		abnormal.IsAbnormal = true
		abnormal.AbnormalType = abType
		abnormal.IsAlarm = false
		abnormal.Action = "drop"
		abnormal.Pattern = formatPattern(sql.Pattern(context.Query))
		abnormal.UUID = m.UUID
		id, err := o.Insert(&abnormal)
		if err == nil {
			logger.Debugf("Abnormal saved, ID: %d", id)
		} else {
			logger.Warningf("Abnormal save error: %s", err.Error())
		}
	}()
	return nil
}

// CheckPattern check if pattern exists
func (m *MySQL) CheckPattern(pattern []byte) error {

	return errors.New("Not Impletement")
}

// PutPattern put pattern
func (m *MySQL) PutPattern(pattern []byte, query []byte) error {

	return errors.New("Not Impletement")
}

// DeletePattern delete pattern
func (m *MySQL) DeletePattern(pattern []byte) error {
	go func() {
		o := orm.NewOrm()
		if num, err := o.Delete(&Pattern{Key: string(pattern), UUID: m.UUID}); err == nil {
			logger.Debugf("Pattern delete, num: %d", num)
		} else {
			logger.Warningf("Pattern delete error: %s", err.Error())
		}
	}()
	return nil
}

// Purge local databases
func (m *MySQL) Purge() error {
	o := orm.NewOrm()
	_, err := o.Raw("DROP TABLE IF EXISTS pattern, query_action, abnormal, state;").Exec()
	if err != nil {
		return err
	}
	logger.Warningf("All tables dropped")
	return nil
}

// SyncAndClose local databases
func (m *MySQL) SyncAndClose() error {
	// 由 go-sql-driver/mysql 控制
	logger.Debug("MySql synced and closed")
	return nil
}

func formatPattern(pattern []byte) string {
	return fmt.Sprintf("%x", pattern)
}

func unformatPattern(patterString string) []byte {
	var dst []byte
	akey := []byte(patterString)
	dst = make([]byte, hex.DecodedLen(len(akey)))
	hex.Decode(dst, akey)
	return dst
}

// AddPattern add
func (m *MySQL) AddPattern(pattern []byte, context sql.QueryContext) error {
	// pattern := sql.Pattern(context.Query)
	patternString := formatPattern(pattern)

	atomic.AddUint64(&QueryCounter, 1)
	o := orm.NewOrm()
	exist := o.QueryTable("pattern").Filter("key", patternString).Filter("uuid", m.UUID).Exist()
	if !exist {
		var aPattern Pattern
		aPattern.Key = patternString
		aPattern.Value = string(context.Query)
		aPattern.UUID = m.UUID
		aPattern.Enable = true
		id, err := o.Insert(&aPattern)
		if err == nil {
			logger.Debugf("Pattern saved, ID: %d", id)
		} else {
			logger.Warningf("Pattern saved error: %s", err.Error())
		}
	}
	uKey := bytes.Buffer{}
	uKey.Write(pattern)
	uKey.WriteString("_user_")
	uKey.Write(context.User)
	uKeyString := formatPattern(uKey.Bytes())

	exist = o.QueryTable("pattern").Filter("key", uKeyString).Filter("uuid", m.UUID).Exist()
	if !exist {
		var aPattern Pattern
		aPattern.Key = uKeyString
		aPattern.Value = formatPattern([]byte{0x11})
		aPattern.UUID = m.UUID
		aPattern.Enable = true
		id, err := o.Insert(&aPattern)
		if err == nil {
			logger.Debugf("Pattern User saved, ID: %d", id)
		} else {
			logger.Warningf("Pattern User saved error: %s", err.Error())
		}
	}

	cKey := bytes.Buffer{}
	cKey.Write(pattern)
	cKey.WriteString("_client_")
	cKey.Write(context.Client)
	cKeyString := formatPattern(cKey.Bytes())

	exist = o.QueryTable("pattern").Filter("key", cKeyString).Filter("uuid", m.UUID).Exist()
	if !exist {
		var aPattern Pattern
		aPattern.Key = cKeyString
		aPattern.Value = formatPattern([]byte{0x11})
		aPattern.UUID = m.UUID
		aPattern.Enable = true
		id, err := o.Insert(&aPattern)
		if err == nil {
			logger.Debugf("Pattern Source saved, ID: %d", id)
		} else {
			logger.Warningf("Pattern Source saved error: %s", err.Error())
		}
	}

	return nil
}

//CheckQuery check query
func (m *MySQL) CheckQuery(context sql.QueryContext, checkUser bool, checkSource bool) bool {
	atomic.AddUint64(&QueryCounter, 1)
	pattern := sql.Pattern(context.Query)
	patternString := formatPattern(pattern)
	o := orm.NewOrm()
	exist := o.QueryTable("pattern").Filter("key", patternString).Filter("enable", true).Filter("uuid", m.UUID).Exist()
	if !exist {
		return false
	}
	key := bytes.Buffer{}
	if checkUser {
		key.Write(pattern)
		key.WriteString("_user_")
		key.Write(context.User)
		exist := o.QueryTable("pattern").Filter("key", formatPattern(key.Bytes())).Filter("enable", true).Filter("uuid", m.UUID).Exist()
		if !exist {
			return false
		}
	}
	if checkSource {
		key.Reset()
		key.Write(pattern)
		key.WriteString("_client_")
		key.Write(context.Client)
		exist := o.QueryTable("pattern").Filter("key", formatPattern(key.Bytes())).Filter("enable", true).Filter("uuid", m.UUID).Exist()
		if !exist {
			return false
		}
	}
	return true
}

//CheckPermission check if has permission
func (m *MySQL) CheckPermission(context sql.QueryContext, q bool, v bool) bool {
	// get statement type
	stmt := sql.GetStmtType(string(context.Query))
	if stmt == sql.StmtUnknown {
		return false
	}
	tables, _ := sql.ExtractTableNames(string(context.Query))
	logger.Debugf("tables: ", tables)
	// verify permission
	o := orm.NewOrm()
	qs := o.QueryTable("permission")
	var exist bool
	if exist = qs.Filter("uuid", m.UUID).Exist(); !exist {
		// 没有规则
		logger.Debug("no rules")
		return true
	}
	if len(tables) > 0 {
		exist = qs.
			Filter("db", string(context.Database)).
			Filter("user", string(context.User)).
			// Filter("client", string(context.Client)).
			Filter("permission__contains", stmt).
			Filter("table__in", tables).
			Filter("enable", true).
			Filter("uuid", m.UUID).
			Exist()
	} else {
		exist = qs.
			Filter("db", string(context.Database)).
			Filter("user", string(context.User)).
			// Filter("client", string(context.Client)).
			Filter("permission__contains", stmt).
			Filter("enable", true).
			Filter("uuid", m.UUID).Exist()
	}
	if !exist {
		return false
	}
	return true
}

//UpdateState update
func (m *MySQL) UpdateState() error {
	o := orm.NewOrm()
	var state State
	err := o.QueryTable("state").Filter("key", "state").Filter("uuid", m.UUID).One(&state)
	if err != nil {
		if err == orm.ErrMultiRows {
			// 多条的时候报错
			logger.Warning("Returned Multi Rows Not One")
		}
		if err == orm.ErrNoRows {
			// 没有找到记录
			logger.Warning("Not row found")
			var newState State
			newState.QueryCounter = QueryCounter
			newState.QueryCounter = AbnormalCounter
			newState.Key = "state"
			newState.UUID = m.UUID
			id, err := o.Insert(&newState)
			if err == nil {
				logger.Warning(id)
				return nil
			}
			return err
		}
		return err
	}
	state.QueryCounter = QueryCounter
	state.AbnormalCounter = AbnormalCounter
	_, err = o.Update(&state)
	if err == nil {
		logger.Debugf("State Updated, QueryCounter:%d AbnormalCounter:%d", QueryCounter, AbnormalCounter)
		return nil
	}
	return err
}

// Abnormals list abnormals
func (m *MySQL) Abnormals() (count int) {
	var abnormals []*QueryAction
	o := orm.NewOrm()
	_, err := o.QueryTable("query_action").Filter("is_abnormal", true).Filter("uuid", m.UUID).All(&abnormals)
	if err == nil && len(abnormals) > 0 {
		logger.Debug("range abnormal")
		for _, element := range abnormals {
			// var c sql.QueryContext
			// c.Unmarshal(unformatPattern(element.Value))
			fmt.Printf("[%s] [User: %s] [CliendIP: %s] [Database: %s] [AbnormalType: %s] %s\n",
				element.Time.Format(time.RFC1123),
				element.User,
				element.ClientIP,
				element.Database,
				element.AbnormalType,
				element.Query)
			count++
		}
	} else {
		logger.Debug("no abnormals")
	}
	return
}

// Patterns list Patterns
func (m *MySQL) Patterns() (count int) {
	logger.Debugf("==> Patterns")
	var patterns []*Pattern
	o := orm.NewOrm()
	_, err := o.QueryTable("pattern").Filter("uuid", m.UUID).All(&patterns)
	if err == nil {
		logger.Debug(patterns)
		for _, element := range patterns {
			elementKey := unformatPattern(element.Key)
			if strings.Index(string(elementKey), "_client_") == -1 && strings.Index(string(elementKey), "_user_") == -1 {
				fmt.Printf(
					`-----Pattern: 0x%s
Sample: %s
`,
					element.Key,
					element.Value,
				)
				count++
			}
		}
	} else {
		logger.Warningf("Pattern error: %s", err.Error())
	}
	return
}

//InitialDB local databases
func (m *MySQL) InitialDB(str string, syncInterval time.Duration, timeout time.Duration) error {
	orm.Debug = false
	//InitLocalDB initail local db
	orm.RegisterDriver("mysql", orm.DRMySQL)

	err := orm.RegisterDataBase("default", "mysql", str, 30)
	if err != nil {
		// logger.Debugf("%s", err.Error())
		return err
	}
	// 注册定义的model
	orm.RegisterModel(new(QueryAction))
	orm.RegisterModel(new(Pattern))
	// orm.RegisterModel(new(Abnormal))
	orm.RegisterModel(new(State))
	orm.RegisterModel(new(Permission))

	// 创建table
	// Database alias.
	name := "default"
	// Drop table and re-create.
	force := false
	// Print log.
	verbose := false
	orm.RunSyncdb(name, force, verbose)
	return nil
}
