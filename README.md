# gh-issue-config-filter

A CLI tool to filter GitHub issues based on creation months from a YAML configuration file.

## Motivation

When managing product backlogs, you often need to create recurring issues on different schedules. Each issue type has its own creation frequency:

- **Monthly issues**: Security reviews, performance monitoring (created every month)
- **Quarterly issues**: Planning sessions, roadmap reviews (created in March, June, September, December)
- **Yearly issues**: Annual reviews, budget planning (created once a year)
- **Custom schedules**: Any combination of specific months

Manually tracking which issues to create in which month is tedious and error-prone. This tool automates the process by:

1. **Defining issue templates in YAML**: Configure which issues to create, their specific creation months, and what fields to set
2. **Filtering by month**: Automatically determine which issues should be created in a given month based on their individual schedules
3. **Integration with GitHub Actions**: Use the filtered results in your CI/CD pipeline to automatically create issues

This tool is designed to work with GitHub Projects, allowing you to automatically set project fields (like Story Points, Status, etc.) when issues are created.

## Installation

```bash
go install github.com/Rindrics/gh-issue-config-filter@latest
```

Or build from source:

```bash
make build
```

## Usage

```bash
gh-issue-config-filter --month <1-12> [--config <config-file>]
```

### Options

- `--month`: Month (1-12) to filter issues (required)
- `--config`: Path to config file (default: `.gh-issue-config-filter.yml`)

## Example

```bash
# Get issues to create in January
gh-issue-config-filter --month 1 --config config-template.yml

# Output:
[
  {
    "name": "Wash My Cat",
    "template_file": ".github/ISSUE_TEMPLATE/wash_my_cat.md",
    "title_suffix": "- $(date +%Y-%m)",
    "fields": {
      "priority": "1",
      "status": "Ready"
    },
    "project_id": "PVT_kwHOAOKHl84BHgin",
    "target_repo": "Rindrics/gh-issue-config-filter"
  }
]
```

```bash
# Two issues are returned for March
gh-issue-config-filter --month 3 --config config-template.yml

# Output:
[
  {
    "fields": {
      "priority": "1",
      "status": "Ready"
    },
    "name": "Wash My Cat",
    "project_id": "PVT_kwHOAOKHl84BHgin",
    "target_repo": "Rindrics/gh-issue-config-filter",
    "template_file": ".github/ISSUE_TEMPLATE/wash_my_cat.md",
    "title_suffix": "- $(date +%Y-%m)"
  },
  {
    "fields": {
      "priority": "2",
      "status": "Backlog"
    },
    "name": "Buy New Shoes",
    "project_id": "other_project_id",
    "target_repo": "other/repo",
    "template_file": ".github/ISSUE_TEMPLATE/buy_shoes.md",
    "title_suffix": "- $(date +%Y)"
  }
]
```


## Configuration File Format

See [`config-template.yml`](./config-template.yml) for an example configuration file.
