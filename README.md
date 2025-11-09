# recurring-backlog-item-creator

Automatically create and manage GitHub issues in your project backlog based on a YAML configuration file

## Overview

This repository contains:

- **`gh-issue-config-filter`**: A CLI tool that filters issues from a YAML configuration file based on creation months
- **GitHub Actions workflow**: Automatically creates issues and registers them to GitHub Projects based on the filtered results

## Motivation

When managing product backlogs, you often need to create recurring issues on different schedules.
Each issue type has its own creation frequency:

- **Monthly issues**: Security reviews, performance monitoring (created every month)
- **Quarterly issues**: Planning sessions, roadmap reviews (created in March, June, September, December)
- **Yearly issues**: Annual reviews, budget planning (created once a year)
- **Custom schedules**: Any combination of specific months

Manually tracking which issues to create in which month is tedious and error-prone.
This tool automates the entire process of creating issues and managing them in GitHub Projects.

## Components

### gh-issue-config-filter

A CLI tool that filters issues from a YAML configuration file based on the current month.
See [`gh-issue-config-filter/README.md`](./gh-issue-config-filter/README.md) for detailed documentation.

## Setup

### 1. Create a configuration file

Create a YAML configuration file in your repository.
See [`config-template.yml`](./config-template.yml) for an example.

### 2. Configure GitHub Variables (Optional)

If you want to use a default project ID for all issues, set the following repository variable:

- `PROJECT_ID`: Your GitHub Project ID (can be overridden per issue in config)

> [NOTE]
> You don't need to set IDs of project item fields (like Story Points, Status) manually
> because they are automatically detected from field names in your configuration file

### 3. Configure GitHub Token Permissions

The GitHub token (`GITHUB_TOKEN` or a custom token) requires the following permissions:

- **`issues: write`** - Required to create issues in the repository
- **`repository-projects: write`** - Required to add issues to projects and update project fields (includes read access)
- **`contents: read`** - Required to read the configuration file from the repository

When using `GITHUB_TOKEN` in GitHub Actions, these permissions must be explicitly granted in your workflow file:

```yaml
permissions:
  contents: read
  issues: write
  repository-projects: write
```

> [NOTE]
> If you're using a Personal Access Token (PAT) or a custom token, ensure it has the `repo` scope (for private repositories) or `public_repo` scope (for public repositories), which includes the necessary permissions.

### 4. Create GitHub Actions workflow

Create `.github/workflows/create-monthly-issues.yml`:

```yaml
name: Create Monthly Issues

on:
  schedule:
    # Run on the 1st day of every month at 00:00 UTC
    - cron: '0 0 1 * *'
  workflow_dispatch: # Allow manual trigger

permissions:
  contents: read
  issues: write
  projects: write

jobs:
  create-issues:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Create recurring backlog items
        uses: Rindrics/recurring-backlog-item-creator@latest
        with:
          config: '.recurrent-backlog-items.yml'
```

## Usage

The workflow runs automatically on the 1st day of each month. You can also trigger it manually from the GitHub Actions tab.

### Inputs

- `config` (required): Path to the YAML configuration file

### How It Works

1. The Action determines the current month
2. The CLI tool filters issues from your configuration file based on the current month
3. For each matching issue:
   - Creates a GitHub issue with the specified title and template
   - Adds the issue to the specified GitHub Project
   - Sets project fields (like Priority, Status, Story Points, etc.) automatically

### Template Variables

The `title_suffix` field in your configuration supports template variables:

- `{{Year}}`: Current year (e.g., `2025`)
- `{{Month}}`: Current month (e.g., `01`)
- `{{YearMonth}}`: Current year and month (e.g., `2025-01`)
- `{{Date}}`: Current date in YYYY-MM-DD format (e.g., `2025-01-15`)

Example:

```yaml
title_suffix: "- {{YearMonth}}"  # Results in "- 2025-01"
```

### Output Format

The CLI tool outputs JSON in a format compatible with GitHub Projects GraphQL API. Field IDs and option IDs are automatically resolved from field names and option names in your configuration file.

## Examples

### Basic Configuration

```yaml
defaults:
  project_id: "PVT_xxx"
  target_repo: "owner/repo"

issues:
  - name: "Monthly Security Review"
    template_file: ".github/ISSUE_TEMPLATE/security-review.md"
    creation_months: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
    title_suffix: "- {{YearMonth}}"
    fields:
      Priority: "High"
      Status: "Ready"
```

### Override Default Project

```yaml
issues:
  - name: "Quarterly Planning"
    template_file: ".github/ISSUE_TEMPLATE/planning.md"
    creation_months: [3, 6, 9, 12]
    project_id: "PVT_yyy"  # Override default project
    target_repo: "owner/other-repo"  # Override default repo
    fields:
      Priority: "Medium"
      Status: "Backlog"
```

### GITHUB_TOKEN Permission Errors

Ensure your workflow has the required permissions:

- `contents: read` - To read the configuration file
- `issues: write` - To create issues
- `repository-projects: write` - To add issues to projects and update fields

## License

See [LICENSE](./LICENSE) file for details.
