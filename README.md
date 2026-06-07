# Eco System

Eco System is a monorepo for a simple online store platform. It includes a
customer storefront plus administration apps for users, products, orders, and
shared UI components.

The project is organized as several business domains that share a common
deployment model. Backend services are written in Go, frontend applications are
Angular apps managed with pnpm, and Kubernetes releases are packaged with Helm.

See [CHANGELOG.md](CHANGELOG.md) for the project change log. The previous
detailed deployment README has been preserved at
[docs/repo/README.md](docs/repo/README.md).

## What Is Deployed

The deployed system is a small store platform:

- `online-storefront` provides the customer-facing web app, product catalog,
  authentication-aware cart, checkout, and order history.
- `user-management` provides user-service, user-gateway, and the admin web app
  for authentication and user administration.
- `product-management` provides product-service, product-gateway, and the
  product admin web app for catalog and promotion management.
- `order-management` provides order-service, order-gateway, and the order admin
  web app for order lifecycle and fulfillment workflows.
- `shared-components` provides the Angular design system packages and Storybook.
- `foundation` provides the shared database, migrations, and object storage.
- `ci-cd` provides Nexus Repository Manager for container images built by the
  GitHub Actions pipeline.

## Repository Structure

```text
.
|-- ci-cd/                  # Nexus and CI/CD infrastructure Helm charts
|-- config/                 # Shell-compatible environment files for deploys
|-- deployment/charts/app/  # Shared reusable Helm app chart
|-- docs/repo/              # Archived detailed operational README
|-- foundation/             # Postgres, Liquibase, MinIO, database migrations
|-- online-storefront/      # Storefront frontend and storefront/cart backends
|-- order-management/       # Order admin frontend and order backends
|-- product-management/     # Product admin frontend and product backends
|-- shared-components/      # Shared Go packages, Angular design system, Storybook
`-- user-management/        # User admin frontend and user backends
```

Most domains follow the same shape:

- `backend/app/*` contains deployable Go services or gateways.
- `frontend/app/*` contains deployable Angular applications.
- `deployment/` contains umbrella Helm charts for the domain.
- Each deployable app has its own nested `deployment/` chart that depends on the
  shared `deployment/charts/app` chart.

## Technology

- **Backend:** Go, Fiber/Iris-style route layers, shared backend packages, and
  Postgres access.
- **Frontend:** Angular, TypeScript, Sass, pnpm workspaces, and Storybook.
- **Database:** Postgres with Liquibase migrations and seed scripts.
- **Storage:** MinIO for S3-compatible product media storage.
- **Containers:** Docker images for every backend, frontend, Storybook, and
  Liquibase migration job.
- **Orchestration:** Kubernetes, Helm, and namespace-scoped releases.
- **CI/CD:** GitHub Actions builds, tests, pushes images to Nexus, and deploys
  Helm releases.
- **Dependencies:** Renovate is configured for frontend, backend, Docker,
  GitHub Actions, and Helm dependency updates.

## Namespaces

The deployment pipeline is split into three namespace layers.

### `eco-cicd`

The CI/CD namespace is deployed first. It runs Nexus Repository Manager and a
hosted Docker registry named `eco-images`.

The GitHub runner builds backend and frontend images, pushes them into Nexus,
and later Kubernetes pulls those images from the same registry. The registry
address is configured in `config/hp-prodesk-homelab.txt`:

```sh
CI_CD_NAMESPACE=eco-cicd
NEXUS_DOCKER_REGISTRY=172.18.0.4:30500
NEXUS_DOCKER_NODE_PORT=30500
IMAGE_NAMESPACE=eco-system
```

`NEXUS_DOCKER_REGISTRY` must be reachable by the GitHub runner and every
Kubernetes node. For multi-node clusters, do not use `127.0.0.1` or
`localhost`.

Nexus is exposed as an HTTP Docker registry. The deploy workflow configures the
GitHub runner Docker daemon as an insecure registry before pushing images. Each
Kubernetes node also needs its container runtime configured to trust the same
HTTP registry address, otherwise deployments can push successfully but pods will
fail when pulling images.

For k3d clusters, the deploy workflow configures each k3d node container with
the matching containerd registry mirror and restarts the node containers so
pulls use HTTP. Other Kubernetes runtimes should be preconfigured manually.

The Nexus admin password can be supplied as a GitHub secret named
`NEXUS_ADMIN_PASSWORD`. If the secret is not set, the workflow uses Nexus'
initial local default of `admin123`. Any `sudo` password prompt comes from the
self-hosted runner operating system, not from Nexus. The runner either needs
passwordless sudo for Docker daemon configuration or `/etc/docker/daemon.json`
must already include the Nexus registry under `insecure-registries`.

For k3s nodes, configure the registry in `/etc/rancher/k3s/registries.yaml`
before deploying workloads:

```yaml
mirrors:
  "172.18.0.4:30500":
    endpoint:
      - "http://172.18.0.4:30500"
```

Then restart k3s on each node so containerd picks up the registry settings.

### `eco-foundation`

The foundation namespace is deployed after CI/CD. It runs shared infrastructure:

- Postgres for application data.
- Liquibase as a Helm hook job for schema migrations and seed data.
- MinIO for product media.

The Liquibase image is built by the workflow, pushed to Nexus, then pulled by
the foundation namespace using a `nexus-docker-registry` pull secret.

### `eco-test`

The application namespace is deployed after CI/CD and foundation. It runs the
storefront, admin apps, gateways, services, and Storybook.

By default, `eco-test` is the main application namespace:

```sh
DEFAULT_APPLICATION_NAMESPACE=eco-test
```

Feature pull requests use a namespace derived from the branch name. Branches
must match `feature/<10 lowercase alphanumeric chars>`, and that suffix becomes
the namespace.

## Deployment Flow

The GitHub Actions `Build and Deploy` workflow can deploy a single layer or all
layers. On a fresh cluster, choose `all`.

The expected order is:

1. `prepare` resolves config, namespaces, and hostnames.
2. `ci_cd` creates `eco-cicd`, deploys Nexus, and creates the Docker registry.
3. `foundation` creates `eco-foundation`, pushes the Liquibase image to Nexus,
   and deploys Postgres, MinIO, and Liquibase.
4. `application` creates `eco-test`, builds and pushes app images to Nexus,
   creates image pull secrets, and deploys the store apps and services.

The workflow also runs Go functional tests, builds Angular apps and Storybooks,
and waits for Kubernetes rollouts before running deployed smoke tests.

## Useful Files

- [config/hp-prodesk-homelab.txt](config/hp-prodesk-homelab.txt) contains
  namespace, host, database, and registry settings.
- [.github/workflows/deploy.yml](.github/workflows/deploy.yml) builds and
  deploys the platform.
- [.github/workflows/delete-deployment.yml](.github/workflows/delete-deployment.yml)
  removes application, foundation, or CI/CD releases.
- [.github/workflows/renovate.yml](.github/workflows/renovate.yml) runs
  Renovate manually from GitHub Actions.
- [renovate.json](renovate.json) configures dependency update behavior.

## Local Notes

Install dependencies from the repository root:

```sh
pnpm install
```

Build all frontend workspace packages:

```sh
pnpm -r build
```

Run Go tests:

```sh
go test ./...
```

For detailed historical Helm rendering and manual deployment notes, see
[docs/repo/README.md](docs/repo/README.md).
