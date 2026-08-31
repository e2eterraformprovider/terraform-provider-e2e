# e2e\_security\_group (Data Source)

The `e2e_security_group` data source allows you to fetch and reference an existing security group by its name, project, and location from your E2E cloud environment.

This is useful when you want to use or reference an existing group in another resource (e.g., attaching to a VM) without creating a new one.

---

## Example Usage

```hcl
data "e2e_security_groups" "web_sg" {
  name       = "web-sg"
  location   = "Delhi"
  project_id = "42914"
}


```

---

## Schema

### Required Attributes

* **`name`** (String): Name of the existing security group.
* **`project_id`** (String): The project ID to which the security group belongs.
* **`location`** (String): The region in which the group resides.

### Read-Only Attributes

* **`id`** (String): Unique ID of the security group.
* **`description`** (String): Description of the security group.
* **`default`** (Boolean): Whether this group is the default security group.
* **`rules`** (List of Rule Blocks): A list of rules defined in the group.

---

## Nested Block: `rules`

Each rule block contains the following fields:

### Read-Only Fields

* **`rule_id`** (Number): ID of the rule.
* **`rule_type`** (String): Direction of traffic (e.g., `Inbound`, `Outbound`).
* **`protocol_name`** (String): Protocol allowed by the rule (e.g., `All`, `Custom_TCP`).
* **`port_range`** (String): Port range this rule applies to.
* **`network`** (String): Type of network source (e.g., `myNetwork`, `manual`, `any`).
* **`network_cidr`** (String): The CIDR for manual or VPC networks. 
* **`size`** (Number): Network size used with `myNetwork` or CIDR.
* **`description`** (String): Description of the rule.

---

## Notes

* Ensure the `name`, `project_id`, and `location` combination is valid and corresponds to an existing security group.
* If no matching group is found, Terraform will return an error during the plan or apply stage.
