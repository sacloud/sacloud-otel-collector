# sacloud-otel-collector Development Notes

## Development Guidelines

1. **Component Documentation**
   - Don't add individual component usage examples in README
   - Maintain the component list table instead
   - Link to official OpenTelemetry documentation for detailed configuration

2. **Git Workflow**
   - Always create a feature branch before making changes
   - Use `git add` with individual file paths (never use `git add -A`)
   - Example:
     ```bash
     git add README.md
     git add builder-config.yaml
     git add cmd/sacloud-otel-collector/components.go
     ```

3. **Testing**
   - No need to write tests for third-party components
   - Successful build is sufficient verification
   - Use `./sacloud-otel-collector validate --config config.yaml` to validate configurations

## Adding New Components

To add new components to the collector:

1. Edit `builder-config.yaml` and add the component's gomod reference
2. Run `make build-src` to regenerate the source
3. Run `make sacloud-otel-collector` to build the binary
4. Update the README.md component table if needed
5. Commit changes with descriptive message

## Component Management

- All available components are listed in `builder-config.yaml`
- The README.md maintains a table of all components with links to official documentation
- Configuration details for each component should reference the official OpenTelemetry documentation

## Upgrading OpenTelemetry Collector Version

When upgrading to a new OpenTelemetry Collector version:

1. **Update ocb binary first** - The ocb version in `Makefile` must match the component versions in `builder-config.yaml`. Using an old ocb with new components may generate incompatible code.

2. **Version mapping** - Provider versions use a different scheme:
   - Collector v0.142.0 → Providers v1.48.0
   - The pattern is: provider version = collector version + ~1.0 offset

3. **Update all versions consistently** in `builder-config.yaml`:
   - All exporters, receivers, processors, extensions: `v0.x.0`
   - All providers: `v1.x.0` (corresponding version)

4. **Fix replace directives** - After `make build-src`, check `cmd/sacloud-otel-collector/go.mod` for absolute paths in replace directives. They should be relative paths like `../../exporter/sacloudexporter`.

5. **Document breaking changes** - Create/update `docs/UPGRADE_vX_to_vY.md` for significant version upgrades with breaking changes.

6. **Verify documentation** - After writing upgrade documentation, fetch each URL and verify that the linked PR/issue content matches the description. PR numbers from changelogs are often incorrect.

7. **Verify build** - Run `make sacloud-otel-collector && make test` to ensure the upgrade works.