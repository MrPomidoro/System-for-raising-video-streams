package rtspsimpleserver

const (
	// 1 - url из конфига
	URLGetConst = "%s/v1/paths/list"

	// 1 - url из конфига, 2 - выполняемое действие (add, remove, edit), 3 - стрим камеры из базы/rtsp
	URLPostConst = "%s/v1/config/paths/%s/%s"
)
