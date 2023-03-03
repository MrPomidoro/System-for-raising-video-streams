package refreshstream

import "database/sql"

// Stream Структура таблицы refresh_stream.
// sql.Null* когда возможен null в столбце
type Stream struct {
	Id           int            `json:"id" db:"id"`
	Auth         sql.NullString `json:"auth" db:"auth"`
	Ip           sql.NullString `json:"ip" db:"ip"`
	Stream       string         `json:"stream" db:"stream"`
	Portsrv      string         `json:"portsrv" db:"portsrv"`
	Sp           sql.NullString `json:"sp" db:"sp"`
	CamId        sql.NullString `json:"camid" db:"camid"`
	RecordStatus sql.NullBool   `json:"record_status" db:"record_status"`
	StreamStatus sql.NullBool   `json:"stream_status" db:"stream_status"`
	RecordState  sql.NullBool   `json:"record_state" db:"record_state"`
	StreamState  sql.NullBool   `json:"stream_state" db:"stream_state"`
	Protocol     sql.NullString `json:"protocol" db:"protocol"`
}
