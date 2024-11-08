package customError

var (
	ErrorConfig = &Error{
		level: FatalLevel,
		code:  "50.0.1",
		desc:  "error at the level of reading and processing the config",
		err:   nil,
	}

	ErrorDatabase = &Error{
		level: FatalLevel,
		code:  "50.1.1",
		desc:  "error at database operation level",
		err:   nil,
	}

	ErrorApp = &Error{
		level: FatalLevel,
		code:  "50.3.1",
		desc:  "error at app operation level",
		err:   nil,
	}

	ErrorRefreshStream = &Error{
		level: ErrorLevel,
		code:  "50.1.2",
		desc:  "refresh stream entity error at database operation level",
		err:   nil,
	}

	ErrorStatusStream = &Error{
		level: ErrorLevel,
		code:  "50.1.3",
		desc:  "status stream entity error at database operation level",
		err:   nil,
	}

	ErrorRTSP = &Error{
		level: ErrorLevel,
		code:  "50.2.1",
		desc:  "error at rtsp-simple-server entity operation level",
		err:   nil,
	}
)
