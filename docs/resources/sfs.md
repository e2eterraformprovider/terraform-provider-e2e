---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "e2e_sfs Resource - terraform-provider-e2e"
subcategory: ""
description: |-
  
---

# e2e_sfs (Resource)
Provides an e2e node resource. provides an on-demand, scalable, and high-performance shared file system for Elastic Cloud Servers.

# Example uses
```hcl
 resource "e2e_sfs" "sfs1" {
    name   = "sfs-999"
    plan   = "5GB"
    vpc_id = "143"
    disk_size = 5
    project_id = "325"
    disk_iops = 75
    region = "Delhi"
 }
 ```





<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `disk_iops` (Number) input output per second
- `disk_size` (Number) size of disk to be created
- `name` (String) The name of the resource, also acts as it's unique ID
- `plan` (String) Details  of the Plan
- `project_id` (String) size of disk to be created
- `vpc_id` (String) virtual private cloud id of sfs

### Optional

- `region` (String) Location where node is to be launched
- `status` (String) status will be updated after creation

### Read-Only

- `id` (String) The ID of this resource.

