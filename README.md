# go-spatialite

Use SQLite3 with the Spatialite extension in Go.

## Use

This registers a `database/sql` driver called `spatialite` by loading the
`spatialite` extension with the the `github.com/mattn/go-sqlite3` module.

To use in your application:

1. Because this uses CGO, you must have `CGO_ENABLED=1`.
2. Make sure you have `libsqlite3` and `libspatialite` installed on your system.
   On some systems (MacOS, for example), `libspatialite` is referred to as
   `modspatialite`.
3. Make sure that the module from the point step above is installed in a place your 
   system can find it. Where that is depends on your OS:
    1. `linux`
       - `$HOMEBREW_PREFIX/lib` (if it exists)
       - the directories listed in `/etc/ld.so.conf`
    2. `darwin`:
       - `$HOMEBREW_PREFIX/lib` (if it exists)
       - If Xcode is installed, `/Applications/Xcode.app/Contents`
       - If the go executable is in a path that ends with `Contents/bin/`, 
         which implies that it's part of a MacOS `.app` file, then the
         `Contents` directory will be searched.
    3. `windows`:
       - `$env:PATH`
   
   If you want to override this, set the environment variable
   `GO_SPATIALITE_SEARCH_PATH` as a set of directories separated by your OS's
   path separator character—`;` for `windows`, `:` for `linux` and `darwin`.
4. Make sure that `github.com/watchedsky-social/go-spatialite` is imported in
   your app, because the driver is registered via an `init` function.

## Credits

This was taken from [`github.com/shaxbee/go-spatialite`](https://github.com/shaxbee/go-spatialite)
for the registration portion. Thanks!