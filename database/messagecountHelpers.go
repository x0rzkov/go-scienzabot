package database

import (
	"database/sql"
	"errors"
)

/*
CREATE TABLE IF NOT EXISTS 'MessageCount' (
	'ID'  INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'UserID'  INTEGER NOT NULL,
	'GroupID'  INTEGER NOT NULL,
	'MessageCount'  INTEGER NOT NULL,
	FOREIGN KEY('UserID') REFERENCES Users('ID'),
	FOREIGN KEY('GroupID') REFERENCES Groups('ID'),
	CONSTRAINT con_msgcoubt_user_group_unique UNIQUE ('UserID','GroupID')
);
*/

//GetMessageCount returns the message number of a user in a group
func (db *SQLiteDB) GetMessageCount(user int64, group int64) (int64, error) {
	var messageCount int64
	err := db.QueryRow("SELECT MessageCount FROM MessageCount WHERE `UserID` AND `GroupID`", user, group).Scan(&messageCount)
	switch {
	case err == sql.ErrNoRows:
		db.AddLogEvent(Log{Event: "GetMessageCount_ErrorNoRows", Message: "Impossible to get rows", Error: err.Error()})
		return messageCount, err
	case err != nil:
		db.AddLogEvent(Log{Event: "GetMessageCount_ErrorUnknown", Message: "Uknown error verified", Error: err.Error()})
		return messageCount, err
	default:
		return messageCount, nil
	}
}

//SetMessageCount sets the message of a user in a group
func (db *SQLiteDB) SetMessageCount(user int64, group int64, messageCount int64) error {
	query, err := db.Exec(
		"INSERT INTO MessageCount (`UserID`, `GroupID`, `MessageCount`) VALUES (?,?,?) " +
			"ON CONFLICT(`UserID`, `GroupID`) DO UPDATE SET `MessageCount` = Excluded.MessageCount",
	)
	if err != nil {
		db.AddLogEvent(Log{Event: "SetMessageCount_QueryFailed", Message: "Impossible to create the execute the query", Error: err.Error()})
		return err
	}
	rows, err := query.RowsAffected()
	if err != nil {
		db.AddLogEvent(Log{Event: "SetMessageCount_RowsInfoNotGot", Message: "Impossible to get afftected rows", Error: err.Error()})
		return err
	}
	if rows < 1 {
		db.AddLogEvent(Log{Event: "SetMessageCount_NoRowsAffected", Message: "No rows affected", Error: err.Error()})
		return NoRowsAffected{error: errors.New("No rows affected from the query")}
	}
	return err
}

//IncrementMessageCount increments by 1 the number of messages from a user
func (db *SQLiteDB) IncrementMessageCount(user int64, group int64) error {
	msgCnt, err := db.GetMessageCount(user, group)
	/*if err != nil {
		//User may not exist yet in DB
		db.AddLogEvent(Log{Event: "IncrementMessageCount_CannotGetMessageCount", Message: "user may not exist", Error: err.Error()})
	}*/
	err = db.SetMessageCount(user, group, msgCnt+1)
	return err
}
