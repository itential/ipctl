# API Command Reference

The `api` command provides direct access to the Itential Platform REST API, allowing you to perform arbitrary HTTP requests against any API endpoint. This is useful for accessing endpoints that don't have dedicated commands or for advanced automation scenarios.

## Overview

The `api` command supports all standard HTTP methods:
- `get` - Retrieve resources
- `post` - Create new resources
- `put` - Update existing resources
- `patch` - Partially update resources
- `delete` - Remove resources

All `api` commands return JSON responses that can be formatted using the standard output options.

## Basic Usage

The basic syntax for the `api` command is:

```bash
ipctl api <method> <path> [options]
```

Where:
- `<method>` is the HTTP method (get, post, put, patch, delete)
- `<path>` is the API endpoint path (e.g., `/api/v2.0/workflows`)
- `[options]` are additional command-line options

## Command Options

### Common Options

All `api` commands support the following options:

#### `--expected-status-code <code>`

Specify the expected HTTP status code for the response. If the response status code doesn't match, the command will return an error.

```bash
ipctl api get /api/v2.0/workflows --expected-status-code 200
```

#### `--params <key=value>`

Pass custom query parameters to the API request. This option can be specified multiple times to add multiple parameters. Parameters are automatically URL-encoded.

```bash
ipctl api get /api/v2.0/workflows --params limit=10 --params offset=20
```

The `--params` flag supports:
- Special characters and URLs (automatic URL encoding)
- Empty values (`key=`)
- Values containing equals signs (`key=value=with=equals`)
- Unicode characters
- Spaces in values (`message=hello world`)

### Method-Specific Options

#### `--data <data>` or `-d <data>`

For `post`, `put`, and `patch` commands, specify the request body data. The data can be:
- Inline JSON: `-d '{"name": "My Workflow"}'`
- File reference: `-d @workflow.json` (reads from file)

```bash
ipctl api post /api/v2.0/workflows -d '{"name": "Test Workflow"}'
ipctl api put /api/v2.0/workflows/123 -d @workflow.json
```

## HTTP Methods

### GET - Retrieve Resources

Use `get` to retrieve resources from the API.

**Basic Usage:**
```bash
ipctl api get /api/v2.0/workflows
```

**With Query Parameters:**
```bash
# Pagination
ipctl api get /api/v2.0/workflows --params limit=50 --params offset=100

# Filtering
ipctl api get /api/v2.0/automations --params filter="name eq 'test'"

# Sorting
ipctl api get /api/v2.0/projects --params sort=-createdAt

# Multiple parameters
ipctl api get /api/v2.0/workflows \
  --params limit=25 \
  --params offset=0 \
  --params status=active \
  --params type=workflow
```

**With Status Code Validation:**
```bash
ipctl api get /api/v2.0/workflows --expected-status-code 200
```

### POST - Create Resources

Use `post` to create new resources.

**Inline JSON Data:**
```bash
ipctl api post /api/v2.0/workflows -d '{
  "name": "New Workflow",
  "description": "Created via API",
  "type": "automation"
}'
```

**From File:**
```bash
ipctl api post /api/v2.0/workflows -d @new-workflow.json
```

**With Query Parameters:**
```bash
ipctl api post /api/v2.0/workflows \
  -d @workflow.json \
  --params validate=true \
  --params async=false
```

**With Expected Status:**
```bash
ipctl api post /api/v2.0/workflows \
  -d @workflow.json \
  --expected-status-code 201
```

### PUT - Update Resources

Use `put` to completely update an existing resource.

**Basic Update:**
```bash
ipctl api put /api/v2.0/workflows/507f1f77bcf86cd799439011 -d '{
  "name": "Updated Workflow",
  "description": "Modified via API"
}'
```

**From File:**
```bash
ipctl api put /api/v2.0/workflows/507f1f77bcf86cd799439011 -d @workflow.json
```

**With Query Parameters:**
```bash
ipctl api put /api/v2.0/workflows/507f1f77bcf86cd799439011 \
  -d @workflow.json \
  --params force=true
```

### PATCH - Partially Update Resources

Use `patch` to partially update a resource (only specified fields are modified).

**Partial Update:**
```bash
ipctl api patch /api/v2.0/workflows/507f1f77bcf86cd799439011 -d '{
  "description": "Updated description only"
}'
```

