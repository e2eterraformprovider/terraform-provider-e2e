# e2e\_dbaas\_postgres (Resource)

The `e2e_dbaas_postgresql` resource allows you to provision and manage PostgreSQL-based DBaaS (Database as a Service) clusters within your E2E cloud environment.

Applying this resource creates a new PostgreSQL DBaaS instance. Destroying it removes the instance.

---

## Example Usage

```hcl
resource "e2e_dbaas_postgresql" "db1" {
  location     = "Delhi"
  project_id   = 12345
  plan         = "DBS.16GB"
  version      = "15.0"
  name         = "mydbname"

  database {
    user         = "admin"
    password     = "SecurePasswordItis@123"
    name         = "mydb"
    dbaas_number = 1
  }

  vpc_list = [e2e_vpc.VPC-TS-01.id]
}

resource "e2e_vpc" "VPC-TS-01" {
  location    = "Delhi"
  vpc_name    = "VPC-TS-01"
  project_id  = "your_project_id"
}
```

---

## Schema

### Required Attributes

- **`location`** (String): Region in which the DBaaS instance will be deployed.
- **`project_id`** (String): The project ID associated with the DBaaS.
- **`plan`** (String): The DBaaS plan (e.g., `"DBS.16GB"`).
- **`version`** (String): Desired PostgreSQL version (e.g., `"15.0"`).
- **`name`** (String): Name of the DBaaS instance.
- **`database`** (Block): Configuration block for the database (see [Database Block](#nested-block-database)).

### Optional Attributes

- **`group`** (String): Group name for organizing resources. Defaults to `"Default"`.
- **`parameter_group_id`** (Number): ID of the parameter group.
- **`detach_public_ip`** (Boolean): Set to `true` or `false` to detach or reattach the public IP. Defaults to `false`.
- **`power_status`** (String): Power control operation. Accepts: `"start"`, `"stop"`, or `"restart"`.
- **`size`** (Number): Disk size (in GB) for upgrades. Note: Instance must be stopped for upgrade.
- **`vpc_list`** (Set of Number): List of VPC IDs to attach. Remove an ID to detach a VPC.
- **`is_encryption_enabled`** (Boolean): Enable encryption. Defaults to `false`. Must be added only during creation of an instance.

### Read-Only Attributes

- **`id`** (String): Unique ID of the DBaaS instance.
- **`status`** (String): Current status of the DBaaS.
- **`status_title`** (String): Human-readable status.
- **`status_actions`** (List of String): Permissible operations (e.g., `delete`, `stop`, `restart`).
- **`num_instances`** (Number): Number of database instances in the cluster.
- **`project_name`** (String): Name of the project.
- **`snapshot_exist`** (Boolean): Indicates if snapshots exist.
- **`connectivity_detail`** (String): Read/write connectivity information.
- **`vector_database_status`** (String): Status of the vector database feature.

---

## Nested Block: `database`

### Required Fields

- **`user`** (String): Username for the database.
- **`password`** (String): Password for the DB user.
- **`name`** (String): Name of the database.
