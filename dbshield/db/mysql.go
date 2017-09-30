package db

import (
	"encoding/binary"
	"errors"
	"net"
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
}

//QueryAction query and action
type QueryAction struct {
	ID     int
	Query  string    `orm:"column(query);null;type(text)"`
	User   string    `orm:"column(user);null;size(128)"`
	Client string    `orm:"column(client);null;size(128)"`
	Db     string    `orm:"column(db);null;size(128)"`
	Time   time.Time `orm:"column(time);type(datetime);size(6)"`
	Action string    `orm:"column(action);size(32)"`
}

//Pattern record trainging set
type Pattern struct {
	ID    int
	Key   string `orm:"column(key);null;type(text)"`
	Value string `orm:"column(value);null;type(text)"`
}

//Abnormal record abnormal set
type Abnormal struct {
	ID    int
	Key   string `orm:"column(key);type(text)"`
	Value string `orm:"column(value);type(text)"`
}

//State record abnormal set
type State struct {
	ID              int
	QueryCounter    uint64 `orm:"column(QueryCounter);type(bigint unsigned)"`
	AbnormalCounter uint64 `orm:"column(AbnormalCounter);type(bigint unsigned)"`
}

func oneFourByteBigEndianToIP(data []byte) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(data))
	return ip.String()
}

// RecordQueryAction record query and action
func (m *MySQL) RecordQueryAction(context sql.QueryContext, action string) error {
	logger.Debugf("action: %s", action)

	// 异步记录
	go func() {
		o := orm.NewOrm()
		var queryAction QueryAction
		queryAction.Query = string(context.Query)
		queryAction.User = string(context.User)
		queryAction.Client = oneFourByteBigEndianToIP(context.Client)
		queryAction.Db = string(context.Database)
		queryAction.Time = context.Time
		queryAction.Action = action
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
func (m *MySQL) RecordAbnormal(context sql.QueryContext) error {
	return errors.New("Not Impletement")
}

// CheckPattern check if pattern exist
func (m *MySQL) CheckPattern(pattern []byte) error {

	return errors.New("Not Impletement")
}

// PutPattern put pattern
func (m *MySQL) PutPattern(pattern []byte, query []byte) error {

	return errors.New("Not Impletement")
}

// DeletePattern delete pattern
func (m *MySQL) DeletePattern(pattern []byte) error {
	return errors.New("Not Impletement")
}

// Purge local databases
func (m *MySQL) Purge() error {
	return errors.New("Not Impletement")
}

// SyncAndClose local databases
func (m *MySQL) SyncAndClose() error {
	return errors.New("Not Impletement")
}

// AddPattern add
func (m *MySQL) AddPattern(pattern []byte, query sql.QueryContext) error {
	return errors.New("Not Impletement")
}

//CheckQuery check query
func (m *MySQL) CheckQuery(context sql.QueryContext) bool {
	return false
}

//UpdateState update
func (m *MySQL) UpdateState() error {
	return errors.New("Not Impletement")
}

// Abnormals list abnormals
func (m *MySQL) Abnormals() (count int) {

	return 0
}

// Patterns list Patterns
func (m *MySQL) Patterns() (count int) {
	return 0
}

//InitialDB local databases
func (m *MySQL) InitialDB(str string, syncInterval time.Duration, timeout time.Duration) error {
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
	orm.RegisterModel(new(Abnormal))
	orm.RegisterModel(new(State))

	// 创建table
	// Database alias.
	name := "default"
	// Drop table and re-create.
	force := false
	// Print log.
	verbose := true
	orm.RunSyncdb(name, force, verbose)
	return nil
}
