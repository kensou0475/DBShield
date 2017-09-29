package db

import (
	"encoding/binary"
	"net"
	"time"
	// mysql orm
	"github.com/astaxie/beego/orm"
	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/nim4/DBShield/dbshield/logger"
	"github.com/nim4/DBShield/dbshield/sql"
)

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

func fourByteBigEndianToIP(data []byte) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(data))
	return ip.String()
}

// RecordQueryAction record query and action
func RecordQueryAction(context sql.QueryContext, action string) error {
	logger.Debugf("action: %s", action)

	// 异步记录
	go func() {
		o := orm.NewOrm()
		var queryAction QueryAction
		queryAction.Query = string(context.Query)
		queryAction.User = string(context.User)
		queryAction.Client = fourByteBigEndianToIP(context.Client)
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
