package database

const (
	DBInfoConst         = "host=%s port=%s user=%s password=%s dbname=%s"
	OpenDBErrConst      = "cannot get connect to database: %s"
	ConnectToDBOkConst  = "Success connect to database %s"
	ConnectToDBErrConst = "cannot connect to database %s"
	WaitForBDErrConst   = "Time waiting of DB connection exceeded limit: %v"
	CloseDBErrConst     = "cannot close DB connection. Error: %v"
	CloseDBOkConst      = "Established closing of connection to DB"
)
