---
page_title: "e2e_scaler_group Data Source - terraform-provider-e2e"
subcategory: "Compute"
description: |-
  Retrieves details of an existing E2E Cloud Auto Scaling (Scaler) Group.
---

<!-- updated-for-3.0-migration: 2025-11-27 -->

# e2e_scaler_group (Data Source)

Provides information about an existing E2E Cloud Auto Scaling Group (Scaler Group).

~> **Migration Notice**: This data source has parameter changes in v3.0.0. See the [v3.0.0 Upgrade Guide](../guides/upgrade-to-v3.md#affected-resources) for migration instructions.

```hcl
data "e2e_scaler_group" "example" {
  id         = "your-scaler-group-id"    # Replace with your scaler group ID
  project_id = "your-project-id"         # Replace with your actual project ID
  region     = "your-region"             # Use 'region' instead of deprecated 'location', e.g., "Delhi"
}
```

## Schema

### Required

- `id` (String) The unique ID of the scaler group.
- `project_id` (String) The project ID under which the scaler group is created.
- `region` (String) The region where the scaler group is deployed.
- `location` (Optional, **Deprecated**) Use `region` instead. Will be removed in v3.0.0.

### Read-Only

- `name` (String) Name of the scaler group.
- `desired` (Number) Desired number of nodes.
- `min_nodes` (Number) Minimum allowed nodes.
- `max_nodes` (Number) Maximum allowed nodes.
- `plan` (String) Name of the plan used by the scaler group.
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
