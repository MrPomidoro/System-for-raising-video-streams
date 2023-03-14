package refreshstream

import (
	"database/sql"

	"github.com/jackc/pgtype"
)

// Stream -  Структура таблицы refresh_stream.
// sql.Null* когда возможен null в столбце
type Stream struct {
	Id           int            `json:"id" db:"id"`
	Login        sql.NullString `json:"login" db:"login"`
	Pass         sql.NullString `json:"pass" db:"pass"`
	Ip           pgtype.Inet    `json:"ip_address_out" db:"ip_address_out"`
	CamPath      sql.NullString `json:"cam_path" db:"cam_path"`
	StatePublic  sql.NullInt16  `json:"state_public" db:"state_public"`
	StatusPublic sql.NullInt16  `json:"status_public" db:"status_public"`
	CodeMp       string         `json:"code_mp" db:"code_mp"`
	Port         string         // 554
	Protocol     sql.NullString `json:"protocol" db:"protocol"`

	// Stream       string         `json:"stream" db:"stream"`
	// CamId       sql.NullString `json:"camid" db:"camid"`
}
