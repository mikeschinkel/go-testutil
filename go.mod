module github.com/mikeschinkel/go-testutil

go 1.25.3

require github.com/mikeschinkel/go-cliutil v0.0.0-20251027221228-4ec65dd988d2

require (
	github.com/mikeschinkel/go-cfgstore v0.0.0-20251027170244-67974a842019 // indirect
	github.com/mikeschinkel/go-dt v0.0.0-20251027222746-b5ea4e0da9da
	github.com/mikeschinkel/go-dt/appinfo v0.0.0-00010101000000-000000000000 // indirect
	github.com/mikeschinkel/go-dt/de v0.0.0-00010101000000-000000000000 // indirect
	github.com/mikeschinkel/go-fsfix v0.1.0 // indirect
)

replace (
	github.com/mikeschinkel/go-cfgstore => ../go-cfgstore

	github.com/mikeschinkel/go-cliutil => ../go-cliutil

	github.com/mikeschinkel/go-dt => ../go-dt
	github.com/mikeschinkel/go-dt/appinfo => ../go-dt/appinfo
	github.com/mikeschinkel/go-dt/de => ../go-dt/de
	github.com/mikeschinkel/go-dt/dtx => ../go-dt/dtx
)
