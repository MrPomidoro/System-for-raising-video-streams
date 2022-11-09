package refreshstream

import "database/sql"

// структура таблицы refresh_stream
// sql.Null* когда возможен null в столбце
type RefreshStream struct {
	Id            int            `json:"id" db:"id"`
	Auth          sql.NullString `json:"auth" db:"auth"`
	Ip            sql.NullString `json:"ip" db:"ip"`
	Stream        sql.NullString `json:"stream" db:"stream"`
	Portsrv       string         `json:"portsrv" db:"portsrv"`
	Sp            sql.NullString `json:"sp" db:"sp"`
	Camid         sql.NullString `json:"camid" db:"camid"`
	Record_status sql.NullBool   `json:"record_status" db:"record_status"`
	Stream_status sql.NullBool   `json:"stream_status" db:"stream_status"`
	Record_state  sql.NullBool   `json:"record_state" db:"record_state"`
	Stream_state  sql.NullBool   `json:"stream_state" db:"stream_state"`
}
