module github.com/organization-service/goorg/v2

go 1.15

require (
	github.com/auth0/go-jwt-middleware v1.0.0
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible
	github.com/julienschmidt/httprouter v1.3.0
	github.com/newrelic/go-agent v3.10.0+incompatible
	github.com/newrelic/go-agent/v3 v3.10.0
	github.com/newrelic/go-agent/v3/integrations/nrhttprouter v1.0.0
	github.com/rakyll/statik v0.1.7
	github.com/savaki/swag v0.0.0-20170722173931-3a75479e44a3
	github.com/stretchr/testify v1.5.1
	github.com/urfave/cli/v2 v2.3.0
	go.elastic.co/apm v1.10.0
	go.elastic.co/apm/module/apmgormv2 v1.10.0
	go.elastic.co/apm/module/apmhttp v1.10.0
	go.elastic.co/apm/module/apmhttprouter v1.10.0
	go.uber.org/dig v1.10.0
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/tools v0.0.0-20200509030707-2212a7e161a5
	gopkg.in/yaml.v2 v2.2.3
	gorm.io/driver/mysql v1.0.2
	gorm.io/driver/postgres v1.0.7
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.12
)
