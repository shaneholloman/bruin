# AWS Redshift

Bruin supports AWS Redshift as a data platform, which means you can use Bruin to build tables and views in your Redshift data warehouse.

## Connection
In order to set up a Redshift connection, you need to add a configuration item to `connections` in the `.bruin.yml` file complying with the following schema
Mind that, despite the connection being at all effects a Postgres connection, the default `port` field of Amazon Redshift is `5439`.

```yaml
    connections:
      redshift:
        - name: "connection_name"
          username: "awsuser"
          password: "XXXXXXXXXX"
          host: "redshift-cluster-1.xxxxxxxxx.eu-north-1.redshift.amazonaws.com"
          port: 5439
          database: "dev"
          ssl_mode: "allow" # optional
```

> [!NOTE]
> `ssl_mode` should be one of the modes describe in the [PostgreSQL documentation](https://www.postgresql.org/docs/current/libpq-ssl.html#LIBPQ-SSL-PROTECTION).


### Making Redshift publicly accessible

Before the connection works properly, you need to ensure that the Redshift cluster can be accessed from the outside. In order to do that you must mark the configuration option in your Redshift cluster

![Make publicly available](/publicly-accessible.png)

In addition to this, you must configure the inbound rules of the security group your Redshift cluster belongs to, to accept inbound connections. In the example below we enabled access for all origins but you can set more restrictive rules for this.

![Inbound Rules](/inbound-rules.png)

If you have trouble setting this up you can check [AWS documentation](https://repost.aws/knowledge-center/redshift-cluster-private-public) on the topic


## AWS Redshift Assets

### `rs.sql`
Runs a materialized AWS Redshift asset or an SQL script. For detailed parameters, you can check [Definition Schema](../assets/definition-schema.md) page.

#### Example: Create a table for product reviews
```bruin-sql
/* @bruin
name: product_reviews.table
type: rs.sql
materialization:
    type: table
@bruin */

create table product_reviews (
    review_id bigint identity(1,1),
    product_id bigint,
    user_id bigint,
    rating int,
    review_text varchar(500),
    review_date timestamp
);
```

#### Example: Run an AWS Redshift script to clean up old data
```bruin-sql
/* @bruin
name: clean_old_data
type: rs.sql
@bruin */

begin transaction;

delete from user_activity
where activity_date < dateadd(year, -2, current_date);

delete from order_history
where order_date < dateadd(year, -5, current_date);

commit transaction;
```

### `rs.sensor.query`

Checks if a query returns any results in Redshift, runs every 5 minutes until this query returns any results.

```yaml
name: string
type: string
parameters:
    query: string
```

**Parameters**:
- `query`: Query you expect to return any results

#### Example: Partitioned upstream table

Checks if the data available in upstream table for end date of the run.
```yaml
name: analytics_123456789.events
type: rs.sensor.query
parameters:
    query: select exists(select 1 from upstream_table where dt = "{{ end_date }}"
```

#### Example: Streaming upstream table

Checks if there is any data after end timestamp, by assuming that older data is not appended to the table.
```yaml
name: analytics_123456789.events
type: rs.sensor.query
parameters:
    query: select exists(select 1 from upstream_table where inserted_at > "{{ end_timestamp }}"
```

### `rs.seed`
`rs.seed` is a special type of asset used to represent CSV files that contain data that is prepared outside of your pipeline that will be loaded into your Redshift database. Bruin supports seed assets natively, allowing you to simply drop a CSV file in your pipeline and ensuring the data is loaded to the Redshift database.

You can define seed assets in a file ending with `.yaml`:
```yaml
name: dashboard.hello
type: rs.seed

parameters:
    path: seed.csv
```

**Parameters**:
- `path`:  The `path` parameter is the path to the CSV file that will be loaded into the data platform. path is relative to the asset definition file.


####  Examples: Load csv into a Redshift database

The examples below show how to load a CSV into a Redshift database.
```yaml
name: dashboard.hello
type: rs.seed

parameters:
    path: seed.csv
```

Example CSV:

```csv
name,networking_through,position,contact_date
Y,LinkedIn,SDE,2024-01-01
B,LinkedIn,SDE 2,2024-01-01
```