**With Query Parameters:**
```bash
ipctl api patch /api/v2.0/workflows/507f1f77bcf86cd799439011 \
  -d '{"status": "active"}' \
  --params partial=true
```

### DELETE - Remove Resources

Use `delete` to remove resources.

**Basic Delete:**
```bash
ipctl api delete /api/v2.0/workflows/507f1f77bcf86cd799439011
```

**With Query Parameters:**
```bash
ipctl api delete /api/v2.0/workflows/507f1f77bcf86cd799439011 \
  --params cascade=true \
  --params force=true
```

**With Expected Status:**
```bash
ipctl api delete /api/v2.0/workflows/507f1f77bcf86cd799439011 \
  --expected-status-code 204
```

## Working with Query Parameters

Query parameters are specified using the `--params` flag in `key=value` format. This flag can be used multiple times to add multiple parameters.

### Basic Parameters

```bash
ipctl api get /api/v2.0/workflows --params limit=10
```

### Multiple Parameters

```bash
ipctl api get /api/v2.0/workflows \
  --params limit=25 \
  --params offset=50 \
  --params status=active
```

### Complex Parameters

**Filtering:**
```bash
ipctl api get /api/v2.0/automations \
  --params filter="name eq 'Production Workflow'"
```

**Sorting:**
```bash
ipctl api get /api/v2.0/projects \
  --params sort=-createdAt \
  --params limit=10
```

**URL Parameters:**
```bash
ipctl api get /api/v2.0/search \
  --params url=https://example.com/api \
  --params callback=https://callback.example.com
```

**Special Characters:**
All special characters are automatically URL-encoded, so you can pass them directly:

```bash
ipctl api get /api/v2.0/search --params query="hello world & goodbye"
```

### Combining with Existing Query Strings

If your API path already includes query parameters, the `--params` flag will append to them:

```bash
# URL already has ?type=workflow
ipctl api get "/api/v2.0/workflows?type=workflow" --params limit=10

# Results in: /api/v2.0/workflows?type=workflow&limit=10
```

## Response Handling

### Default JSON Output

By default, responses are returned as formatted JSON:

```bash
ipctl api get /api/v2.0/workflows
```

Output:
```json
{
  "workflows": [
    {
      "id": "507f1f77bcf86cd799439011",
      "name": "Production Workflow",
      "type": "automation"
    }
  ],
  "total": 1
}
```

### Output Formatting

While the `api` command returns raw JSON responses, you can use `jq` or other tools for additional formatting:

```bash
# Extract specific fields
ipctl api get /api/v2.0/workflows | jq '.workflows[].name'

# Pretty print with jq
ipctl api get /api/v2.0/workflows | jq '.'

# Filter results
ipctl api get /api/v2.0/workflows | jq '.workflows[] | select(.type == "automation")'
```

## Error Handling

### Status Code Validation

Use `--expected-status-code` to validate the response:

```bash
ipctl api get /api/v2.0/workflows --expected-status-code 200
```

If the response status code doesn't match, the command exits with an error.

### Common Status Codes

- `200 OK` - Successful GET, PUT, PATCH
- `201 Created` - Successful POST
- `204 No Content` - Successful DELETE
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource doesn't exist
- `500 Internal Server Error` - Server error

## Practical Examples

### List All Workflows with Pagination

```bash
# Get first page
ipctl api get /api/v2.0/workflows \
  --params limit=50 \
  --params offset=0

# Get second page
ipctl api get /api/v2.0/workflows \
  --params limit=50 \
  --params offset=50
```

### Search for Resources

```bash
ipctl api get /api/v2.0/workflows \
  --params filter="name contains 'production'" \
  --params limit=10
```

### Create a Workflow from Template

```bash
ipctl api post /api/v2.0/workflows \
  -d @templates/workflow-template.json \
  --params validate=true \
  --expected-status-code 201
```

### Update Multiple Fields

```bash
ipctl api patch /api/v2.0/workflows/507f1f77bcf86cd799439011 -d '{
  "description": "Updated workflow description",
  "status": "active",
  "tags": ["production", "critical"]
}'
```

### Batch Delete with Parameters

