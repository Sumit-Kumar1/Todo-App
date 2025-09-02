package migrations

// nolint:gochecknoglobals // required this as a global but is not exported
var migrations = map[string]migrator{
	"20241013015650": M20241013015650(""),
}
