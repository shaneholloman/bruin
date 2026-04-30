/* @bruin
name: "{{ var.client }}_raw"
type: duckdb.sql
@bruin */
select '{{ var.client }}' as client, '{{ var.region }}' as region;