```bash
# Delete with cascade to remove dependencies
ipctl api delete /api/v2.0/projects/507f1f77bcf86cd799439011 \
  --params cascade=true \
  --params force=true \
  --expected-status-code 204
```

### Export Data with Filtering

```bash
ipctl api get /api/v2.0/automations \
  --params status=active \
  --params created_after=2024-01-01 \
  --params limit=1000 > active-automations.json
```

## Best Practices

### Use Descriptive File Names

When using `-d @file.json`, use descriptive names:

```bash
ipctl api post /api/v2.0/workflows -d @create-production-workflow.json
```

### Validate Responses

Always use `--expected-status-code` for automation scripts:

```bash
ipctl api post /api/v2.0/workflows \
  -d @workflow.json \
  --expected-status-code 201
```

### Use Query Parameters for Filtering

Instead of filtering in code, use API query parameters:

```bash
# Good - filter at API level
ipctl api get /api/v2.0/workflows \
  --params status=active \
  --params type=automation

# Less efficient - retrieve all and filter client-side
ipctl api get /api/v2.0/workflows | jq '.workflows[] | select(.status == "active")'
```

### Store JSON Data in Files

For complex requests, store JSON in files:

```bash
# workflow.json
{
  "name": "Production Workflow",
  "description": "Main production workflow",
  "type": "automation",
  "steps": [...]
}

ipctl api post /api/v2.0/workflows -d @workflow.json
```

### Use Variables in Scripts

For reusable scripts, use variables:

```bash
#!/bin/bash
WORKFLOW_ID="507f1f77bcf86cd799439011"
API_PATH="/api/v2.0/workflows/${WORKFLOW_ID}"

ipctl api get "${API_PATH}" --expected-status-code 200
```

### Handle Errors Gracefully

Check exit codes in scripts:

```bash
if ipctl api get /api/v2.0/workflows --expected-status-code 200; then
  echo "Successfully retrieved workflows"
else
  echo "Failed to retrieve workflows" >&2
  exit 1
fi
```

## Advanced Usage

### Pipeline Processing

Combine with other commands for complex workflows:

```bash
# Get all workflows and extract IDs
WORKFLOW_IDS=$(ipctl api get /api/v2.0/workflows --params limit=1000 | \
  jq -r '.workflows[].id')

# Process each workflow
for id in $WORKFLOW_IDS; do
  ipctl api get "/api/v2.0/workflows/${id}" --expected-status-code 200
done
```

### Bulk Operations

```bash
# Create multiple resources from a directory
for file in workflows/*.json; do
  echo "Creating workflow from ${file}"
  ipctl api post /api/v2.0/workflows -d "@${file}" --expected-status-code 201
done
```

### Conditional Updates

```bash
# Get current resource
CURRENT=$(ipctl api get /api/v2.0/workflows/507f1f77bcf86cd799439011)

# Check condition and update
if echo "$CURRENT" | jq -e '.status == "draft"' > /dev/null; then
  ipctl api patch /api/v2.0/workflows/507f1f77bcf86cd799439011 \
    -d '{"status": "active"}' \
    --expected-status-code 200
fi
```

## Troubleshooting

### Invalid JSON Data

Ensure JSON is properly formatted:

```bash
# Validate JSON before sending
cat workflow.json | jq '.' && ipctl api post /api/v2.0/workflows -d @workflow.json
```

### URL Encoding Issues

The `--params` flag automatically handles URL encoding, but if you're including parameters directly in the path, ensure they're properly encoded:

```bash
# Good - use --params
ipctl api get /api/v2.0/search --params query="hello world"

# Manual encoding needed if in path
ipctl api get "/api/v2.0/search?query=hello%20world"
```

### Authentication Errors

Ensure your profile is configured correctly:

```bash
# Check current profile
ipctl config get profile

# Test authentication
ipctl api get /api/v2.0/health --expected-status-code 200
```

### Large Payloads

For large data files, ensure sufficient timeout:

```bash
# Configure timeout in profile
ipctl config set timeout 300  # 5 minutes

ipctl api post /api/v2.0/large-resource -d @large-file.json
```

## See Also

- [Configuration Reference](configuration-reference.md) - Configure `ipctl` settings
- [Commands Quick Reference](commands-quick-reference.md) - Overview of all commands
- [Working with Repositories](working-with-repositories.md) - Git integration
