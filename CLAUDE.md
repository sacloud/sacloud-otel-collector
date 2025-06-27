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