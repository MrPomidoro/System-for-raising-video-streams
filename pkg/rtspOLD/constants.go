package rtsp

const (
	// 1 - cfg.Run из конфига, 2 - port_srv, 3 - sp, 4 - cam_id из бд
	RunOnReadyConst = "%s localhost --port %s --stream_path %s --camera_id %s"

	// 1 - host, 2 - port из конфига
	URLGetConst = "http://%s%s/v1/paths/list"

	// 1 - host, 2 - port из конфига, 3 - стрим камеры из базы/rtsp
	URLPostAddConst = "http://%s%s/v1/config/paths/add/%s"

	// 1 - host, 2 - port из конфига, 3 - стрим камеры из базы/rtsp
	URLPostRemoveConst = "http://%s%s/v1/config/paths/remove/%s"

	// 1 - host, 2 - port из конфига, 3 - стрим камеры из базы/rtsp
	URLPostEditConst = "http://%s%s/v1/config/paths/edit/%s"
)
