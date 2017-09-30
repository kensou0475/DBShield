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
		config.Config.LocalDB.RecordAbnormal(context)
		return false
	}
	return true
}
