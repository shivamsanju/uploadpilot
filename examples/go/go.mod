module test

go 1.23.2

replace github.com/uploadpilot/sdk/go-client v0.0.0 => ../sdk/go-client

require github.com/uploadpilot/sdk/go-client v0.0.0

require (
	github.com/eventials/go-tus v0.0.0-20220610120217-05d0564bb571 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
