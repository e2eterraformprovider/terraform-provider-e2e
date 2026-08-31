---
page_title: "e2e_scaler_group Data Source - terraform-provider-e2e"
subcategory: ""
description: |-
  Retrieves details of an existing E2E Cloud Auto Scaling (Scaler) Group.
---

# e2e_scaler_group (Data Source)

Provides information about an existing E2E Cloud Auto Scaling Group (Scaler Group).

```hcl
data "e2e_scaler_group" "example" {
  id         = "your-scaler-group-id"    # Replace with your scaler group ID
  project_id = "your-project-id"         # Replace with your actual project ID
  location   = "your-region"             # Replace with the region, e.g., "us-west-1"
}
```

## Schema

### Required

- `id` (String) The unique ID of the scaler group.
- `project_id` (String) The project ID under which the scaler group is created.
- `location` (String) The region where the scaler group is deployed.

### Read-Only

- `name` (String) Name of the scaler group.
- `desired` (Number) Desired number of nodes.
- `min_nodes` (Number) Minimum allowed nodes.
- `max_nodes` (Number) Maximum allowed nodes.
- `plan_name` (String) Name of the plan used by the scaler group.
- `vm_image_name` (String) Name of the VM image used by nodes.
- `provision_status` (String) Current provision status (e.g., "Running", "Stopped").
- `policy_type` (String) Type of scaling policy in use.
- `policy` (List) List of elastic scaling policies with fields:
  - `type` (String)
  - `adjust` (Number)
  - `parameter` (String)
  - `operator` (String)
  - `value` (String)
  - `period_number` (String)
  - `period_seconds` (String)
  - `cooldown` (String)
- `scheduled_policy` (List) List of scheduled scaling policies with fields:
  - `type` (String)
  - `adjust` (String)
  - `recurrence` (String)
