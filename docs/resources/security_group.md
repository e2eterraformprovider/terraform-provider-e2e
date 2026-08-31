# e2e\_security\_group (Resource)

The `e2e_security_group` resource allows you to provision and manage security groups within your E2E cloud environment. Security groups define inbound and outbound rules to control network traffic at the instance level.

Creating this resource provisions a new security group with specified firewall rules. Destroying it deletes the group.

---

## Example Usage

```hcl
resource "e2e_security_groups" "web_sg" {
  name        = "web-sg"
  location    = "Delhi"
  project_id  = "42914" 
  description = "Web Tier Security Group"
  default     = false

  rules{
      rule_type     = "Inbound"
      protocol_name = "Custom_TCP"
      port_range    = "80"
      network       = "myNetwork"
      network_cidr  = "vpc_4012" #provide your network_cidr. Find the format below in optional feilds.
      size          = 24
      description   = "Allow HTTP traffic"
    }

    rules{
      rule_type     = "Outbound"
      protocol_name = "All"
      port_range    = "All"
      network       = "any"
      description   = "Allow all outbound"
    }
  
}
```

---


## Schema

### Required Attributes

* **`name`** (String): Name of the security group.
* **`project_id`** (String): The project ID associated with the security group.
* **`rules`** (List of Rule Blocks): A list of rule blocks defining firewall rules.
* **`location`** (String): The region in which to create the security group. 

### Optional Attributes

* **`description`** (String): Description of the security group.
* **`default`** (Boolean): Whether this group is the default group. Defaults to `false`.

### Read-Only Attributes

* **`id`** (String): Unique ID of the security group.

---

## Nested Block: `rules`

Each element in the `rules` list supports the following attributes:

### Required Fields

* **`rule_type`** (String): Direction of traffic. Allowed values: `"Inbound"`, `"Outbound"`.

### Optional Fields

* **`rule_id`** (Number): ID of the rule (computed).
* **`protocol_name`** (String): Protocol to allow. Allowed values: `"All"`, `"All_TCP"`, `"All_UDP"`, `"ICMP"`, `"Custom_TCP"`, `"Custom_UDP"`. Defaults to `"All"`.
* **`port_range`** (String): Port range to allow. Defaults to `"All"`.
* **`network`** (String): Network type. Allowed values: `"myNetwork"`, `"manual"`, `"any"`. Defaults to `"any"`.
* **`network_cidr`** (String): The CIDR block for the rule. If VPC network then format must be 'vpc_<`vpc_id`>'. For Manual network one may go for the IP address.
* **`size`** (Number): Size of the network if `myNetwork` is used or manual CIDR provided.
* **`description`** (String): Description of the rule.

---

## Notes

* Rules with `network = "myNetwork"`, `size` is automatically set to `512` if not specified. `network_cidr` must be explicitly provided.
* Rules with `network = "myNetwork"`  must provide `network_cidr` explicitly.
* Rules with `network = "manual"` must provide `size` and `network_cidr` explicitly.
