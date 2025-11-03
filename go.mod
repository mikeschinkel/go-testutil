module github.com/mikeschinkel/go-testutil

go 1.25.3

require github.com/mikeschinkel/go-cliutil v0.0.0-20251103044509-f4857eae9f54

require (
	github.com/mikeschinkel/go-dt v0.0.0-20251103073248-cc1248280ed9
	github.com/mikeschinkel/go-dt/appinfo v0.0.0-20251103073248-cc1248280ed9 // indirect
	github.com/mikeschinkel/go-dt/de v0.0.0-20251103073248-cc1248280ed9 // indirect
)

replace (
	github.com/mikeschinkel/go-cfgstore => ../go-cfgstore

	github.com/mikeschinkel/go-cliutil => ../go-cliutil

	github.com/mikeschinkel/go-dt => ../go-dt
	github.com/mikeschinkel/go-dt/appinfo => ../go-dt/appinfo
	github.com/mikeschinkel/go-dt/de => ../go-dt/de
	github.com/mikeschinkel/go-dt/dtx => ../go-dt/dtx
)
