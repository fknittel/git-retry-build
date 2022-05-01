CREATE OR REPLACE VIEW `APP_ID.PROJECT_NAME.step_status_transitions`
AS
/*
Step status transition table.
This view represents status transitions for build steps over time.
Each row represents a build where the step produced a different status
output that it did in the previous run (on that builder etc).
*/
WITH
  step_lag AS (
  SELECT
    b.end_time AS end_time,
    b.builder.project AS project,
    b.builder.bucket AS bucket,
    b.builder.builder AS builder,
    JSON_EXTRACT_SCALAR(b.input.properties, "$.builder_group") as buildergroup,
    b.id,
    b.number,
    b.critical as critical,
    b.status as status,
    output.gitiles_commit as output_commit,
    input.gitiles_commit as input_commit,
    step.name AS step_name,
    step.status AS step_status,
    LAG(step.status) OVER (PARTITION BY b.builder.project, b.builder.bucket, b.builder.builder, b.output.gitiles_commit.host, b.output.gitiles_commit.project, b.output.gitiles_commit.ref, step.name ORDER BY b.output.gitiles_commit.position, b.id desc, b.number) AS previous_step_status,
    LAG(b.output.gitiles_commit) OVER (PARTITION BY b.builder.project, b.builder.bucket, b.builder.builder, b.output.gitiles_commit.host, b.output.gitiles_commit.project, b.output.gitiles_commit.ref, step.name ORDER BY b.output.gitiles_commit.position, b.id desc, b.number) AS previous_output_commit,
     LAG(b.input.gitiles_commit) OVER (PARTITION BY b.builder.project, b.builder.bucket, b.builder.builder, b.input.gitiles_commit.host, b.input.gitiles_commit.project, b.input.gitiles_commit.ref, step.name ORDER BY b.input.gitiles_commit.position, b.id desc, b.number) AS previous_input_commit,
    LAG(b.id) OVER (PARTITION BY b.builder.project, b.builder.bucket, b.builder.builder, b.output.gitiles_commit.host, b.output.gitiles_commit.project, b.output.gitiles_commit.ref, step.name ORDER BY b.output.gitiles_commit.position, b.number) AS previous_id
    FROM
    `cr-buildbucket.raw.completed_builds_prod` b,
    UNNEST(steps) AS step
  WHERE
    PROJECT_FILTER_CONDITIONS
)
SELECT
  end_time,
  project,
  bucket,
  builder,
  buildergroup,
  number,
  id,
  critical,
  status,
  output_commit,
  input_commit,
  step_name,
  step_status,
  previous_output_commit,
  previous_input_commit,
  previous_id,
  previous_step_status
FROM
  step_lag s
WHERE
  s.previous_step_status != s.step_status
