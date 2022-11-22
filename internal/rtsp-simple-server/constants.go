package rtspsimpleserver

const (
	// 1 - host, 2 - port из конфига
	URLGetConst = "http://%s%s/v1/paths/list"

	// 1 - host, 2 - port из конфига, 3 - выполняемое действие (add, remove, edit), 4 - стрим камеры из базы/rtsp
	URLPostConst = "http://%s%s/v1/config/paths/%s/%s"

	// // 1 - host, 2 - port из конфига, 3 - стрим камеры из базы/rtsp
	// URLPostRemoveConst = "http://%s%s/v1/config/paths/remove/%s"

	// // 1 - host, 2 - port из конфига, 3 - стрим камеры из базы/rtsp
	// URLPostEditConst = "http://%s%s/v1/config/paths/edit/%s"
)
