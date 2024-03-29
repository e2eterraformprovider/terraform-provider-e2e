---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "e2e_node Data Source - terraform-provider-e2e"
subcategory: ""
description: |-
  
---

# e2e_node (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `label` (String) The name of the group
- `name` (String) The name of the resource, also acts as it's unique ID
- `os` (String) OS and its version  format : <OS>-<version>
- `plan` (String) The name of the Plan

### Optional

- `backup` (Boolean) Tells you the state of your backups
- `default_public_ip` (Boolean) Tells us the state of default public ip
- `disable_password` (Boolean)
- `enable_bitninja` (Boolean)
- `image` (String) The name of the image you have selected
- `is_ipv6_availed` (Boolean)
- `is_saved_image` (Boolean)
- `lock_node` (Boolean)
- `ngc_container_id` (Number)
- `power_status` (String)
- `reboot_node` (Boolean)
- `region` (String)
- `reinstall_node` (Boolean)
- `reserve_ip` (String)
- `saved_image_template_id` (Number)
- `security_group_id` (Number)
- `ssh_keys` (Set of String)
- `vpc_id` (String)

### Read-Only

- `created_at` (String)
- `disk` (String)
- `id` (String) The ID of this resource.
- `is_active` (Boolean)
- `is_bitninja_license_active` (Boolean)
- `is_monitored` (Boolean)
- `memory` (String)
- `price` (String)
- `private_ip_address` (String)
- `public_ip_address` (String)
- `status` (String)


