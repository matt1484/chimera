module github.com/matt1484/chimera

go 1.19

retract (
	[v0.0.0, v0.0.5] // pre-release versions
	v0.1.0 // middleware bug
)

require (
	github.com/go-chi/chi/v5 v5.0.8
	github.com/go-playground/form/v4 v4.2.1
	github.com/invopop/jsonschema v0.12.0
	github.com/matt1484/spectagular v1.0.4
	github.com/stretchr/testify v1.8.1
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/swaggest/swgui v1.8.0 // indirect
	github.com/vearutop/statigz v1.4.0 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
