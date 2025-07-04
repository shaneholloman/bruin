package clickhouse

import (
	"context"
	"fmt"
	"sort"

	click_house "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/bruin-data/bruin/pkg/ansisql"
	"github.com/bruin-data/bruin/pkg/query"
	"github.com/pkg/errors"
)

// Rowscanner exists since clickhouse library requires us to scan either to a specific type or an implementor of the
// interface sql.Scanner, cannot scan directly to interface{}.
type RowScanner struct {
	values []any
}

func (s *RowScanner) SetValues(values []any) {
	s.values = values
}

func (s *RowScanner) Scan(src any) error {
	s.values = append(s.values, src)
	return nil
}

type Client struct {
	connection connection
	config     ClickHouseConfig
}

type ClickHouseConfig interface {
	ToClickHouseOptions() *click_house.Options
	GetIngestrURI() string
}

type connection interface {
	Query(ctx context.Context, sql string, args ...any) (driver.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) error
}

func NewClient(c ClickHouseConfig) (*Client, error) {
	conn, err := click_house.Open(c.ToClickHouseOptions())
	if err != nil {
		return nil, err
	}

	return &Client{connection: conn, config: c}, nil
}

func (c *Client) RunQueryWithoutResult(ctx context.Context, query *query.Query) error {
	err := c.connection.Exec(ctx, query.String())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetIngestrURI() (string, error) {
	return c.config.GetIngestrURI(), nil
}

// Select runs a query and returns the results.
func (c *Client) Select(ctx context.Context, query *query.Query) ([][]interface{}, error) {
	rows, err := c.connection.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	collectedRows := make([][]interface{}, 0)
	for rows.Next() {
		result := RowScanner{}
		if err := rows.Scan(&result); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}

		collectedRows = append(collectedRows, result.values)
	}

	return collectedRows, nil
}

func (c *Client) SelectWithSchema(ctx context.Context, queryObj *query.Query) (*query.QueryResult, error) {
	rows, err := c.connection.Query(ctx, queryObj.String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute query")
	}
	defer rows.Close()

	fieldDescriptions := rows.ColumnTypes()
	if fieldDescriptions == nil {
		return nil, errors.New("field descriptions are not available")
	}

	columns := make([]string, len(fieldDescriptions))
	columnTypes := make([]string, len(fieldDescriptions))
	for i, field := range fieldDescriptions {
		columns[i] = field.Name()
		columnTypes[i] = field.DatabaseTypeName()
	}

	collectedRows := make([][]interface{}, 0)
	for rows.Next() {
		result := RowScanner{}
		if err := rows.Scan(&result); err != nil {
			return nil, errors.Wrap(err, "failed to scan row")
		}
		collectedRows = append(collectedRows, result.values)
	}

	return &query.QueryResult{
		Columns:     columns,
		ColumnTypes: columnTypes,
		Rows:        collectedRows,
	}, nil
}

// Test runs a simple query (SELECT 1) to validate the connection.
func (c *Client) Ping(ctx context.Context) error {
	q := query.Query{
		Query: "SELECT 1",
	}
	err := c.RunQueryWithoutResult(ctx, &q)
	if err != nil {
		return errors.Wrap(err, "failed to run test query on ClickHouse connection")
	}

	return nil
}

func (c *Client) GetDatabaseSummary(ctx context.Context) (*ansisql.DBDatabase, error) {
	// ClickHouse has databases and tables
	// We'll query system.tables to get all databases and tables
	q := `
SELECT
    database,
    name as table_name
FROM
    system.tables
WHERE
    database NOT IN ('system', 'information_schema', 'INFORMATION_SCHEMA')
ORDER BY database, name;
`

	result, err := c.Select(ctx, &query.Query{Query: q})
	if err != nil {
		return nil, fmt.Errorf("failed to query ClickHouse system tables: %w", err)
	}

	summary := &ansisql.DBDatabase{
		Name:    "clickhouse", // ClickHouse instance
		Schemas: []*ansisql.DBSchema{},
	}
	schemas := make(map[string]*ansisql.DBSchema)

	for _, row := range result {
		if len(row) != 2 {
			continue
		}

		schemaName, ok := row[0].(string)
		if !ok {
			continue
		}
		tableName, ok := row[1].(string)
		if !ok {
			continue
		}

		// Create schema if it doesn't exist
		if _, exists := schemas[schemaName]; !exists {
			schema := &ansisql.DBSchema{
				Name:   schemaName,
				Tables: []*ansisql.DBTable{},
			}
			schemas[schemaName] = schema
		}

		// Add table to schema
		table := &ansisql.DBTable{
			Name:    tableName,
			Columns: []*ansisql.DBColumn{}, // Initialize empty columns array
		}
		schemas[schemaName].Tables = append(schemas[schemaName].Tables, table)
	}

	for _, schema := range schemas {
		summary.Schemas = append(summary.Schemas, schema)
	}

	// Sort schemas by name
	sort.Slice(summary.Schemas, func(i, j int) bool {
		return summary.Schemas[i].Name < summary.Schemas[j].Name
	})

	return summary, nil
}
