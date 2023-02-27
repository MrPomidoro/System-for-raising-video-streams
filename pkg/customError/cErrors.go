package customError

var (
	ErrorConfig = &Error{
		level: FatalLevel,
		code:  "50.0.1",
		desc:  "error at the level of reading and processing the config",
		err:   nil,
		deep:  nil,
	}

	ErrorDatabase = &Error{
		level: FatalLevel,
		code:  "50.1.1",
		desc:  "error at database operation level",
		err:   nil,
		deep:  nil,
	}

	ErrorStorage = &Error{
		level: FatalLevel,
		code:  "50.2.1",
		desc:  "error at database operation level",
		err:   nil,
		deep:  nil,
	}

	ErrorRefreshStream = &Error{
		level: ErrorLevel,
		code:  "50.2.2",
		desc:  "refresh stream entity error at database operation level",
		err:   nil,
		deep:  nil,
	}

	ErrorStatusStream = &Error{
		level: ErrorLevel,
		code:  "50.2.3",
		desc:  "status stream entity error at database operation level",
		err:   nil,
		deep:  nil,
	}

	ErrorRTSP = &Error{
		level: ErrorLevel,
		code:  "50.3.1",
		desc:  "rtsp-simple-server entity error at database operation level",
		err:   nil,
		deep:  nil,
	}
)
