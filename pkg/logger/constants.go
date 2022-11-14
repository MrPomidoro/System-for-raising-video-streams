package logger

const (
	FileNameConst                = "out.log"
	ServTimestampFormatConst     = "2006-01-02 15:04:05"
	ServLogFormatConst           = "[%lvl%]: (%time%) %msg%\n"
	ServLogFormatStatusCodeConst = "[%lvl%]: (%time%) {Request Method: %meth%. Status code - %status%} - %msg%\n"
) //
