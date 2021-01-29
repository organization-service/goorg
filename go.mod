module github.com/organization-service/goorg

go 1.15

require (
	github.com/julienschmidt/httprouter v1.3.0
	github.com/newrelic/go-agent/v3 v3.10.0
	github.com/newrelic/go-agent/v3/integrations/nrhttprouter v1.0.0
	github.com/stretchr/testify v1.5.1
	go.elastic.co/apm/module/apmgormv2 v1.10.0
	go.elastic.co/apm/module/apmhttprouter v1.10.0
	gopkg.in/yaml.v2 v2.2.2
	gorm.io/driver/mysql v1.0.2
	gorm.io/driver/postgres v1.0.7
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.20.12
)
