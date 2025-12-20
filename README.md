# Hatchet Workflow Examples

A collection of workflow examples using Hatchet in Go, demonstrating distributed task orchestration and execution patterns.

## Overview

This repository contains practical examples of implementing workflows with [Hatchet](https://hatchet.run), a distributed task orchestration engine. The examples are written in Go and showcase various workflow patterns and task execution scenarios.

## Project Structure

```
.
├── cmd/                   # Application entry points
├── pkg/tasks/             # Task implementations and workflow definitions
├── .devcontainer/         # Development container configuration
├── .vscode/               # VS Code workspace settings
├── docker-compose.dev.yml # Docker Compose setup for development
├── Dockerfile.dev         # Development Docker image
└── .env.example           # Environment variables template
```

## Getting Started

### Prerequisites

- Go 1.x or higher
- Docker and Docker Compose (for containerized development)
- Hatchet server instance (local or cloud)

### Environment Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Configure your Hatchet connection settings in the `.env` file

### Running with Docker

The project includes a complete Docker development environment:

```bash
docker-compose -f docker-compose.dev.yml up
```

### Development Container

This repository supports VS Code Dev Containers for a consistent development environment. Open the project in VS Code and select "Reopen in Container" when prompted.

## Features

- Workflow orchestration examples
- Task execution patterns
- Distributed processing demonstrations
- Docker-based development environment
- Dev Container support for VS Code

## Technology Stack

- **Language**: Go (93.3%)
- **Orchestration**: Hatchet
- **Containerization**: Docker

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.


## Resources

- [Hatchet Documentation](https://docs.hatchet.run)
- [Hatchet GitHub](https://github.com/hatchet-dev/hatchet)

---

**Note**: This is a personal project by [@raykavin](https://github.com/raykavin) for learning and demonstrating Hatchet workflow patterns.