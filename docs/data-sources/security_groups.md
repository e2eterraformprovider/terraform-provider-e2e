---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "e2e_security_groups Data Source - terraform-provider-e2e"
subcategory: ""
description: |-
  
---

# e2e_security_groups (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `security_group_list` (List of Object) (see [below for nested schema](#nestedatt--security_group_list))

<a id="nestedatt--security_group_list"></a>
### Nested Schema for `security_group_list`

Read-Only:

- `description` (String)
- `id` (Number)
- `is_default` (Boolean)
- `name` (String)
- `rules` (List of Object) (see [below for nested schema](#nestedobjatt--security_group_list--rules))

<a id="nestedobjatt--security_group_list--rules"></a>
### Nested Schema for `security_group_list.rules`

Read-Only:

- `created_at` (String)
- `deleted` (Boolean)
- `id` (Number)
- `is_active` (Boolean)
- `network` (String)
- `network_cidr` (String)
- `network_size` (Number)
- `port_range` (String)
- `protocol_name` (String)
- `rule_type` (String)
- `security_group` (Number)
- `updated_at` (String)


