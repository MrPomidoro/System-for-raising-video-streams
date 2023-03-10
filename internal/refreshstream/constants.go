package refreshstream

const (
	QueryStateTrue = `SELECT *
	FROM public."refresh_stream"
	WHERE "stream" IS NOT null AND "stream_state" = true`

	QueryStateFalse = `SELECT *
	FROM public."refresh_stream"
	WHERE "stream" IS NOT null AND "stream_state" = false`

	QueryEditStatus = `UPDATE public."refresh_stream"
	SET "stream_status"='true'
	WHERE "stream"='%s'`
)
