package statusstream

import (
	"time"
)

// Структура таблицы refresh_stream.
// sql.Null* когда возможен null в столбце
type StatusStream struct {
	StreamId       int       `json:"stream_id" db:"stream_id"`
	DateCreate     time.Time `json:"date_create" db:"date_create"`
	StatusResponse bool      `json:"status_response" db:"status_response"`
}
