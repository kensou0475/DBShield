package training

import (
	"github.com/qiwihui/DBShield/dbshield/config"
	"github.com/qiwihui/DBShield/dbshield/sql"
)

//AddToTrainingSet records query context in local database
func AddToTrainingSet(context sql.QueryContext) error {
	pattern := sql.Pattern(context.Query)
	return config.Config.LocalDB.AddPattern(pattern, context)
}

//CheckQuery pattern, returns true if it finds the pattern
//We should keep it as fast as possible
func CheckQuery(context sql.QueryContext) bool {
	hasQuery := config.Config.LocalDB.CheckQuery(context, config.Config.CheckUser, config.Config.CheckSource)
	if !hasQuery {
		// pattern := sql.Pattern(context.Query)
		config.Config.LocalDB.RecordAbnormal(context, "pattern")
		return false
	}
	return true
}

// CheckPermission if has statement permission
func CheckPermission(context sql.QueryContext) bool {
	hasPermission := config.Config.LocalDB.CheckPermission(context, true, true)
	if !hasPermission {
		// pattern := sql.Pattern(context.Query)
		config.Config.LocalDB.RecordAbnormal(context, "permission")
		return false
	}
	return true
}
