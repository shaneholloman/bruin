package connection

import (
	"context"
	"maps"
	"os"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/bruin-data/bruin/pkg/adjust"
	"github.com/bruin-data/bruin/pkg/airtable"
	"github.com/bruin-data/bruin/pkg/applovin"
	"github.com/bruin-data/bruin/pkg/applovinmax"
	"github.com/bruin-data/bruin/pkg/appsflyer"
	"github.com/bruin-data/bruin/pkg/appstore"
	"github.com/bruin-data/bruin/pkg/asana"
	"github.com/bruin-data/bruin/pkg/athena"
	"github.com/bruin-data/bruin/pkg/attio"
	"github.com/bruin-data/bruin/pkg/bigquery"
	"github.com/bruin-data/bruin/pkg/chess"
	"github.com/bruin-data/bruin/pkg/clickhouse"
	"github.com/bruin-data/bruin/pkg/config"
	"github.com/bruin-data/bruin/pkg/databricks"
	"github.com/bruin-data/bruin/pkg/db2"
	duck "github.com/bruin-data/bruin/pkg/duckdb"
	"github.com/bruin-data/bruin/pkg/dynamodb"
	"github.com/bruin-data/bruin/pkg/elasticsearch"
	"github.com/bruin-data/bruin/pkg/emr_serverless"
	"github.com/bruin-data/bruin/pkg/facebookads"
	"github.com/bruin-data/bruin/pkg/frankfurter"
	"github.com/bruin-data/bruin/pkg/gcs"
	"github.com/bruin-data/bruin/pkg/github"
	"github.com/bruin-data/bruin/pkg/googleads"
	"github.com/bruin-data/bruin/pkg/googleanalytics"
	"github.com/bruin-data/bruin/pkg/gorgias"
	"github.com/bruin-data/bruin/pkg/gsheets"
	"github.com/bruin-data/bruin/pkg/hana"
	"github.com/bruin-data/bruin/pkg/hubspot"
	"github.com/bruin-data/bruin/pkg/isocpulse"
	"github.com/bruin-data/bruin/pkg/kafka"
	"github.com/bruin-data/bruin/pkg/kinesis"
	"github.com/bruin-data/bruin/pkg/klaviyo"
	"github.com/bruin-data/bruin/pkg/linkedinads"
	"github.com/bruin-data/bruin/pkg/mixpanel"
	"github.com/bruin-data/bruin/pkg/mongo"
	"github.com/bruin-data/bruin/pkg/mssql"
	"github.com/bruin-data/bruin/pkg/mysql"
	"github.com/bruin-data/bruin/pkg/notion"
	"github.com/bruin-data/bruin/pkg/oracle"
	"github.com/bruin-data/bruin/pkg/personio"
	"github.com/bruin-data/bruin/pkg/phantombuster"
	"github.com/bruin-data/bruin/pkg/pipedrive"
	"github.com/bruin-data/bruin/pkg/postgres"
	"github.com/bruin-data/bruin/pkg/quickbooks"
	"github.com/bruin-data/bruin/pkg/s3"
	"github.com/bruin-data/bruin/pkg/salesforce"
	"github.com/bruin-data/bruin/pkg/sftp"
	"github.com/bruin-data/bruin/pkg/shopify"
	"github.com/bruin-data/bruin/pkg/slack"
	"github.com/bruin-data/bruin/pkg/smartsheet"
	"github.com/bruin-data/bruin/pkg/snowflake"
	"github.com/bruin-data/bruin/pkg/solidgate"
	"github.com/bruin-data/bruin/pkg/spanner"
	"github.com/bruin-data/bruin/pkg/sqlite"
	"github.com/bruin-data/bruin/pkg/stripe"
	"github.com/bruin-data/bruin/pkg/tiktokads"
	"github.com/bruin-data/bruin/pkg/zendesk"
	"github.com/bruin-data/bruin/pkg/zoom"
	"github.com/pkg/errors"
	"github.com/sourcegraph/conc"
)

type Manager struct {
	BigQuery        map[string]*bigquery.Client
	Snowflake       map[string]*snowflake.DB
	Postgres        map[string]*postgres.Client
	MsSQL           map[string]*mssql.DB
	Databricks      map[string]*databricks.DB
	Mongo           map[string]*mongo.DB
	Mysql           map[string]*mysql.Client
	Notion          map[string]*notion.Client
	HANA            map[string]*hana.Client
	Shopify         map[string]*shopify.Client
	Gorgias         map[string]*gorgias.Client
	Klaviyo         map[string]*klaviyo.Client
	Adjust          map[string]*adjust.Client
	Athena          map[string]*athena.DB
	FacebookAds     map[string]*facebookads.Client
	Stripe          map[string]*stripe.Client
	Appsflyer       map[string]*appsflyer.Client
	Kafka           map[string]*kafka.Client
	Airtable        map[string]*airtable.Client
	DuckDB          map[string]*duck.Client
	Hubspot         map[string]*hubspot.Client
	GoogleSheets    map[string]*gsheets.Client
	Chess           map[string]*chess.Client
	S3              map[string]*s3.Client
	Slack           map[string]*slack.Client
	Asana           map[string]*asana.Client
	DynamoDB        map[string]*dynamodb.Client
	Zendesk         map[string]*zendesk.Client
	GoogleAds       map[string]*googleads.Client
	TikTokAds       map[string]*tiktokads.Client
	GitHub          map[string]*github.Client
	AppStore        map[string]*appstore.Client
	LinkedInAds     map[string]*linkedinads.Client
	ClickHouse      map[string]*clickhouse.Client
	GCS             map[string]*gcs.Client
	ApplovinMax     map[string]*applovinmax.Client
	Personio        map[string]*personio.Client
	Kinesis         map[string]*kinesis.Client
	Pipedrive       map[string]*pipedrive.Client
	Mixpanel        map[string]*mixpanel.Client
	QuickBooks      map[string]*quickbooks.Client
	Zoom            map[string]*zoom.Client
	Frankfurter     map[string]*frankfurter.Client
	EMRSeverless    map[string]*emr_serverless.Client
	GoogleAnalytics map[string]*googleanalytics.Client
	AppLovin        map[string]*applovin.Client
	Salesforce      map[string]*salesforce.Client
	SQLite          map[string]*sqlite.Client
	DB2             map[string]*db2.Client
	Oracle          map[string]*oracle.Client
	Phantombuster   map[string]*phantombuster.Client
	Elasticsearch   map[string]*elasticsearch.Client
	Spanner         map[string]*spanner.Client
	Solidgate       map[string]*solidgate.Client
	Smartsheet      map[string]*smartsheet.Client
	Attio           map[string]*attio.Client
	Sftp            map[string]*sftp.Client
	ISOCPulse       map[string]*isocpulse.Client
	mutex           sync.Mutex
}

