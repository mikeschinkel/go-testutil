module github.com/mikeschinkel/go-testutil

go 1.25.3

require (
	github.com/mikeschinkel/go-cfgstore v0.0.0-20251027170244-67974a842019
	github.com/mikeschinkel/go-cliutil v0.0.0-20251027170801-82399064d27f
	github.com/mikeschinkel/go-fsfix v0.1.0
)

require github.com/mikeschinkel/go-dt v0.0.0-20251027170931-0f47f0479185 // indirect

replace (
	github.com/mikeschinkel/go-cfgstore => ../go-cfgstore

	github.com/mikeschinkel/go-cliutil => ../go-cliutil

	github.com/mikeschinkel/go-dt => ../go-dt
)
