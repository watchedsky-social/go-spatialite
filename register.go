package spatialite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	"github.com/watchedsky-social/go-spatialite/sharedlibs"
)

var (
	libNames    = []string{"mod_spatialite", "libspatialite"}
	symbolNames = []string{"sqlite3_modspatialite_init", "spatialite_init_ex"}
)

var ErrSpatialiteNotFound = errors.New("spatialite: spatialite extension not found")

func init() {
	sql.Register("spatialite", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			var errs []error
			for _, lib := range libNames {
				absPath, err := sharedlibs.FindSharedLibraryFiles(lib, true)
				if err == nil && len(absPath) > 0 {
					for _, symbol := range symbolNames {
						if err = conn.LoadExtension(absPath[0], symbol); err == nil {
							return nil
						}

						errs = append(errs, err)
					}
				}
			}
			return fmt.Errorf("%w: %w", ErrSpatialiteNotFound, errors.Join(errs...))
		},
	})
}
