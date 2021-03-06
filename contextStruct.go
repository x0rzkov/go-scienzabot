package main

import (
	"scienzabot/database"
	"scienzabot/embtypes"
)

// Context is a struct that contains the current redis client, the bot and a update.
// It is used as a paramater, passed withafter the routing
type Context struct {
	Bot      *embtypes.Tgbotapi
	Update   *embtypes.Tgupdate
	Database *database.SQLiteDB
}
