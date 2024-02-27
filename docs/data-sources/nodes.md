---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "e2e_nodes Data Source - terraform-provider-e2e"
subcategory: ""
description: |-
  
---

# e2e_nodes (Data Source)
DataSource will list all the created Resources of nodes.

# Example uses
```hcl 
data "e2e_nodes" "nodes324" {
   region = "Delhi"
   project_id = "325"
  }
output "all_nodes_list" {
value= data.e2e_nodes.sfs111.nodes324
}
```





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `region` (String) Region should specified
- `project_id` (String) id of project associated with nodes

### Read-Only

- `id` (String) The ID of this resource.
- `nodes_list` (List of Object) List of all the Nodes of your account . (see [below for nested schema](#nestedatt--nodes_list))

<a id="nestedatt--nodes_list"></a>
### Nested Schema for `nodes_list`

Read-Only:

- `id` (Number)
- `is_locked` (Boolean)
- `name` (String)
- `private_ip_address` (String)
- `public_ip_address` (String)
- `rescue_mode_status` (String)
- `status` (String)