func (m *Manager) GetConnection(name string) (interface{}, error) {
	availableConnectionNames := make([]string, 0)

	// todo(turtledev): make this DRY
	connBigQuery, err := m.GetBqConnectionWithoutDefault(name)
	if err == nil {
		return connBigQuery, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.BigQuery))...)

	connSnowflake, err := m.GetSfConnectionWithoutDefault(name)
	if err == nil {
		return connSnowflake, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Snowflake))...)

	connPostgres, err := m.GetPgConnectionWithoutDefault(name)
	if err == nil {
		return connPostgres, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Postgres))...)

	connMSSql, err := m.GetMsConnectionWithoutDefault(name)
	if err == nil {
		return connMSSql, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.MsSQL))...)

	connDatabricks, err := m.GetDatabricksConnectionWithoutDefault(name)
	if err == nil {
		return connDatabricks, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Databricks))...)

	connMongo, err := m.GetMongoConnectionWithoutDefault(name)
	if err == nil {
		return connMongo, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Mongo))...)

	connMysql, err := m.GetMySQLConnectionWithoutDefault(name)
	if err == nil {
		return connMysql, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Mysql))...)

	connNotion, err := m.GetNotionConnectionWithoutDefault(name)
	if err == nil {
		return connNotion, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Notion))...)

	connHANA, err := m.GetHANAConnectionWithoutDefault(name)
	if err == nil {
		return connHANA, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.HANA))...)

	connShopify, err := m.GetShopifyConnectionWithoutDefault(name)
	if err == nil {
		return connShopify, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Shopify))...)

	connGorgias, err := m.GetGorgiasConnectionWithoutDefault(name)
	if err == nil {
		return connGorgias, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Gorgias))...)

	connKlaviyo, err := m.GetKlaviyoConnectionWithoutDefault(name)
	if err == nil {
		return connKlaviyo, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Klaviyo))...)

	connAdjust, err := m.GetAdjustConnectionWithoutDefault(name)
	if err == nil {
		return connAdjust, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Adjust))...)

	athenaConnection, err := m.GetAthenaConnectionWithoutDefault(name)
	if err == nil {
		return athenaConnection, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Athena))...)

	connFacebookAds, err := m.GetFacebookAdsConnectionWithoutDefault(name)
	if err == nil {
		return connFacebookAds, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.FacebookAds))...)

	connStripe, err := m.GetStripeConnectionWithoutDefault(name)
	if err == nil {
		return connStripe, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Stripe))...)

	connAppsflyer, err := m.GetAppsflyerConnectionWithoutDefault(name)
	if err == nil {
		return connAppsflyer, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Appsflyer))...)

	connKafka, err := m.GetKafkaConnectionWithoutDefault(name)
	if err == nil {
		return connKafka, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Kafka))...)

	connDuckDB, err := m.GetDuckDBConnectionWithoutDefault(name)
	if err == nil {
		return connDuckDB, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.DuckDB))...)

	connClickHouse, err := m.GetClickHouseConnectionWithoutDefault(name)
	if err == nil {
		return connClickHouse, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.ClickHouse))...)

	connHubspot, err := m.GetHubspotConnectionWithoutDefault(name)
	if err == nil {
		return connHubspot, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Hubspot))...)

	connGoogleSheets, err := m.GetGoogleSheetsConnectionWithoutDefault(name)
	if err == nil {
		return connGoogleSheets, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.GoogleSheets))...)

	connChess, err := m.GetChessConnectionWithoutDefault(name)
	if err == nil {
		return connChess, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Chess))...)

	connAirtable, err := m.GetAirtableConnectionWithoutDefault(name)
	if err == nil {
		return connAirtable, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Airtable))...)

	connS3, err := m.GetS3ConnectionWithoutDefault(name)
	if err == nil {
		return connS3, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.S3))...)

	connSlack, err := m.GetSlackConnectionWithoutDefault(name)
	if err == nil {
		return connSlack, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Slack))...)

	connAsana, err := m.GetAsanaConnectionWithoutDefault(name)
	if err == nil {
		return connAsana, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Asana))...)

	connDynamoDB, err := m.GetDynamoDBConnectionWithoutDefault(name)
	if err == nil {
		return connDynamoDB, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.DynamoDB))...)

	connZendesk, err := m.GetZendeskConnectionWithoutDefault(name)
	if err == nil {
		return connZendesk, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Zendesk))...)

	connGoogleAds, err := m.GetGoogleAdsConnectionWithoutDefault(name)
	if err == nil {
		return connGoogleAds, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.GoogleAds))...)
	connTikTokAds, err := m.GetTikTokAdsConnectionWithoutDefault(name)
	if err == nil {
		return connTikTokAds, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.TikTokAds))...)

	connGitHub, err := m.GetGitHubConnectionWithoutDefault(name)
	if err == nil {
		return connGitHub, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.GitHub))...)

	connAppStore, err := m.GetAppStoreConnectionWithoutDefault(name)
	if err == nil {
		return connAppStore, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.AppStore))...)

	connLinkedInAds, err := m.GetLinkedInAdsConnectionWithoutDefault(name)
	if err == nil {
		return connLinkedInAds, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.LinkedInAds))...)

	connApplovinMax, err := m.GetApplovinMaxConnectionWithoutDefault(name)
	if err == nil {
		return connApplovinMax, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.ApplovinMax))...)

	connPersonio, err := m.GetPersonioConnectionWithoutDefault(name)
	if err == nil {
		return connPersonio, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Personio))...)

	connGCS, err := m.GetGCSConnectionWithoutDefault(name)
	if err == nil {
		return connGCS, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.GCS))...)

	connKinesis, err := m.GetKinesisConnectionWithoutDefault(name)
	if err == nil {
		return connKinesis, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Kinesis))...)

	connPipedrive, err := m.GetPipedriveConnectionWithoutDefault(name)
	if err == nil {
		return connPipedrive, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Pipedrive))...)

	connMixpanel, err := m.GetMixpanelConnectionWithoutDefault(name)
	if err == nil {
		return connMixpanel, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Mixpanel))...)

	connQuickBooks, err := m.GetQuickBooksConnectionWithoutDefault(name)
	if err == nil {
		return connQuickBooks, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.QuickBooks))...)

	connZoom, err := m.GetZoomConnectionWithoutDefault(name)
	if err == nil {
		return connZoom, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Zoom))...)

	connEMRServerless, err := m.GetEMRServerlessConnectionWithoutDefault(name)
	if err == nil {
		return connEMRServerless, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.EMRSeverless))...)

	connGoogleAnalytics, err := m.GetGoogleAnalyticsConnectionWithoutDefault(name)
	if err == nil {
		return connGoogleAnalytics, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.GoogleAnalytics))...)

	connAppLovin, err := m.GetAppLovinConnectionWithoutDefault(name)
	if err == nil {
		return connAppLovin, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.AppLovin))...)

	connFrankfurter, err := m.GetFrankfurterConnectionWithoutDefault(name)
	if err == nil {
		return connFrankfurter, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Frankfurter))...)
	connSalesforce, err := m.GetSalesforceConnectionWithoutDefault(name)
	if err == nil {
		return connSalesforce, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Salesforce))...)
	connSQLite, err := m.GetSQLiteConnectionWithoutDefault(name)
	if err == nil {
		return connSQLite, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.SQLite))...)

	connDB2, err := m.GetDB2ConnectionWithoutDefault(name)
	if err == nil {
		return connDB2, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.DB2))...)

	connOracle, err := m.GetOracleConnectionWithoutDefault(name)
	if err == nil {
		return connOracle, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Oracle))...)

	connPhantombuster, err := m.GetPhantombusterConnectionWithoutDefault(name)
	if err == nil {
		return connPhantombuster, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Phantombuster))...)

	connElasticsearch, err := m.GetElasticsearchConnectionWithoutDefault(name)
	if err == nil {
		return connElasticsearch, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Elasticsearch))...)

	connSpanner, err := m.GetSpannerConnectionWithoutDefault(name)
	if err == nil {
		return connSpanner, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Spanner))...)

	connSolidgate, err := m.GetSolidgateConnectionWithoutDefault(name)
	if err == nil {
		return connSolidgate, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Solidgate))...)

	connSmartsheet, err := m.GetSmartsheetConnectionWithoutDefault(name)
	if err == nil {
		return connSmartsheet, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Smartsheet))...)

	connAttio, err := m.GetAttioConnectionWithoutDefault(name)
	if err == nil {
		return connAttio, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Attio))...)

	connSftp, err := m.GetSftpConnectionWithoutDefault(name)
	if err == nil {
		return connSftp, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.Sftp))...)

	connPulse, err := m.GetISOCPulseConnectionWithoutDefault(name)
	if err == nil {
		return connPulse, nil
	}
	availableConnectionNames = append(availableConnectionNames, slices.Collect(maps.Keys(m.ISOCPulse))...)

	return nil, errors.Errorf("connection '%s' not found, available connection names are: %v", name, availableConnectionNames)
}

func (m *Manager) GetAthenaConnection(name string) (athena.Client, error) {
	db, err := m.GetAthenaConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetAthenaConnectionWithoutDefault("athena-default")
}

