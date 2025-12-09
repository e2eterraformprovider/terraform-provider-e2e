# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### 🚧 Upcoming: v3.0.0

**🎉 Major Release - Schema Consistency & Standardization**

> **⚠️ Breaking Changes Ahead**: v3.0.0 standardizes schema field naming across all resources.
>
> **Not ready to migrate?** Stay on **v2.2.7** (latest stable). Terraform version pinning allows you to control when you upgrade.

**What's changing**:

- All resources: `location` → `region` parameter
- Load Balancer, Autoscaling, MariaDB: `plan_name` → `plan` parameter
- Block Storage, Kubernetes: `project_id` type change (integer → string)
- Image: `creation_time` → `created_at` attribute
- Object Store: `created_on` → `created_at` attribute
- Container Registry: `setup_status` → `status` attribute
- Autoscaling: `plan_id`/`sku_id` become computed-only (no longer inputs)

### Changed

- **Internal**: Completed comprehensive attribute standardization review and documentation
  - Reviewed and documented all resource-specific vs. common computed attributes
  - Confirmed consistent naming conventions (snake_case) across all resources
  - Documented common attribute patterns for future resource development
  - Identified minor inconsistencies for future consideration (VPC `state` vs `status`, DBaaS `disk` vs `disk_size`)

**Migration tools**:

- 📚 [Complete v3.0 Upgrade Guide](./docs/guides/v3-upgrade-guide.md)
- 🤖 Automated migration script (15-30 min migration time)
- 📋 Step-by-step manual migration instructions

**Release date**: TBD

---

## [2.2.7] - 2024-10-24

**Current Stable Release** - Last v2.x version

### Fixed

- **SSH Key**: Fixed read operation for SSH key resource ([#86](https://github.com/e2eterraformprovider/terraform-provider-e2e/pull/86))

### Notes

- 🔒 This is the last v2.x release
- ✅ Recommended for production workloads not ready for v3.0.0
- 📌 Pin to this version if you prefer the existing schema: `version = "2.2.7"`

---

## [2.2.6] - 2024-09-12

### Fixed

- **Documentation**: Updated resource documentation and examples
- **Image Resource**: Fixed image destroy operation
- **MySQL DBaaS**: Fixed MySQL database resource issues
- **SFS**: Fixed location parameter handling in Shared File Storage

### Changed

- Code cleanup and improvements across multiple resources

---

## [2.2.5] - 2024-08-19

### Fixed

- **VPC**: Fixed VPC creation bug that was preventing successful resource creation

---

## [2.2.4] - 2024-08-XX

### Changed

- Minor improvements and bug fixes

---

## [2.2.3] - 2024-07-XX

### Changed

- Stability improvements

---

## [2.2.2] - 2024-06-XX

### Changed

- Bug fixes and improvements

---

## [2.2.1] - 2024-06-XX

### Changed

- Minor updates

---

## [2.2.0] - 2024-05-XX

### Added

- Initial v2.2 release with enhanced features
- Foundation for future schema standardization

---

## Version History

**Choosing the Right Version**:

| Version    | Status   | Use Case                              | Support                      |
| ---------- | -------- | ------------------------------------- | ---------------------------- |
| **v3.0.0** | Upcoming | New projects, standardized schema     | ✅ Active development        |
| **v2.2.7** | Stable   | Production workloads, existing schema | ✅ Critical fixes + security |
| **v2.2.x** | Legacy   | Older deployments                     | ⚠️ Upgrade recommended       |

### Version Links

[Unreleased]: https://github.com/e2enetworks/terraform-provider-e2e/compare/v2.2.7...HEAD
[2.2.7]: https://github.com/e2enetworks/terraform-provider-e2e/compare/v2.2.6...v2.2.7
[2.2.6]: https://github.com/e2enetworks/terraform-provider-e2e/compare/v2.2.5...v2.2.6
[2.2.5]: https://github.com/e2enetworks/terraform-provider-e2e/compare/v2.2.4...v2.2.5
[2.2.4]: https://github.com/e2enetworks/terraform-provider-e2e/compare/v2.2.3...v2.2.4
[2.2.3]: https://github.com/e2enetworks/terraform-provider-e2e/compare/v2.2.2...v2.2.3
[2.2.2]: https://github.com/e2enetworks/terraform-provider-e2e/compare/v2.2.1...v2.2.2
[2.2.1]: https://github.com/e2enetworks/terraform-provider-e2e/compare/v2.2.0...v2.2.1
[2.2.0]: https://github.com/e2enetworks/terraform-provider-e2e/releases/tag/v2.2.0

---

## Generating Changelog Entries

To add a new version entry to this changelog:

```bash
make changelog VERSION=2.2.8
```

This will:

1. Extract commits since the last tag
2. Add a new changelog entry with the current date
3. Format the entry according to Keep a Changelog standards

---

## Support & Contributing

- **🐛 Issues**: [GitHub Issues](https://github.com/e2enetworks/terraform-provider-e2e/issues)
- **🚀 v3 Migration Help**: Tag issues with `v3-migration` label
- **📖 Documentation**: [Provider Docs](./docs/)
