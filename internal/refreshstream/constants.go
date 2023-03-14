package refreshstream

const (
	QueryStateTrue = `SELECT "id", "login", "pass", "ip_address_out", "cam_path",
	"code_mp", "state_public", "status_public"
	FROM public."%s"
	WHERE "cam_path" IS NOT null AND "state_public" = 1`

	QueryStateFalse = `SELECT "id", "login", "pass", "ip_address_out", "cam_path",
	"code_mp", "state_public", "status_public"
	FROM public."%s"
	WHERE "cam_path" IS NOT null AND "state_public" = 0`
)
