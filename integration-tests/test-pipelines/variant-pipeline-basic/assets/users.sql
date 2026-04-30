/* @bruin
name: "{{ var.client }}_users"
type: duckdb.sql
depends:
  - "{{ var.client }}_raw"
@bruin */
select * from {{ var.client }}_raw;
