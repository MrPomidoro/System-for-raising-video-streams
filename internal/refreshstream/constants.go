package refreshstream

// queries
// const (
// 	SELECT_COL_FROM_TBL               = "SELECT %s FROM %s"
// 	SELECT_COL_FROM_TBL_WHERE_CND     = "SELECT %s FROM %s WHERE %s"
// 	INSERT_INTO_TBL_VALUES_VAL        = "INSERT INTO %s(%s) VALUES (%s) ON CONFLICT DO NOTHING"
// 	UPDATE_TBL_SET_VAL_WHERE_CND      = `UPDATE public."refresh_stream" SET %v WHERE %v`
// 	DELETE_FROM_TBL_WHERE_CND         = "DELETE FROM %s WHERE %s"
// 	DELETE_CASCADE_FROM_TBL_WHERE_CND = "DELETE CASCADE FROM %s WHERE %s"
// )

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
