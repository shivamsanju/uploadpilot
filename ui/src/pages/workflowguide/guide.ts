export const workflowGuideMd = `
The guide describes how to define the workflows using yaml

----

### Workflow components

#### 1. Workflow Variables

- \`variables\`: Key-value map containing workflow variables.

#### 2. Root Statement

Defines the main execution path of the workflow.

- \`sequence\`: Executes activities in order.
- \`parallel\`: Executes activities in parallel.
- \`activity\`: Executes a single activity.

#### 3. Activity Invocation

Each \`activity\` includes:

- \`key\`: Unique key for the activity.
- \`uses\`: The activity name to execute.
- \`with\`: Input parameters.
- \`input\`: key of the activity to be used as input.
- \`save_output\`: Whether to save the output.
- \`schedule_to_close_timeout_seconds\`: Timeout from scheduling to completion.
- \`start_to_close_timeout_seconds\`: Execution timeout.
- \`max_retries\`: Number of retries.
- \`retry_backoff_coefficient\`: Backoff coefficient for retries.
- \`retry_max_interval_seconds\`: Maximum retry interval.
- \`retry_initial_interval_seconds\`: Initial retry interval.
- \`on_success\`: Defines an action on success.
- \`on_error\`: Defines an action on failure.

#### 4. Workflow error and success handlers

- \`on_workflow_success\`: Defines steps to execute if the workflow succeeds.
- \`on_workflow_failure\`: Defines steps to execute if the workflow fails.

----

### Important Points

1. The keys should be unique within the workflow.
2. If \`save_output\` is not specified, the output will not be saved.
3. If \`input\` is not specified, the uploaded file will be used as the input.
4. If \`input\` is specified, it output of input activity will be used as the input.
5. Inside \`on_workflow_failure\`, extra variable named \`workflow_error\` will be available.

----

### Workflow Structure

A workflow YAML file should be structured as follows:

\`\`\`yaml
variables:
  key1: value1
  key2: value2

root:
  sequence:
    elements:
      - activity:
          key: activity-key
          uses: activity-name
          with:
            param1: value1
            param2: value2
          input: input-value
          save_output: true
          schedule_to_close_timeout_seconds: 300
          start_to_close_timeout_seconds: 200
          max_retries: 3
          retry_backoff_coefficient: 2.0
          retry_max_interval_seconds: 60
          retry_initial_interval_seconds: 5
          on_success:
            activity:
              key: success-key
              uses: success-activity
          on_error:
            activity:
              key: error-key
              uses: error-activity

on_workflow_success:
  activity:
    key: success-handler
    uses: workflow-success-handler

on_workflow_failure:
  activity:
    key: failure-handler
    uses: workflow-failure-handler
\`\`\`

----

### Example Use Cases

#### 1. Simple Workflow with a Single Activity

\`\`\`yaml
root:
  activity:
    key: converttojpeg
    uses: ImageFormatConverterV1
    with:
      format: jpeg
    save_output: true
\`\`\`

#### 2. Sequential Execution

\`\`\`yaml
root:
  sequence:
    elements:
      - activity:
          key: step1
          uses: ProcessStep1
      - activity:
          key: step2
          uses: ProcessStep2
\`\`\`

#### 3. Parallel Execution

\`\`\`yaml
root:
  parallel:
    branches:
      - activity:
          key: branch1
          uses: TaskA
      - activity:
          key: branch2
          uses: TaskB
\`\`\`


`;
