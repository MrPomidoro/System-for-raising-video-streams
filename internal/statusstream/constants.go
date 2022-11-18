package statusstream

const (
	InsertToStatusStream = `INSERT INTO public."status_stream" ("stream_id", "date_create", "status_response")
	VALUES (%d, current_timestamp, %t)`
)
