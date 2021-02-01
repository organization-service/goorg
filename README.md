# CEAP(シープ)

Connect enterprise with agentur and people

企業とフリーランスや労働者を結ぶサービスで使用しているフレームワーク

## Contents

- [CEAP(シープ)](#ceapシープ)
  - [Contents](#contents)
  - [Environment](#environment)

## Environment

| key                      | overview                                   |
| :----------------------- | :----------------------------------------- |
| DB_CONFIG_DIR            | dbconfig.ymlの存在しているディレクトリパス |
| APM_NAME                 | 使用しているAPMの名称(elastic or newrelic) |
| IDP_NAME                 | IDP名(auth0)                               |
| AUTH0_DOMAIN             | Auth0のドメイン名                          |
| AUTH0_AUD                | Auth0で設定しているaudience                |
| AUTH0_ISS                | Auth0で設定されるIssuer                    |
| NEW_RELIC_APP_NAME       | newrelic appliacation name                 |
| NEW_RELIC_LICENSE_KEY    | newrelic license key                       |
| ELASTIC_APM_SERVER_URL   | elastic apm server URL                     |
| ELASTIC_APM_SECRET_TOKEN | elastic apm secret token                   |
| ELASTIC_APM_SERVICE_NAME | elastic service name                       |

[elastic apm environment variables](https://www.elastic.co/guide/en/apm/agent/go/current/configuration.html#configuration)

[newrelic environment variables](https://docs.newrelic.com/docs/agents/go-agent/configuration/go-agent-configuration)