package migrations

import "gofr.dev/pkg/gofr/migration"

func All() map[int64]migration.Migrate {
	return map[int64]migration.Migrate{
		// 20241013015640: createTableUser(),
		20241013015650: createTableTask(),
		// 20241013015656: createTableSession(),
	}
}
