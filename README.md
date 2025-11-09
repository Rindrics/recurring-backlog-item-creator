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
- **`projects: write`** - Required to add issues to projects and update project fields
- **`contents: read`** - Required to read the configuration file from the repository

When using `GITHUB_TOKEN` in GitHub Actions, these permissions must be explicitly granted in your workflow file:

```yaml
permissions:
  contents: read
  issues: write
  projects: write
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
        uses: Rindrics/recurring-backlog-item-creator@main
        with:
          config: '.recurrent-backlog-items.yml'
```

## Usage

The workflow runs automatically on the 1st day of each month. You can also trigger it manually from the GitHub Actions tab.

### Inputs

- `config` (required): Path to the YAML configuration file

## License

See [LICENSE](./LICENSE) file for details.