func (m *Manager) GetAthenaConnectionWithoutDefault(name string) (athena.Client, error) {
	if m.Athena == nil {
		return nil, errors.New("no Athena connections found")
	}

	db, ok := m.Athena[name]
	if !ok {
		return nil, errors.Errorf("Athena connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetDuckDBConnection(name string) (duck.DuckDBClient, error) {
	db, err := m.GetDuckDBConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetDuckDBConnectionWithoutDefault("duckdb-default")
}

func (m *Manager) GetDuckDBConnectionWithoutDefault(name string) (duck.DuckDBClient, error) {
	if m.DuckDB == nil {
		return nil, errors.New("no DuckDB connections found")
	}

	db, ok := m.DuckDB[name]
	if !ok {
		return nil, errors.Errorf("DuckDB connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetClickHouseConnection(name string) (clickhouse.ClickHouseClient, error) {
	db, err := m.GetClickHouseConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetClickHouseConnectionWithoutDefault("clickhouse-default")
}

func (m *Manager) GetClickHouseConnectionWithoutDefault(name string) (clickhouse.ClickHouseClient, error) {
	if m.ClickHouse == nil {
		return nil, errors.New("no clickhouse connections found")
	}
	db, ok := m.ClickHouse[name]
	if !ok {
		return nil, errors.Errorf("clickhouse connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetBqConnection(name string) (bigquery.DB, error) {
	return m.GetBqConnectionWithoutDefault(name)
}

func (m *Manager) GetBqConnectionWithoutDefault(name string) (bigquery.DB, error) {
	if m.BigQuery == nil {
		return nil, errors.New("no bigquery connections found")
	}

	db, ok := m.BigQuery[name]
	if !ok {
		return nil, errors.Errorf("bigquery connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetSfConnection(name string) (snowflake.SfClient, error) {
	db, err := m.GetSfConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetSfConnectionWithoutDefault("snowflake-default")
}

func (m *Manager) GetSfConnectionWithoutDefault(name string) (snowflake.SfClient, error) {
	if m.Snowflake == nil {
		return nil, errors.New("no snowflake connections found")
	}

	db, ok := m.Snowflake[name]
	if !ok {
		return nil, errors.Errorf("snowflake connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetPgConnection(name string) (postgres.PgClient, error) {
	db, err := m.GetPgConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetPgConnectionWithoutDefault("postgres-default")
}

func (m *Manager) GetPgConnectionWithoutDefault(name string) (postgres.PgClient, error) {
	if m.Postgres == nil {
		return nil, errors.New("no postgres/redshift connections found")
	}

	db, ok := m.Postgres[name]
	if !ok {
		return nil, errors.Errorf("postgres/redshift connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetMsConnection(name string) (mssql.MsClient, error) {
	db, err := m.GetMsConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetMsConnectionWithoutDefault("mssql-default")
}

func (m *Manager) GetMsConnectionWithoutDefault(name string) (mssql.MsClient, error) {
	if m.MsSQL == nil {
		return nil, errors.New("no mssql connections found")
	}

	db, ok := m.MsSQL[name]
	if !ok {
		return nil, errors.Errorf("mssql connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetDatabricksConnection(name string) (databricks.Client, error) {
	db, err := m.GetDatabricksConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetDatabricksConnectionWithoutDefault("databricks-default")
}

func (m *Manager) GetDatabricksConnectionWithoutDefault(name string) (databricks.Client, error) {
	if m.Databricks == nil {
		return nil, errors.New("no databricks connections found")
	}

	db, ok := m.Databricks[name]
	if !ok {
		return nil, errors.Errorf("databricks connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetMongoConnection(name string) (*mongo.DB, error) {
	db, err := m.GetMongoConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetMongoConnectionWithoutDefault("mongo-default")
}

func (m *Manager) GetMongoConnectionWithoutDefault(name string) (*mongo.DB, error) {
	if m.Mongo == nil {
		return nil, errors.New("no mongo connections found")
	}

	db, ok := m.Mongo[name]
	if !ok {
		return nil, errors.Errorf("mongo connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetMySQLConnection(name string) (*mysql.Client, error) {
	db, err := m.GetMySQLConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetMySQLConnectionWithoutDefault("mysql-default")
}

func (m *Manager) GetMySQLConnectionWithoutDefault(name string) (*mysql.Client, error) {
	if m.Mysql == nil {
		return nil, errors.New("no mysql connections found")
	}

	db, ok := m.Mysql[name]
	if !ok {
		return nil, errors.Errorf("mysql connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetNotionConnection(name string) (*notion.Client, error) {
	db, err := m.GetNotionConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetNotionConnectionWithoutDefault("notion-default")
}

func (m *Manager) GetNotionConnectionWithoutDefault(name string) (*notion.Client, error) {
	if m.Notion == nil {
		return nil, errors.New("no notion connections found")
	}

	db, ok := m.Notion[name]
	if !ok {
		return nil, errors.Errorf("notion connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetHANAConnection(name string) (*hana.Client, error) {
	db, err := m.GetHANAConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetHANAConnectionWithoutDefault("hana-default")
}

func (m *Manager) GetHANAConnectionWithoutDefault(name string) (*hana.Client, error) {
	if m.HANA == nil {
		return nil, errors.New("no hana connections found")
	}

	db, ok := m.HANA[name]
	if !ok {
		return nil, errors.Errorf("hana connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetShopifyConnection(name string) (*shopify.Client, error) {
	db, err := m.GetShopifyConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetShopifyConnectionWithoutDefault("shopify-default")
}

func (m *Manager) GetShopifyConnectionWithoutDefault(name string) (*shopify.Client, error) {
	if m.Shopify == nil {
		return nil, errors.New("no shopify connections found")
	}

	db, ok := m.Shopify[name]
	if !ok {
		return nil, errors.Errorf("shopify connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetKlaviyoConnection(name string) (*klaviyo.Client, error) {
	db, err := m.GetKlaviyoConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetKlaviyoConnectionWithoutDefault("klaviyo-default")
}

func (m *Manager) GetKlaviyoConnectionWithoutDefault(name string) (*klaviyo.Client, error) {
	if m.Klaviyo == nil {
		return nil, errors.New("no klaviyo connections found")
	}

	db, ok := m.Klaviyo[name]
	if !ok {
		return nil, errors.Errorf("klaviyo connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetSpannerConnection(name string) (*spanner.Client, error) {
	db, err := m.GetSpannerConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetSpannerConnectionWithoutDefault("spanner-default")
}

func (m *Manager) GetSpannerConnectionWithoutDefault(name string) (*spanner.Client, error) {
	if m.Spanner == nil {
		return nil, errors.New("no spanner connections found")
	}

	db, ok := m.Spanner[name]
	if !ok {
		return nil, errors.Errorf("spanner connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetSolidgateConnection(name string) (*solidgate.Client, error) {
	db, err := m.GetSolidgateConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetSolidgateConnectionWithoutDefault("solidgate-default")
}

func (m *Manager) GetSolidgateConnectionWithoutDefault(name string) (*solidgate.Client, error) {
	if m.Solidgate == nil {
		return nil, errors.New("no solidgate connections found")
	}

	db, ok := m.Solidgate[name]
	if !ok {
		return nil, errors.Errorf("solidgate connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetSmartsheetConnection(name string) (*smartsheet.Client, error) {
	db, err := m.GetSmartsheetConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetSmartsheetConnectionWithoutDefault("smartsheet-default")
}

func (m *Manager) GetSmartsheetConnectionWithoutDefault(name string) (*smartsheet.Client, error) {
	if m.Smartsheet == nil {
		return nil, errors.New("no smartsheet connections found")
	}

	db, ok := m.Smartsheet[name]
	if !ok {
		return nil, errors.Errorf("smartsheet connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetAttioConnection(name string) (*attio.Client, error) {
	db, err := m.GetAttioConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetAttioConnectionWithoutDefault("attio-default")
}

func (m *Manager) GetAttioConnectionWithoutDefault(name string) (*attio.Client, error) {
	if m.Attio == nil {
		return nil, errors.New("no attio connections found")
	}

	db, ok := m.Attio[name]
	if !ok {
		return nil, errors.Errorf("attio connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetSftpConnection(name string) (*sftp.Client, error) {
	db, err := m.GetSftpConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetSftpConnectionWithoutDefault("sftp-default")
}

func (m *Manager) GetSftpConnectionWithoutDefault(name string) (*sftp.Client, error) {
	if m.Sftp == nil {
		return nil, errors.New("no sftp connections found")
	}
	db, ok := m.Sftp[name]
	if !ok {
		return nil, errors.Errorf("sftp connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetAdjustConnection(name string) (*adjust.Client, error) {
	db, err := m.GetAdjustConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetAdjustConnectionWithoutDefault("adjust-default")
}

func (m *Manager) GetAdjustConnectionWithoutDefault(name string) (*adjust.Client, error) {
	if m.Adjust == nil {
		return nil, errors.New("no adjust connections found")
	}

	db, ok := m.Adjust[name]
	if !ok {
		return nil, errors.Errorf("adjust connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetStripeConnection(name string) (*stripe.Client, error) {
	db, err := m.GetStripeConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetStripeConnectionWithoutDefault("stripe-default")
}

func (m *Manager) GetStripeConnectionWithoutDefault(name string) (*stripe.Client, error) {
	if m.Stripe == nil {
		return nil, errors.New("no stripe connections found")
	}

	db, ok := m.Stripe[name]
	if !ok {
		return nil, errors.Errorf("stripe connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetGorgiasConnection(name string) (*gorgias.Client, error) {
	db, err := m.GetGorgiasConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetGorgiasConnectionWithoutDefault("gorgias-default")
}

func (m *Manager) GetGorgiasConnectionWithoutDefault(name string) (*gorgias.Client, error) {
	if m.Gorgias == nil {
		return nil, errors.New("no gorgias connections found")
	}

	db, ok := m.Gorgias[name]
	if !ok {
		return nil, errors.Errorf("hana gorgias not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetFacebookAdsConnection(name string) (*facebookads.Client, error) {
	db, err := m.GetFacebookAdsConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetFacebookAdsConnectionWithoutDefault("facebookads-default")
}

func (m *Manager) GetFacebookAdsConnectionWithoutDefault(name string) (*facebookads.Client, error) {
	if m.FacebookAds == nil {
		return nil, errors.New("no facebookads connections found")
	}

	db, ok := m.FacebookAds[name]
	if !ok {
		return nil, errors.Errorf("facebookads connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetAppsflyerConnection(name string) (*appsflyer.Client, error) {
	db, err := m.GetAppsflyerConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetAppsflyerConnectionWithoutDefault("appsflyer-default")
}

func (m *Manager) GetAppsflyerConnectionWithoutDefault(name string) (*appsflyer.Client, error) {
	if m.Appsflyer == nil {
		return nil, errors.New("no appsflyer connections found")
	}

	db, ok := m.Appsflyer[name]
	if !ok {
		return nil, errors.Errorf("appsflyer connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetKafkaConnection(name string) (*kafka.Client, error) {
	db, err := m.GetKafkaConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetKafkaConnectionWithoutDefault("kafka-default")
}

func (m *Manager) GetKafkaConnectionWithoutDefault(name string) (*kafka.Client, error) {
	if m.Kafka == nil {
		return nil, errors.New("no kafka connections found")
	}

	db, ok := m.Kafka[name]
	if !ok {
		return nil, errors.Errorf("kafka connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetHubspotConnection(name string) (*hubspot.Client, error) {
	db, err := m.GetHubspotConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetHubspotConnectionWithoutDefault("hubspot-default")
}

func (m *Manager) GetHubspotConnectionWithoutDefault(name string) (*hubspot.Client, error) {
	if m.Hubspot == nil {
		return nil, errors.New("no Hubspot connections found")
	}
	db, ok := m.Hubspot[name]
	if !ok {
		return nil, errors.Errorf("hubspot connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetAirtableConnection(name string) (*airtable.Client, error) {
	db, err := m.GetAirtableConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetAirtableConnectionWithoutDefault("airtable-default")
}

func (m *Manager) GetAirtableConnectionWithoutDefault(name string) (*airtable.Client, error) {
	if m.Airtable == nil {
		return nil, errors.New("no airtable connections found")
	}
	db, ok := m.Airtable[name]
	if !ok {
		return nil, errors.Errorf("airtable connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetGoogleSheetsConnection(name string) (*gsheets.Client, error) {
	db, err := m.GetGoogleSheetsConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}

	return m.GetGoogleSheetsConnectionWithoutDefault("google-sheets-default")
}

func (m *Manager) GetGoogleSheetsConnectionWithoutDefault(name string) (*gsheets.Client, error) {
	if m.GoogleSheets == nil {
		return nil, errors.New("no google sheets connections found")
	}
	db, ok := m.GoogleSheets[name]
	if !ok {
		return nil, errors.Errorf("google sheets connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetChessConnection(name string) (*chess.Client, error) {
	db, err := m.GetChessConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetChessConnectionWithoutDefault("chess-default")
}

func (m *Manager) GetChessConnectionWithoutDefault(name string) (*chess.Client, error) {
	if m.Chess == nil {
		return nil, errors.New("no chess connections found")
	}
	db, ok := m.Chess[name]
	if !ok {
		return nil, errors.Errorf("chess connection not found for '%s'", name)
	}

	return db, nil
}

func (m *Manager) GetZendeskConnection(name string) (*zendesk.Client, error) {
	if m.Zendesk == nil {
		return nil, errors.New("no zendesk connections found")
	}
	db, ok := m.Zendesk[name]
	if !ok {
		return nil, errors.Errorf("zendesk connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetZendeskConnectionWithoutDefault(name string) (*zendesk.Client, error) {
	if m.Zendesk == nil {
		return nil, errors.New("no zendesk connections found")
	}
	db, ok := m.Zendesk[name]
	if !ok {
		return nil, errors.Errorf("zendesk connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetS3Connection(name string) (*s3.Client, error) {
	db, err := m.GetS3ConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetS3ConnectionWithoutDefault("s3-default")
}

func (m *Manager) GetS3ConnectionWithoutDefault(name string) (*s3.Client, error) {
	if m.S3 == nil {
		return nil, errors.New("no s3 connections found")
	}
	db, ok := m.S3[name]
	if !ok {
		return nil, errors.Errorf("s3 connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetSlackConnection(name string) (*slack.Client, error) {
	db, err := m.GetSlackConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetSlackConnectionWithoutDefault("slack-default")
}

func (m *Manager) GetSlackConnectionWithoutDefault(name string) (*slack.Client, error) {
	if m.Slack == nil {
		return nil, errors.New("no slack connections found")
	}
	db, ok := m.Slack[name]
	if !ok {
		return nil, errors.Errorf("slack connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetAsanaConnection(name string) (*asana.Client, error) {
	db, err := m.GetAsanaConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetAsanaConnectionWithoutDefault("asana-default")
}

func (m *Manager) GetAsanaConnectionWithoutDefault(name string) (*asana.Client, error) {
	if m.Asana == nil {
		return nil, errors.New("no asana connections found")
	}
	db, ok := m.Asana[name]
	if !ok {
		return nil, errors.Errorf("asana connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetDynamoDBConnection(name string) (*dynamodb.Client, error) {
	db, err := m.GetDynamoDBConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetDynamoDBConnectionWithoutDefault("dynamodb-default")
}

func (m *Manager) GetDynamoDBConnectionWithoutDefault(name string) (*dynamodb.Client, error) {
	if m.DynamoDB == nil {
		return nil, errors.New("no dynamodb connections found")
	}
	db, ok := m.DynamoDB[name]
	if !ok {
		return nil, errors.Errorf("dynamodb connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetGoogleAdsConnection(name string) (*googleads.Client, error) {
	db, err := m.GetGoogleAdsConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetGoogleAdsConnectionWithoutDefault("googleads-default")
}

func (m *Manager) GetGoogleAdsConnectionWithoutDefault(name string) (*googleads.Client, error) {
	if m.GoogleAds == nil {
		return nil, errors.New("no googleads connections found")
	}
	db, ok := m.GoogleAds[name]
	if !ok {
		return nil, errors.Errorf("googleads connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetGitHubConnection(name string) (*github.Client, error) {
	db, err := m.GetGitHubConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetGitHubConnectionWithoutDefault("github-default")
}

func (m *Manager) GetGitHubConnectionWithoutDefault(name string) (*github.Client, error) {
	if m.GitHub == nil {
		return nil, errors.New("no github connections found")
	}
	db, ok := m.GitHub[name]
	if !ok {
		return nil, errors.Errorf("github connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetTikTokAdsConnection(name string) (*tiktokads.Client, error) {
	db, err := m.GetTikTokAdsConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetTikTokAdsConnectionWithoutDefault("tiktokads-default")
}

func (m *Manager) GetTikTokAdsConnectionWithoutDefault(name string) (*tiktokads.Client, error) {
	if m.TikTokAds == nil {
		return nil, errors.New("no tiktokads connections found")
	}
	db, ok := m.TikTokAds[name]
	if !ok {
		return nil, errors.Errorf("tiktokads connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetAppStoreConnection(name string) (*appstore.Client, error) {
	db, err := m.GetAppStoreConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetAppStoreConnectionWithoutDefault("appstore-default")
}

func (m *Manager) GetAppStoreConnectionWithoutDefault(name string) (*appstore.Client, error) {
	if m.AppStore == nil {
		return nil, errors.New("no appstore connections found")
	}
	db, ok := m.AppStore[name]
	if !ok {
		return nil, errors.Errorf("appstore connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetLinkedInAdsConnection(name string) (*linkedinads.Client, error) {
	db, err := m.GetLinkedInAdsConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetLinkedInAdsConnectionWithoutDefault("linkedinads-default")
}

func (m *Manager) GetLinkedInAdsConnectionWithoutDefault(name string) (*linkedinads.Client, error) {
	if m.LinkedInAds == nil {
		return nil, errors.New("no linkedinads connections found")
	}
	db, ok := m.LinkedInAds[name]
	if !ok {
		return nil, errors.Errorf("linkedinads connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetGCSConnection(name string) (*gcs.Client, error) {
	db, err := m.GetGCSConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetGCSConnectionWithoutDefault("gcs-default")
}

func (m *Manager) GetGCSConnectionWithoutDefault(name string) (*gcs.Client, error) {
	if m.GCS == nil {
		return nil, errors.New("no gcs connections found")
	}
	db, ok := m.GCS[name]
	if !ok {
		return nil, errors.Errorf("gcs connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetApplovinMaxConnection(name string) (*applovinmax.Client, error) {
	db, err := m.GetApplovinMaxConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetApplovinMaxConnectionWithoutDefault("applovinmax-default")
}

func (m *Manager) GetApplovinMaxConnectionWithoutDefault(name string) (*applovinmax.Client, error) {
	if m.ApplovinMax == nil {
		return nil, errors.New("no applovinmax connections found")
	}
	db, ok := m.ApplovinMax[name]
	if !ok {
		return nil, errors.Errorf("applovinmax connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetPersonioConnection(name string) (*personio.Client, error) {
	db, err := m.GetPersonioConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetPersonioConnectionWithoutDefault("personio-default")
}

func (m *Manager) GetPersonioConnectionWithoutDefault(name string) (*personio.Client, error) {
	if m.Personio == nil {
		return nil, errors.New("no personio connections found")
	}
	db, ok := m.Personio[name]
	if !ok {
		return nil, errors.Errorf("personio connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetKinesisConnection(name string) (*kinesis.Client, error) {
	db, err := m.GetKinesisConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetKinesisConnectionWithoutDefault("kinesis-default")
}

func (m *Manager) GetKinesisConnectionWithoutDefault(name string) (*kinesis.Client, error) {
	if m.Kinesis == nil {
		return nil, errors.New("no kinesis connections found")
	}
	db, ok := m.Kinesis[name]
	if !ok {
		return nil, errors.Errorf("kinesis connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetPipedriveConnection(name string) (*pipedrive.Client, error) {
	db, err := m.GetPipedriveConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetPipedriveConnectionWithoutDefault("pipedrive-default")
}

func (m *Manager) GetPipedriveConnectionWithoutDefault(name string) (*pipedrive.Client, error) {
	if m.Pipedrive == nil {
		return nil, errors.New("no pipedrive connections found")
	}
	db, ok := m.Pipedrive[name]
	if !ok {
		return nil, errors.Errorf("pipedrive connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetMixpanelConnection(name string) (*mixpanel.Client, error) {
	db, err := m.GetMixpanelConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetMixpanelConnectionWithoutDefault("mixpanel-default")
}

func (m *Manager) GetMixpanelConnectionWithoutDefault(name string) (*mixpanel.Client, error) {
	if m.Mixpanel == nil {
		return nil, errors.New("no mixpanel connections found")
	}
	db, ok := m.Mixpanel[name]
	if !ok {
		return nil, errors.Errorf("mixpanel connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetQuickBooksConnection(name string) (*quickbooks.Client, error) {
	db, err := m.GetQuickBooksConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetQuickBooksConnectionWithoutDefault("quickbooks-default")
}

func (m *Manager) GetQuickBooksConnectionWithoutDefault(name string) (*quickbooks.Client, error) {
	if m.QuickBooks == nil {
		return nil, errors.New("no quickbooks connections found")
	}
	db, ok := m.QuickBooks[name]
	if !ok {
		return nil, errors.Errorf("quickbooks connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetZoomConnection(name string) (*zoom.Client, error) {
	db, err := m.GetZoomConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetZoomConnectionWithoutDefault("zoom-default")
}

func (m *Manager) GetZoomConnectionWithoutDefault(name string) (*zoom.Client, error) {
	if m.Zoom == nil {
		return nil, errors.New("no zoom connections found")
	}
	db, ok := m.Zoom[name]
	if !ok {
		return nil, errors.Errorf("zoom connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetEMRServerlessConnection(name string) (*emr_serverless.Client, error) {
	db, err := m.GetEMRServerlessConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetEMRServerlessConnectionWithoutDefault("emr_serverless-default")
}

func (m *Manager) GetEMRServerlessConnectionWithoutDefault(name string) (*emr_serverless.Client, error) {
	if m.EMRSeverless == nil {
		return nil, errors.New("no EMR Serverless connections found")
	}
	db, ok := m.EMRSeverless[name]
	if !ok {
		return nil, errors.Errorf("EMR Serverless connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetISOCPulseConnection(name string) (*isocpulse.Client, error) {
	db, err := m.GetISOCPulseConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetISOCPulseConnectionWithoutDefault("isoc_pulse-default")
}

func (m *Manager) GetISOCPulseConnectionWithoutDefault(name string) (*isocpulse.Client, error) {
	if m.ISOCPulse == nil {
		return nil, errors.New("no ISOC Pulse connections found")
	}
	db, ok := m.ISOCPulse[name]
	if !ok {
		return nil, errors.Errorf("ISOC Pulse connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetGoogleAnalyticsConnection(name string) (*googleanalytics.Client, error) {
	db, err := m.GetGoogleAnalyticsConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetGoogleAnalyticsConnectionWithoutDefault("googleanalytics-default")
}

func (m *Manager) GetGoogleAnalyticsConnectionWithoutDefault(name string) (*googleanalytics.Client, error) {
	if m.GoogleAnalytics == nil {
		return nil, errors.New("no googleanalytics connections found")
	}
	db, ok := m.GoogleAnalytics[name]
	if !ok {
		return nil, errors.Errorf("googleanalytics connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetAppLovinConnection(name string) (*applovin.Client, error) {
	db, err := m.GetAppLovinConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetAppLovinConnectionWithoutDefault("applovin-default")
}

func (m *Manager) GetAppLovinConnectionWithoutDefault(name string) (*applovin.Client, error) {
	if m.AppLovin == nil {
		return nil, errors.New("no applovin connections found")
	}
	db, ok := m.AppLovin[name]
	if !ok {
		return nil, errors.Errorf("applovin connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetFrankfurterConnection(name string) (*frankfurter.Client, error) {
	db, err := m.GetFrankfurterConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetFrankfurterConnectionWithoutDefault("frankfurter-default")
}

func (m *Manager) GetFrankfurterConnectionWithoutDefault(name string) (*frankfurter.Client, error) {
	if m.Frankfurter == nil {
		return nil, errors.New("no frankfurter connections found")
	}
	db, ok := m.Frankfurter[name]
	if !ok {
		return nil, errors.Errorf("frankfurter connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetSalesforceConnection(name string) (*salesforce.Client, error) {
	db, err := m.GetSalesforceConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetSalesforceConnectionWithoutDefault("salesforce-default")
}

func (m *Manager) GetSalesforceConnectionWithoutDefault(name string) (*salesforce.Client, error) {
	if m.Salesforce == nil {
		return nil, errors.New("no salesforce connections found")
	}
	db, ok := m.Salesforce[name]
	if !ok {
		return nil, errors.Errorf("salesforce connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetSQLiteConnection(name string) (*sqlite.Client, error) {
	db, err := m.GetSQLiteConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetSQLiteConnectionWithoutDefault("sqlite-default")
}

func (m *Manager) GetSQLiteConnectionWithoutDefault(name string) (*sqlite.Client, error) {
	if m.SQLite == nil {
		return nil, errors.New("no sqlite connections found")
	}
	db, ok := m.SQLite[name]
	if !ok {
		return nil, errors.Errorf("sqlite connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetOracleConnection(name string) (*oracle.Client, error) {
	db, err := m.GetOracleConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetOracleConnectionWithoutDefault("oracle-default")
}

func (m *Manager) GetOracleConnectionWithoutDefault(name string) (*oracle.Client, error) {
	if m.Oracle == nil {
		return nil, errors.New("no oracle connections found")
	}
	db, ok := m.Oracle[name]
	if !ok {
		return nil, errors.Errorf("oracle connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetPhantombusterConnection(name string) (*phantombuster.Client, error) {
	db, err := m.GetPhantombusterConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetPhantombusterConnectionWithoutDefault("phantombuster-default")
}

func (m *Manager) GetPhantombusterConnectionWithoutDefault(name string) (*phantombuster.Client, error) {
	if m.Phantombuster == nil {
		return nil, errors.New("no phantombuster connections found")
	}
	db, ok := m.Phantombuster[name]
	if !ok {
		return nil, errors.Errorf("phantombuster connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetElasticsearchConnection(name string) (*elasticsearch.Client, error) {
	db, err := m.GetElasticsearchConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetElasticsearchConnectionWithoutDefault("elasticsearch-default")
}

func (m *Manager) GetElasticsearchConnectionWithoutDefault(name string) (*elasticsearch.Client, error) {
	if m.Elasticsearch == nil {
		return nil, errors.New("no elasticsearch connections found")
	}
	db, ok := m.Elasticsearch[name]
	if !ok {
		return nil, errors.Errorf("elasticsearch connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) GetDB2Connection(name string) (*db2.Client, error) {
	db, err := m.GetDB2ConnectionWithoutDefault(name)
	if err == nil {
		return db, nil
	}
	return m.GetDB2ConnectionWithoutDefault("db2-default")
}

func (m *Manager) GetDB2ConnectionWithoutDefault(name string) (*db2.Client, error) {
	if m.DB2 == nil {
		return nil, errors.New("no db2 connections found")
	}
	db, ok := m.DB2[name]
	if !ok {
		return nil, errors.Errorf("db2 connection not found for '%s'", name)
	}
	return db, nil
}

func (m *Manager) AddBqConnectionFromConfig(connection *config.GoogleCloudPlatformConnection) error {
	m.mutex.Lock()
	if m.BigQuery == nil {
		m.BigQuery = make(map[string]*bigquery.Client)
	}
	m.mutex.Unlock()

	// Check if either ServiceAccountFile or ServiceAccountJSON is provided, prioritizing ServiceAccountFile.
	if len(connection.ServiceAccountFile) == 0 && len(connection.ServiceAccountJSON) == 0 {
		return errors.New("credentials are required: provide either service_account_file or service_account_json")
	}

	// Validate ServiceAccountFile if provided.
	if len(connection.ServiceAccountFile) > 0 {
		if err := validateServiceAccountFile(connection.ServiceAccountFile); err != nil {
			return err
		}
	}

	// Validate ServiceAccountJSON if provided.
	if len(connection.ServiceAccountJSON) > 0 {
		if err := validateServiceAccountJSON(connection.ServiceAccountJSON); err != nil {
			return err
		}
	}

	// Set up the BigQuery client using the preferred credentials.
	db, err := bigquery.NewDB(&bigquery.Config{
		ProjectID:           connection.ProjectID,
		CredentialsFilePath: connection.ServiceAccountFile,
		CredentialsJSON:     connection.ServiceAccountJSON,
		Credentials:         connection.GetCredentials(),
		Location:            connection.Location,
	})
	if err != nil {
		return err
	}

	// Lock and store the new BigQuery client.
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.BigQuery[connection.Name] = db

	return nil
}

func (m *Manager) AddSfConnectionFromConfig(connection *config.SnowflakeConnection) error {
	m.mutex.Lock()
	if m.Snowflake == nil {
		m.Snowflake = make(map[string]*snowflake.DB)
	}
	m.mutex.Unlock()

	privateKey := ""

	// Prioritize the direct PrivateKey field over PrivateKeyPath
	if connection.PrivateKey != "" {
		privateKey = connection.PrivateKey
	} else if connection.PrivateKeyPath != "" {
		privateKeyBytes, err := os.ReadFile(connection.PrivateKeyPath)
		if err != nil {
			return err
		}
		privateKey = string(privateKeyBytes)
	}

	db, err := snowflake.NewDB(&snowflake.Config{
		Account:    connection.Account,
		Username:   connection.Username,
		Password:   connection.Password,
		Region:     connection.Region,
		Role:       connection.Role,
		Database:   connection.Database,
		Schema:     connection.Schema,
		Warehouse:  connection.Warehouse,
		PrivateKey: privateKey,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Snowflake[connection.Name] = db

	return nil
}

func (m *Manager) AddAthenaConnectionFromConfig(connection *config.AthenaConnection) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.Athena == nil {
		m.Athena = make(map[string]*athena.DB)
	}

	if connection.Profile != "" {
		if connection.AccessKey != "" || connection.SecretKey != "" {
			return errors.New("access_key and secret_key cannot be provided when profile is specified, you need specify either profile or access_key and secret_key")
		}
		err := connection.LoadCredentialsFromProfile(context.Background())
		if err != nil {
			return err
		}
	}

	m.Athena[connection.Name] = athena.NewDB(&athena.Config{
		Region:          connection.Region,
		OutputBucket:    connection.QueryResultsPath,
		AccessID:        connection.AccessKey,
		SecretAccessKey: connection.SecretKey,
		SessionToken:    connection.SessionToken,
		Database:        connection.Database,
	})

	return nil
}

func (m *Manager) AddPgConnectionFromConfig(connection *config.PostgresConnection) error {
	return m.addPgLikeConnectionFromConfig(connection, false)
}

func (m *Manager) AddRedshiftConnectionFromConfig(connection *config.RedshiftConnection) error {
	return m.addRedshiftConnectionFromConfig(connection)
}

func (m *Manager) addRedshiftConnectionFromConfig(connection *config.RedshiftConnection) error {
	m.mutex.Lock()
	if m.Postgres == nil {
		m.Postgres = make(map[string]*postgres.Client)
	}
	m.mutex.Unlock()

	var client *postgres.Client
	var err error
	client, err = postgres.NewClient(context.TODO(), postgres.RedShiftConfig{
		Username: connection.Username,
		Password: connection.Password,
		Host:     connection.Host,
		Port:     connection.Port,
		Database: connection.Database,
		Schema:   connection.Schema,
		SslMode:  connection.SslMode,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Postgres[connection.Name] = client

	return nil
}

func (m *Manager) addPgLikeConnectionFromConfig(connection *config.PostgresConnection, redshift bool) error {
	m.mutex.Lock()
	if m.Postgres == nil {
		m.Postgres = make(map[string]*postgres.Client)
	}
	m.mutex.Unlock()

	poolMaxConns := connection.PoolMaxConns
	if connection.PoolMaxConns == 0 {
		poolMaxConns = 10
	}

	var client *postgres.Client
	var err error
	if redshift {
		client, err = postgres.NewClient(context.TODO(), postgres.RedShiftConfig{
			Username: connection.Username,
			Password: connection.Password,
			Host:     connection.Host,
			Port:     connection.Port,
			Database: connection.Database,
			Schema:   connection.Schema,
			SslMode:  connection.SslMode,
		})
	} else {
		client, err = postgres.NewClient(context.TODO(), postgres.Config{
			Username:     connection.Username,
			Password:     connection.Password,
			Host:         connection.Host,
			Port:         connection.Port,
			Database:     connection.Database,
			Schema:       connection.Schema,
			PoolMaxConns: poolMaxConns,
			SslMode:      connection.SslMode,
		})
	}

	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Postgres[connection.Name] = client

	return nil
}

func (m *Manager) AddMsSQLConnectionFromConfig(connection *config.MsSQLConnection) error {
	m.mutex.Lock()
	if m.MsSQL == nil {
		m.MsSQL = make(map[string]*mssql.DB)
	}
	m.mutex.Unlock()

	client, err := mssql.NewDB(&mssql.Config{
		Username: connection.Username,
		Password: connection.Password,
		Host:     connection.Host,
		Port:     connection.Port,
		Database: connection.Database,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.MsSQL[connection.Name] = client

	return nil
}

func (m *Manager) AddSynapseSQLConnectionFromConfig(connection *config.SynapseConnection) error {
	m.mutex.Lock()
	if m.MsSQL == nil {
		m.MsSQL = make(map[string]*mssql.DB)
	}
	m.mutex.Unlock()

	client, err := mssql.NewDB(&mssql.Config{
		Username: connection.Username,
		Password: connection.Password,
		Host:     connection.Host,
		Port:     connection.Port,
		Database: connection.Database,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.MsSQL[connection.Name] = client

	return nil
}

func (m *Manager) AddDatabricksConnectionFromConfig(connection *config.DatabricksConnection) error {
	m.mutex.Lock()
	if m.Databricks == nil {
		m.Databricks = make(map[string]*databricks.DB)
	}
	m.mutex.Unlock()

	client, err := databricks.NewDB(&databricks.Config{
		Token:   connection.Token,
		Host:    connection.Host,
		Path:    connection.Path,
		Port:    connection.Port,
		Catalog: connection.Catalog,
		Schema:  connection.Schema,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Databricks[connection.Name] = client

	return nil
}

func (m *Manager) AddMongoConnectionFromConfig(connection *config.MongoConnection) error {
	m.mutex.Lock()
	if m.Mongo == nil {
		m.Mongo = make(map[string]*mongo.DB)
	}
	m.mutex.Unlock()

	client, err := mongo.NewDB(&mongo.Config{
		Username: connection.Username,
		Password: connection.Password,
		Host:     connection.Host,
		Port:     connection.Port,
		Database: connection.Database,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Mongo[connection.Name] = client

	return nil
}

func (m *Manager) AddMySQLConnectionFromConfig(connection *config.MySQLConnection) error {
	m.mutex.Lock()
	if m.Mysql == nil {
		m.Mysql = make(map[string]*mysql.Client)
	}
	m.mutex.Unlock()

	client, err := mysql.NewClient(&mysql.Config{
		Username: connection.Username,
		Password: connection.Password,
		Host:     connection.Host,
		Port:     connection.Port,
		Database: connection.Database,
		Driver:   connection.Driver,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Mysql[connection.Name] = client

	return nil
}

func (m *Manager) AddNotionConnectionFromConfig(connection *config.NotionConnection) error {
	m.mutex.Lock()
	if m.Notion == nil {
		m.Notion = make(map[string]*notion.Client)
	}
	m.mutex.Unlock()

	client, err := notion.NewClient(&notion.Config{
		APIKey: connection.APIKey,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Notion[connection.Name] = client

	return nil
}

func (m *Manager) AddShopifyConnectionFromConfig(connection *config.ShopifyConnection) error {
	m.mutex.Lock()
	if m.Shopify == nil {
		m.Shopify = make(map[string]*shopify.Client)
	}
	m.mutex.Unlock()

	client, err := shopify.NewClient(&shopify.Config{
		APIKey: connection.APIKey,
		URL:    connection.URL,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Shopify[connection.Name] = client

	return nil
}

func (m *Manager) AddGorgiasConnectionFromConfig(connection *config.GorgiasConnection) error {
	m.mutex.Lock()
	if m.Gorgias == nil {
		m.Gorgias = make(map[string]*gorgias.Client)
	}
	m.mutex.Unlock()

	client, err := gorgias.NewClient(&gorgias.Config{
		APIKey: connection.APIKey,
		Domain: connection.Domain,
		Email:  connection.Email,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Gorgias[connection.Name] = client

	return nil
}

func (m *Manager) AddKlaviyoConnectionFromConfig(connection *config.KlaviyoConnection) error {
	m.mutex.Lock()
	if m.Klaviyo == nil {
		m.Klaviyo = make(map[string]*klaviyo.Client)
	}
	m.mutex.Unlock()

	client, err := klaviyo.NewClient(klaviyo.Config{
		APIKey: connection.APIKey,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Klaviyo[connection.Name] = client

	return nil
}

func (m *Manager) AddAdjustConnectionFromConfig(connection *config.AdjustConnection) error {
	m.mutex.Lock()
	if m.Adjust == nil {
		m.Adjust = make(map[string]*adjust.Client)
	}
	m.mutex.Unlock()

	client, err := adjust.NewClient(adjust.Config{
		APIKey: connection.APIKey,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Adjust[connection.Name] = client

	return nil
}

func (m *Manager) AddHANAConnectionFromConfig(connection *config.HANAConnection) error {
	m.mutex.Lock()
	if m.HANA == nil {
		m.HANA = make(map[string]*hana.Client)
	}
	m.mutex.Unlock()

	client, err := hana.NewClient(&hana.Config{
		Username: connection.Username,
		Password: connection.Password,
		Host:     connection.Host,
		Port:     connection.Port,
		Database: connection.Database,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.HANA[connection.Name] = client

	return nil
}

func (m *Manager) AddFacebookAdsConnectionFromConfig(connection *config.FacebookAdsConnection) error {
	m.mutex.Lock()
	if m.FacebookAds == nil {
		m.FacebookAds = make(map[string]*facebookads.Client)
	}
	m.mutex.Unlock()

	client, err := facebookads.NewClient(facebookads.Config{
		AccessToken: connection.AccessToken,
		AccountID:   connection.AccountID,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.FacebookAds[connection.Name] = client

	return nil
}

func (m *Manager) AddStripeConnectionFromConfig(connection *config.StripeConnection) error {
	m.mutex.Lock()
	if m.Stripe == nil {
		m.Stripe = make(map[string]*stripe.Client)
	}
	m.mutex.Unlock()

	client, err := stripe.NewClient(&stripe.Config{
		APIKey: connection.APIKey,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Stripe[connection.Name] = client

	return nil
}

func (m *Manager) AddAppsflyerConnectionFromConfig(connection *config.AppsflyerConnection) error {
	m.mutex.Lock()
	if m.Appsflyer == nil {
		m.Appsflyer = make(map[string]*appsflyer.Client)
	}
	m.mutex.Unlock()

	client, err := appsflyer.NewClient(appsflyer.Config{
		APIKey: connection.APIKey,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Appsflyer[connection.Name] = client

	return nil
}

func (m *Manager) AddGoogleSheetsConnectionFromConfig(connection *config.GoogleSheetsConnection) error {
	m.mutex.Lock()
	if m.GoogleSheets == nil {
		m.GoogleSheets = make(map[string]*gsheets.Client)
	}
	m.mutex.Unlock()
	if len(connection.ServiceAccountFile) == 0 && len(connection.ServiceAccountJSON) == 0 {
		return errors.New("credentials are required: provide either service_account_file or service_account_json")
	}

	// Validate ServiceAccountFile if provided.
	if len(connection.ServiceAccountFile) > 0 {
		if err := validateServiceAccountFile(connection.ServiceAccountFile); err != nil {
			return err
		}
	}

	// Validate ServiceAccountJSON if provided.
	if len(connection.ServiceAccountJSON) > 0 {
		if err := validateServiceAccountJSON(connection.ServiceAccountJSON); err != nil {
			return err
		}
	}
	client, err := gsheets.NewClient(gsheets.Config{
		ServiceAccountFile: connection.ServiceAccountFile,
		ServiceAccountJSON: connection.ServiceAccountJSON,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.GoogleSheets[connection.Name] = client
	return nil
}

func (m *Manager) AddSpannerConnectionFromConfig(connection *config.SpannerConnection) error {
	m.mutex.Lock()
	if m.Spanner == nil {
		m.Spanner = make(map[string]*spanner.Client)
	}
	m.mutex.Unlock()

	if len(connection.ServiceAccountJSON) == 0 && len(connection.ServiceAccountFile) == 0 {
		return errors.New("credentials are required: provide either service account file or service account json")
	}

	client, err := spanner.NewClient(spanner.Config{
		ProjectID:          connection.ProjectID,
		InstanceID:         connection.InstanceID,
		Database:           connection.Database,
		ServiceAccountJSON: connection.ServiceAccountJSON,
		ServiceAccountFile: connection.ServiceAccountFile,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Spanner[connection.Name] = client

	return nil
}

func (m *Manager) AddSolidgateConnectionFromConfig(connection *config.SolidgateConnection) error {
	m.mutex.Lock()
	if m.Solidgate == nil {
		m.Solidgate = make(map[string]*solidgate.Client)
	}
	m.mutex.Unlock()

	client, err := solidgate.NewClient(solidgate.Config{
		SecretKey: connection.SecretKey,
		PublicKey: connection.PublicKey,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Solidgate[connection.Name] = client

	return nil
}

func (m *Manager) AddSmartsheetConnectionFromConfig(connection *config.SmartsheetConnection) error {
	m.mutex.Lock()
	if m.Smartsheet == nil {
		m.Smartsheet = make(map[string]*smartsheet.Client)
	}
	m.mutex.Unlock()

	client, err := smartsheet.NewClient(smartsheet.Config{
		AccessToken: connection.AccessToken,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Smartsheet[connection.Name] = client

	return nil
}

func (m *Manager) AddSftpConnectionFromConfig(connection *config.SFTPConnection) error {
	m.mutex.Lock()
	if m.Sftp == nil {
		m.Sftp = make(map[string]*sftp.Client)
	}
	m.mutex.Unlock()

	client, err := sftp.NewClient(sftp.Config{
		Host:     connection.Host,
		Port:     connection.Port,
		Username: connection.Username,
		Password: connection.Password,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Sftp[connection.Name] = client

	return nil
}

func (m *Manager) AddKafkaConnectionFromConfig(connection *config.KafkaConnection) error {
	m.mutex.Lock()
	if m.Kafka == nil {
		m.Kafka = make(map[string]*kafka.Client)
	}
	m.mutex.Unlock()

	client, err := kafka.NewClient(kafka.Config{
		BootstrapServers: connection.BootstrapServers,
		GroupID:          connection.GroupID,
		BatchSize:        connection.BatchSize,
		SaslMechanisms:   connection.SaslMechanisms,
		SaslUsername:     connection.SaslUsername,
		SaslPassword:     connection.SaslPassword,
		BatchTimeout:     connection.BatchTimeout,
	})
	if err != nil {
		return err
	}

	m.Kafka[connection.Name] = client

	return nil
}

func (m *Manager) AddDuckDBConnectionFromConfig(connection *config.DuckDBConnection) error {
	m.mutex.Lock()
	if m.DuckDB == nil {
		m.DuckDB = make(map[string]*duck.Client)
	}
	m.mutex.Unlock()

	client, err := duck.NewClient(duck.Config{
		Path: connection.Path,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.DuckDB[connection.Name] = client

	return nil
}

func (m *Manager) AddClickHouseConnectionFromConfig(connection *config.ClickHouseConnection) error {
	m.mutex.Lock()
	if m.ClickHouse == nil {
		m.ClickHouse = make(map[string]*clickhouse.Client)
	}
	m.mutex.Unlock()

	client, err := clickhouse.NewClient(&clickhouse.Config{
		Host:     connection.Host,
		Port:     connection.Port,
		Username: connection.Username,
		Password: connection.Password,
		Database: connection.Database,
		HTTPPort: connection.HTTPPort,
		Secure:   connection.Secure,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.ClickHouse[connection.Name] = client

	return nil
}

func (m *Manager) AddChessConnectionFromConfig(connection *config.ChessConnection) error {
	m.mutex.Lock()
	if m.Chess == nil {
		m.Chess = make(map[string]*chess.Client)
	}
	m.mutex.Unlock()
	client, err := chess.NewClient(chess.Config{
		Players: connection.Players,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Chess[connection.Name] = client
	return nil
}

func (m *Manager) AddHubspotConnectionFromConfig(connection *config.HubspotConnection) error {
	m.mutex.Lock()
	if m.Hubspot == nil {
		m.Hubspot = make(map[string]*hubspot.Client)
	}
	m.mutex.Unlock()

	client, err := hubspot.NewClient(hubspot.Config{
		APIKey: connection.APIKey,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Hubspot[connection.Name] = client

	return nil
}

func (m *Manager) AddAirtableConnectionFromConfig(connection *config.AirtableConnection) error {
	m.mutex.Lock()
	if m.Airtable == nil {
		m.Airtable = make(map[string]*airtable.Client)
	}
	m.mutex.Unlock()
	client, err := airtable.NewClient(airtable.Config{
		BaseID:      connection.BaseID,
		AccessToken: connection.AccessToken,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Airtable[connection.Name] = client

	return nil
}

func (m *Manager) AddS3ConnectionFromConfig(connection *config.S3Connection) error {
	m.mutex.Lock()
	if m.S3 == nil {
		m.S3 = make(map[string]*s3.Client)
	}
	m.mutex.Unlock()
	client, err := s3.NewClient(s3.Config{
		BucketName:      connection.BucketName,
		PathToFile:      connection.PathToFile,
		AccessKeyID:     connection.AccessKeyID,
		SecretAccessKey: connection.SecretAccessKey,
		EndpointURL:     connection.EndpointURL,
		Layout:          connection.Layout,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.S3[connection.Name] = client
	return nil
}

func (m *Manager) AddSlackConnectionFromConfig(connection *config.SlackConnection) error {
	m.mutex.Lock()
	if m.Slack == nil {
		m.Slack = make(map[string]*slack.Client)
	}
	m.mutex.Unlock()
	client, err := slack.NewClient(slack.Config{
		APIKey: connection.APIKey,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Slack[connection.Name] = client
	return nil
}

func (m *Manager) AddAsanaConnectionFromConfig(connection *config.AsanaConnection) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.Asana == nil {
		m.Asana = make(map[string]*asana.Client)
	}

	client, err := asana.NewClient(asana.Config{
		WorkspaceID: connection.WorkspaceID,
		AccessToken: connection.AccessToken,
	})
	if err != nil {
		return err
	}
	m.Asana[connection.Name] = client
	return nil
}

func (m *Manager) AddDynamoDBConnectionFromConfig(connection *config.DynamoDBConnection) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.DynamoDB == nil {
		m.DynamoDB = make(map[string]*dynamodb.Client)
	}

	client, err := dynamodb.NewClient(dynamodb.Config{
		AccessKeyID:     connection.AccessKeyID,
		SecretAccessKey: connection.SecretAccessKey,
		Region:          connection.Region,
	})
	if err != nil {
		return err
	}
	m.DynamoDB[connection.Name] = client
	return nil
}

func (m *Manager) AddGoogleAdsConnectionFromConfig(connection *config.GoogleAdsConnection) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.GoogleAds == nil {
		m.GoogleAds = make(map[string]*googleads.Client)
	}

	client, err := googleads.NewClient(googleads.Config{
		CustomerID:         connection.CustomerID,
		DeveloperToken:     connection.DeveloperToken,
		ServiceAccountFile: connection.ServiceAccountFile,
		ServiceAccountJSON: connection.ServiceAccountJSON,
	})
	if err != nil {
		return err
	}
	m.GoogleAds[connection.Name] = client
	return nil
}

func (m *Manager) AddZendeskConnectionFromConfig(connection *config.ZendeskConnection) error {
	m.mutex.Lock()
	if m.Zendesk == nil {
		m.Zendesk = make(map[string]*zendesk.Client)
	}
	m.mutex.Unlock()
	client, err := zendesk.NewClient(zendesk.Config{
		APIToken:   connection.APIToken,
		Email:      connection.Email,
		OAuthToken: connection.OAuthToken,
		Subdomain:  connection.Subdomain,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Zendesk[connection.Name] = client

	return nil
}

func (m *Manager) AddTikTokAdsConnectionFromConfig(connection *config.TikTokAdsConnection) error {
	m.mutex.Lock()
	if m.TikTokAds == nil {
		m.TikTokAds = make(map[string]*tiktokads.Client)
	}
	m.mutex.Unlock()
	client, err := tiktokads.NewClient(tiktokads.Config{
		AccessToken:   connection.AccessToken,
		AdvertiserIDs: connection.AdvertiserIDs,
		Timezone:      connection.Timezone,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.TikTokAds[connection.Name] = client

	return nil
}

func (m *Manager) AddGitHubConnectionFromConfig(connection *config.GitHubConnection) error {
	m.mutex.Lock()
	if m.GitHub == nil {
		m.GitHub = make(map[string]*github.Client)
	}
	m.mutex.Unlock()

	client, err := github.NewClient(github.Config{
		AccessToken: connection.AccessToken,
		Owner:       connection.Owner,
		Repo:        connection.Repo,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.GitHub[connection.Name] = client

	return nil
}

func (m *Manager) AddAppStoreConnectionFromConfig(connection *config.AppStoreConnection) error {
	m.mutex.Lock()
	if m.AppStore == nil {
		m.AppStore = make(map[string]*appstore.Client)
	}
	m.mutex.Unlock()

	client, err := appstore.NewClient(appstore.Config{
		IssuerID: connection.IssuerID,
		KeyID:    connection.KeyID,
		KeyPath:  connection.KeyPath,
		Key:      connection.Key,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.AppStore[connection.Name] = client

	return nil
}

func (m *Manager) AddLinkedInAdsConnectionFromConfig(connection *config.LinkedInAdsConnection) error {
	m.mutex.Lock()
	if m.LinkedInAds == nil {
		m.LinkedInAds = make(map[string]*linkedinads.Client)
	}
	m.mutex.Unlock()
	client, err := linkedinads.NewClient(linkedinads.Config{
		AccessToken: connection.AccessToken,
		AccountIds:  connection.AccountIds,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.LinkedInAds[connection.Name] = client
	return nil
}

func (m *Manager) AddGCSConnectionFromConfig(connection *config.GCSConnection) error {
	m.mutex.Lock()
	if m.GCS == nil {
		m.GCS = make(map[string]*gcs.Client)
	}
	m.mutex.Unlock()

	client, err := gcs.NewClient(gcs.Config{
		ServiceAccountFile: connection.ServiceAccountFile,
		ServiceAccountJSON: connection.ServiceAccountJSON,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.GCS[connection.Name] = client

	return nil
}

func (m *Manager) AddPersonioConnectionFromConfig(connection *config.PersonioConnection) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.Personio == nil {
		m.Personio = make(map[string]*personio.Client)
	}

	m.Personio[connection.Name] = personio.NewClient(personio.Config{
		ClientID:     connection.ClientID,
		ClientSecret: connection.ClientSecret,
	})
	return nil
}

func (m *Manager) AddApplovinMaxConnectionFromConfig(connection *config.ApplovinMaxConnection) error {
	m.mutex.Lock()
	if m.ApplovinMax == nil {
		m.ApplovinMax = make(map[string]*applovinmax.Client)
	}
	m.mutex.Unlock()

	client, err := applovinmax.NewClient(applovinmax.Config{
		APIKey: connection.APIKey,
	})
	if err != nil {
		return err
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ApplovinMax[connection.Name] = client

	return nil
}

func (m *Manager) AddAppLovinConnectionFromConfig(connection *config.AppLovinConnection) error {
	m.mutex.Lock()
	if m.AppLovin == nil {
		m.AppLovin = make(map[string]*applovin.Client)
	}
	m.mutex.Unlock()

	client, err := applovin.NewClient(applovin.Config{
		APIKey: connection.APIKey,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.AppLovin[connection.Name] = client

	return nil
}

func (m *Manager) AddPipedriveConnectionFromConfig(connection *config.PipedriveConnection) error {
	m.mutex.Lock()
	if m.Pipedrive == nil {
		m.Pipedrive = make(map[string]*pipedrive.Client)
	}
	m.mutex.Unlock()

	client, err := pipedrive.NewClient(pipedrive.Config{
		APIToken: connection.APIToken,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Pipedrive[connection.Name] = client
	return nil
}

func (m *Manager) AddQuickBooksConnectionFromConfig(connection *config.QuickBooksConnection) error {
	m.mutex.Lock()
	if m.QuickBooks == nil {
		m.QuickBooks = make(map[string]*quickbooks.Client)
	}
	m.mutex.Unlock()

	client, err := quickbooks.NewClient(quickbooks.Config{
		CompanyID:    connection.CompanyID,
		ClientID:     connection.ClientID,
		ClientSecret: connection.ClientSecret,
		RefreshToken: connection.RefreshToken,
		Environment:  connection.Environment,
		MinorVersion: connection.MinorVersion,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.QuickBooks[connection.Name] = client
	return nil
}

func (m *Manager) AddISOCPulseConnectionFromConfig(connection *config.ISOCPulseConnection) error {
	m.mutex.Lock()
	if m.ISOCPulse == nil {
		m.ISOCPulse = make(map[string]*isocpulse.Client)
	}
	m.mutex.Unlock()

	client, err := isocpulse.NewClient(isocpulse.Config{
		Token: connection.Token,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.ISOCPulse[connection.Name] = client
	return nil
}

func (m *Manager) AddZoomConnectionFromConfig(connection *config.ZoomConnection) error {
	m.mutex.Lock()
	if m.Zoom == nil {
		m.Zoom = make(map[string]*zoom.Client)
	}
	m.mutex.Unlock()

	client, err := zoom.NewClient(zoom.Config{
		ClientID:     connection.ClientID,
		ClientSecret: connection.ClientSecret,
		AccountID:    connection.AccountID,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Zoom[connection.Name] = client
	return nil
}

func (m *Manager) AddMixpanelConnectionFromConfig(connection *config.MixpanelConnection) error {
	m.mutex.Lock()
	if m.Mixpanel == nil {
		m.Mixpanel = make(map[string]*mixpanel.Client)
	}
	m.mutex.Unlock()

	client, err := mixpanel.NewClient(mixpanel.Config{
		Username:  connection.Username,
		Password:  connection.Password,
		ProjectID: connection.ProjectID,
		Server:    connection.Server,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Mixpanel[connection.Name] = client
	return nil
}

func (m *Manager) AddGoogleAnalyticsConnectionFromConfig(connection *config.GoogleAnalyticsConnection) error {
	m.mutex.Lock()
	if m.GoogleAnalytics == nil {
		m.GoogleAnalytics = make(map[string]*googleanalytics.Client)
	}
	m.mutex.Unlock()

	client, err := googleanalytics.NewClient(googleanalytics.Config{
		ServiceAccountFile: connection.ServiceAccountFile,
		PropertyID:         connection.PropertyID,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.GoogleAnalytics[connection.Name] = client
	return nil
}

func (m *Manager) AddSalesforceConnectionFromConfig(connection *config.SalesforceConnection) error {
	m.mutex.Lock()
	if m.Salesforce == nil {
		m.Salesforce = make(map[string]*salesforce.Client)
	}
	m.mutex.Unlock()

	client, err := salesforce.NewClient(salesforce.Config{
		Username: connection.Username,
		Password: connection.Password,
		Token:    connection.Token,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Salesforce[connection.Name] = client
	return nil
}

func (m *Manager) AddSQLiteConnectionFromConfig(connection *config.SQLiteConnection) error {
	m.mutex.Lock()
	if m.SQLite == nil {
		m.SQLite = make(map[string]*sqlite.Client)
	}
	m.mutex.Unlock()

	client, err := sqlite.NewClient(sqlite.Config{
		Path: connection.Path,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.SQLite[connection.Name] = client
	return nil
}

func (m *Manager) AddOracleConnectionFromConfig(connection *config.OracleConnection) error {
	m.mutex.Lock()
	if m.Oracle == nil {
		m.Oracle = make(map[string]*oracle.Client)
	}
	m.mutex.Unlock()

	client, err := oracle.NewClient(oracle.Config{
		Username: connection.Username,
		Password: connection.Password,
		Host:     connection.Host,
		Port:     connection.Port,
		DBName:   connection.DBName,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Oracle[connection.Name] = client
	return nil
}

func (m *Manager) AddKinesisConnectionFromConfig(connection *config.KinesisConnection) error {
	m.mutex.Lock()
	if m.Kinesis == nil {
		m.Kinesis = make(map[string]*kinesis.Client)
	}
	m.mutex.Unlock()

	client, err := kinesis.NewClient(kinesis.Config{
		AccessKeyID:     connection.AccessKeyID,
		SecretAccessKey: connection.SecretAccessKey,
		Region:          connection.Region,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Kinesis[connection.Name] = client
	return nil
}

func (m *Manager) AddDB2ConnectionFromConfig(connection *config.DB2Connection) error {
	m.mutex.Lock()
	if m.DB2 == nil {
		m.DB2 = make(map[string]*db2.Client)
	}
	m.mutex.Unlock()

	client, err := db2.NewClient(db2.Config{
		Username: connection.Username,
		Password: connection.Password,
		Host:     connection.Host,
		Port:     connection.Port,
		Database: connection.Database,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.DB2[connection.Name] = client
	return nil
}

func (m *Manager) AddPhantombusterConnectionFromConfig(connection *config.PhantombusterConnection) error {
	m.mutex.Lock()
	if m.Phantombuster == nil {
		m.Phantombuster = make(map[string]*phantombuster.Client)
	}
	m.mutex.Unlock()

	client, err := phantombuster.NewClient(phantombuster.Config{
		APIKey: connection.APIKey,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Phantombuster[connection.Name] = client
	return nil
}

func (m *Manager) AddAttioConnectionFromConfig(connection *config.AttioConnection) error {
	m.mutex.Lock()
	if m.Attio == nil {
		m.Attio = make(map[string]*attio.Client)
	}
	m.mutex.Unlock()

	client, err := attio.NewClient(attio.Config{
		APIKey: connection.APIKey,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Attio[connection.Name] = client
	return nil
}

func (m *Manager) AddElasticsearchConnectionFromConfig(connection *config.ElasticsearchConnection) error {
	m.mutex.Lock()
	if m.Elasticsearch == nil {
		m.Elasticsearch = make(map[string]*elasticsearch.Client)
	}
	m.mutex.Unlock()

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Username:    connection.Username,
		Password:    connection.Password,
		Host:        connection.Host,
		Port:        connection.Port,
		Secure:      connection.Secure,
		VerifyCerts: connection.VerifyCerts,
	})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Elasticsearch[connection.Name] = client
	return nil
}

func (m *Manager) AddFrankfurterConnectionFromConfig(connection *config.FrankfurterConnection) error {
	m.mutex.Lock()
	if m.Frankfurter == nil {
		m.Frankfurter = make(map[string]*frankfurter.Client)
	}
	m.mutex.Unlock()

	client, err := frankfurter.NewClient(frankfurter.Config{})
	if err != nil {
		return err
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.Frankfurter[connection.Name] = client
	m.Frankfurter[connection.Name] = client
	return nil
}

func (m *Manager) AddEMRServerlessConnectionFromConfig(connection *config.EMRServerlessConnection) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.EMRSeverless == nil {
		m.EMRSeverless = make(map[string]*emr_serverless.Client)
	}
	client, err := emr_serverless.NewClient(emr_serverless.Config{
		AccessKey:     connection.AccessKey,
		SecretKey:     connection.SecretKey,
		ApplicationID: connection.ApplicationID,
		ExecutionRole: connection.ExecutionRole,
		Region:        connection.Region,
		Workspace:     connection.Workspace,
	})
	if err != nil {
		return err
	}
	m.EMRSeverless[connection.Name] = client

	return nil
}

var envVarRegex = regexp.MustCompile(`\${([^}]+)}`)

func processConnections[T config.Named](connections []T, adder func(*T) error, wg *conc.WaitGroup, errList *[]error, mu *sync.Mutex) {
	if connections == nil {
		return
	}
	for i := range connections {
		conn := &connections[i]
		wg.Go(func() {
			// Check for environment variable placeholders in connection fields
			v := reflect.ValueOf(conn).Elem()
			for i := range v.NumField() {
				field := v.Field(i)

				// Only process string fields
				if field.Kind() != reflect.String {
					continue
				}
				if !field.CanSet() {
					continue
				}
				strValue := strings.TrimSpace(field.String())
				matches := envVarRegex.FindStringSubmatch(strValue)
				for len(matches) > 0 {
					envVarName := matches[1]
					envValue := os.Getenv(envVarName)
					strValue = strings.Replace(strValue, matches[0], envValue, 1)
					field.SetString(strValue)
					// Look for more matches after replacement
					matches = envVarRegex.FindStringSubmatch(strValue)
				}
			}

			err := adder(conn)
			if err != nil {
				mu.Lock()
				*errList = append(*errList, errors.Wrapf(err, "failed to add connection %q", (*conn).GetName()))
				mu.Unlock()
			}
		})
	}
}

func NewManagerFromConfig(cm *config.Config) (*Manager, []error) {
	connectionManager := &Manager{}

	var wg conc.WaitGroup
	var errList []error
	var mu sync.Mutex

	processConnections(cm.SelectedEnvironment.Connections.AthenaConnection, connectionManager.AddAthenaConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.GoogleCloudPlatform, connectionManager.AddBqConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Snowflake, connectionManager.AddSfConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Postgres, connectionManager.AddPgConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.RedShift, connectionManager.AddRedshiftConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.MsSQL, connectionManager.AddMsSQLConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Databricks, connectionManager.AddDatabricksConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Synapse, connectionManager.AddSynapseSQLConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Mongo, connectionManager.AddMongoConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.MySQL, connectionManager.AddMySQLConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Notion, connectionManager.AddNotionConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Shopify, connectionManager.AddShopifyConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Gorgias, connectionManager.AddGorgiasConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Klaviyo, connectionManager.AddKlaviyoConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Adjust, connectionManager.AddAdjustConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.FacebookAds, connectionManager.AddFacebookAdsConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Stripe, connectionManager.AddStripeConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Appsflyer, connectionManager.AddAppsflyerConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Kafka, connectionManager.AddKafkaConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.GoogleSheets, connectionManager.AddGoogleSheetsConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.DuckDB, connectionManager.AddDuckDBConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.ClickHouse, connectionManager.AddClickHouseConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Hubspot, connectionManager.AddHubspotConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Chess, connectionManager.AddChessConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Airtable, connectionManager.AddAirtableConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.S3, connectionManager.AddS3ConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Slack, connectionManager.AddSlackConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Asana, connectionManager.AddAsanaConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.DynamoDB, connectionManager.AddDynamoDBConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Zendesk, connectionManager.AddZendeskConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.GoogleAds, connectionManager.AddGoogleAdsConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.TikTokAds, connectionManager.AddTikTokAdsConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.GitHub, connectionManager.AddGitHubConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.AppStore, connectionManager.AddAppStoreConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.LinkedInAds, connectionManager.AddLinkedInAdsConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.GCS, connectionManager.AddGCSConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Personio, connectionManager.AddPersonioConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.ApplovinMax, connectionManager.AddApplovinMaxConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Kinesis, connectionManager.AddKinesisConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Pipedrive, connectionManager.AddPipedriveConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Mixpanel, connectionManager.AddMixpanelConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.QuickBooks, connectionManager.AddQuickBooksConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Zoom, connectionManager.AddZoomConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.EMRServerless, connectionManager.AddEMRServerlessConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.GoogleAnalytics, connectionManager.AddGoogleAnalyticsConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.AppLovin, connectionManager.AddAppLovinConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Frankfurter, connectionManager.AddFrankfurterConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Salesforce, connectionManager.AddSalesforceConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.SQLite, connectionManager.AddSQLiteConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Oracle, connectionManager.AddOracleConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.DB2, connectionManager.AddDB2ConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Phantombuster, connectionManager.AddPhantombusterConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Elasticsearch, connectionManager.AddElasticsearchConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Spanner, connectionManager.AddSpannerConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Solidgate, connectionManager.AddSolidgateConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Smartsheet, connectionManager.AddSmartsheetConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Attio, connectionManager.AddAttioConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.Sftp, connectionManager.AddSftpConnectionFromConfig, &wg, &errList, &mu)
	processConnections(cm.SelectedEnvironment.Connections.ISOCPulse, connectionManager.AddISOCPulseConnectionFromConfig, &wg, &errList, &mu)
	wg.Wait()
	return connectionManager, errList
}
