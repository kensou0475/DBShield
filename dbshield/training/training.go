package training

import (
	"errors"

	"github.com/qiwihui/DBShield/dbshield/config"
	"github.com/qiwihui/DBShield/dbshield/sql"
)

var (
	errInvalidParrent = errors.New("Invalid pattern")
	errInvalidUser    = errors.New("Invalid user")
	errInvalidClient  = errors.New("Invalid client")
)

//AddToTrainingSet records query context in local database
func AddToTrainingSet(context sql.QueryContext) error {
	pattern := sql.Pattern(context.Query)
	return config.Config.LocalDB.AddPattern(pattern, context)
}

//CheckQuery pattern, returns true if it finds the pattern
//We should keep it as fast as possible
func CheckQuery(context sql.QueryContext) bool {
	hasQuery := config.Config.LocalDB.CheckQuery(context)
	if !hasQuery {
		// pattern := sql.Pattern(context.Query)
		config.Config.LocalDB.RecordAbnormal(context)
		return false
	}
	return true
}
