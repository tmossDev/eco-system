# Project Change Log

## v1.11.0 - (7 Changes)
- Added product-management as a new domain-driven monorepo use case for online store catalog administration.
- Added product-management product-gateway and product-service backend applications with product CRUD routes.
- Added product domain models, services, repository interfaces and Postgres persistence for catalog products.
- Added Liquibase product table creation and seed scripts for default catalog data.
- Added product-admin-web-app for authenticated admin users to view, edit and manage products.
- Added product-management Helm deployment charts and test values for backend services and product admin web app.
- Added reusable box icon support to the angular design system for product navigation.

## v1.10.0 - (8 Changes)
- Added test deployment values files for shared-components and user-management to reduce long Helm `--set` deployment commands.
- Added generated runtime values file instructions for image tags, hostnames, database settings and secret values.
- Added reusable app chart support for externally managed ConfigMaps through `configMap.existingName`.
- Added reusable app chart support for nginx `conf.d` ConfigMaps through `nginxConfig`.
- Added a separately managed admin web app runtime config ConfigMap for test deployments.
- Added durable admin web app ingress routing for same-origin `/api` calls to user-gateway.
- Added shared Storybook host routing for admin web app stories at `/stories/admin-web-app`.
- Moved shared Storybook static files to `/var/www/storybook` and added nginx routing for `/storybook/`.

## v1.9.0 - (4 Changes)
- Added an authenticated user menu to the admin web app navigation with current user details, profile navigation and logout.
- Added logout handling in the admin web app to clear the stored session and return users to the login page.
- Improved backend access logs to emit compact single-line JSON with request IDs, stdout output and redacted request/response bodies.
- Added in-process and deployed user-gateway functional tests for login, logout, list users, get user, create user, edit user and delete user flows to the deployment pipeline.

## v1.8.0 - (8 Changes)
- Added k3d-aware local image import and browser access troubleshooting to the deployment README.
- Changed Liquibase Kubernetes execution from a long-running deployment to a Helm hook job.
- Added Liquibase changesets for application database grants and beta user role assignments.
- Changed Go application build instructions to produce static Linux binaries for Alpine runtime images.
- Added gateway routes for admin web app auth, dashboard, settings and user management API calls.
- Updated gateway login compatibility for frontend email login requests and `accessToken` responses.
- Fixed Postgres user lookup projection and app user table permissions for seeded local users.
- Improved backend request id propagation and database error logging across login repository calls.

## v1.7.0 - (9 Changes)
- Split deployment automation into foundation and application layers.
- Added `config/hp-prodesk-homelab.txt` for shell-compatible homelab deployment values.
- Added feature pull request namespace handling for `feature/<10 lowercase alphanumeric chars>` branches.
- Changed application Dockerfiles to copy prebuilt Go, Angular and Storybook artifacts instead of building from registry-pulled images.
- Added runtime frontend config map support for `BACKEND_API_URL` and mock API toggling.
- Added ingress host wiring for application services and the shared design-system Storybook deploy.
- Added k3s image import support for locally built images that are not pushed to a registry.
- Changed homelab ingress defaults to Traefik with HTTP local hostnames.
- Updated Docker ignore rules so prebuilt backend, admin web app and Storybook artifacts are available to runtime image builds.

## v1.6.1 - (4 Changes)
- Added README instructions for manually deploying Helm releases to Kubernetes.
- Added GitHub Actions workflow for building container images and deploying Helm releases to Kubernetes.
- Added frontend Dockerfiles for admin-web-app and storybook deployment images.
- Added Kubernetes postgres init script wiring for database users used by liquibase and application services.

## v1.6.0 - (6 Changes)
- Added reusable Helm app chart for common Kubernetes deployment resources.
- Added Helm deployment charts for foundation postgres and liquibase services.
- Added Helm deployment charts for shared-components storybook.
- Added Helm deployment charts for user-management user-service, user-gateway and admin-web-app.
- Added umbrella Helm deployment charts for foundation, shared-components and user-management.
- Added README instructions for building Helm dependencies and rendering Kubernetes resource files from Helm templates.

## v1.5.0 - (6 Changes)
- Renamed/reworked Header Component to Navigation Component.
- Added Login Form Component to angular design system.
- Added Login, Forgot Password, User Management and Settings pages to user-management admin-web-app.
- Added Guard to user-management admin-web-app forcing users to login first.
- Added Interceptors for authentication, error handling and api-mocking.
- Added mocking for user-service and auth-service.

## v1.4.1 - (3 Changes)
- Added Icon component to angular design system.
- Added Header component to angular design system which uses Icon component.
- Using Header component in user-management admin-web-app.

## v1.4.0 - (5 Changes)
- Added pnpm for frontend package management.
- Added user-management web-app which is just hello world angular project
- Added angular design system for all angular web-apps in mono repo.
- Added button to angular design system.
- Added storybook for angular design system.

## v1.3.0 - (2 Changes)
- Added user-service controller for user management external calls.
- Added user-service routes for user management external calls using user-service controller.

## v1.2.0 - (4 Changes)
- Added user package for user management code reuse.
- Added user-gateway controller for user management external calls.
- Added user-gateway middleware for user management external calls.
- Added user-gateway routes for user management external calls using user-gateway controller.

## v1.1.0 - (4 Changes)
- Added shared-components for common components.
- Added GoLang common libraries to shared-components/backend/lib.
- Added user-service shell application for user management internal calls.
- Added user-gateway shell application for user management external calls.

## v1.0.0 - (4 Changes)
- Added foundation use case for getting local postgres database up and running.
- Added database Init Script for starting up the shop database from scratch
- Added liquibase setup and seed scripts for postgres database.
- Added a basic docker-compose file for foundation use case.
