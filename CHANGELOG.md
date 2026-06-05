# Project Change Log

## v1.18.0 - (5 Changes)
- Added shared backend order request and response contracts for reuse across storefront and order management.
- Added the order-management backend use case with order-gateway and order-service applications.
- Added order-management APIs for order listing, order detail lookup, and fulfillment status updates.
- Added the order-admin-web-app for order management and fulfillment workflows.
- Added Helm deployment charts and ingress wiring for order-management backend services.

## v1.17.0 - (5 Changes)
- Added order and order item persistence for storefront checkout with product snapshots and order totals.
- Added a cart checkout API route that creates an order, reserves product inventory, and closes the active cart transactionally.
- Added storefront checkout UI on the cart page with order confirmation feedback.
- Added signed-in storefront order history for reviewing previous checkout orders.
- Added checkout coverage to cart gateway route and deployed functional tests.

## v1.16.0 - (4 Changes)
- Added a storefront cart service with backend synchronization, authentication-aware state, and cart item counts.
- Added cart navigation, add-to-cart actions, and a responsive cart page for quantity updates, removal, and clearing.
- Updated storefront account messaging now that cart functionality is available.
- Prefixed Kubernetes resource names for the user admin web app and storefront cart gateway so deployed workloads are easier to identify.

## v1.15.0 - (5 Changes)
- Added the online storefront cart-gateway backend application with authenticated customer cart tracking.
- Added cart and cart item persistence with one active cart per customer, product links, quantity checks, and Liquibase grants.
- Added cart API routes for reading, adding, updating, removing, and clearing the current customer's cart.
- Added Helm, local deployment, and GitHub Actions wiring for building, routing, deploying, and rolling out cart-gateway.
- Added deployed functional tests for storefront-gateway catalog and authentication flows plus cart-gateway cart lifecycle flows.

## v1.14.0 - (3 Changes)
- Added the online storefront Go gateway with public active-product catalog reads, signup, login, logout, current-user lookup, product media proxying, focused functional tests, and Helm deployment charts.
- Added the online storefront Angular web app with a public catalog, product detail pages, account registration, sign-in, sign-out, session-aware navigation, and Helm deployment wiring.
- Added GitHub Actions and local deployment automation for building, importing, deploying, rolling out, and smoke testing the online storefront gateway and web app.

## v1.13.1 - (9 Changes)
- Merge public and private service into single user service for user-management
- Added promotion settings support for enabling and disabling promotions from the backend and product admin web app.
- Added quantity bonus discount support with buy quantity, free quantity, and minimum product count fields.
- Added product admin promotion management UI updates for discount creation, editing, activation, archiving, and settings.
- Split promotion persistence into separate discount and promotion settings repositories.
- Refactored product-management and user-management repository layers to use domain models instead of request DTOs for create and update writes.
- Changed repository create and update methods to return only possible errors while services keep ownership of request-to-entity mapping.
- Moved product media storage into the product repository layer so image persistence sits with the other storage adapters.
- Added Style Dictionary design tokens for the Angular design system with generated SCSS variables, Sass maps, and CSS custom property exports.

## v1.13.0 - (6 Changes)
- Added product labels across the product API, Postgres persistence, Liquibase migration, mocks, and product admin UI.
- Added a Promotions page for creating, scheduling, viewing, activating, and archiving product promotions.
- Moved discount configuration out of product list/edit screens so products stay focused on catalog data.
- Added promotion targeting for all products, categories, product labels, and selected products.
- Updated product list and detail screens to show labels plus discounted final price beside the original price.
- Updated discount price calculations to respect active status plus promotion start and end dates.

## v1.12.1 - (5 Changes)
- Enhanced deployment pipeline to automatically detect and deploy foundation layer changes on push to main and feature branches.
- Feature branches with foundation changes now deploy both foundation and application layers to the same namespace for better isolation.
- Fixed workflow ordering so foundation deploy finishes before application deploy when both are triggered.
- Aligned Liquibase and Postgres passwords during foundation deploy to prevent `liquibase` authentication failures.
- Updated pull request decorator with styled HTML formatting displaying useful links and namespace information at the top of the PR overview section.
- Added `scripts/local-deploy.sh` to automate local k3d/k3s image builds, imports, namespace creation, and Helm deploys for foundation/application.
- added discount functionality to product-management 

## v1.12.0 - (5 Changes)
- Added product short descriptions and product photo support.
- Added MinIO-based S3-compatible object storage in the foundation layer.
- Added product photo upload flow with generated thumbnail and detail image variants.
- Added product media serving through the product gateway with cache headers.
- Updated product list, detail, and edit screens to show uploaded product images efficiently.

## v1.11.1 - (7 Changes)
- Moved shared user request, response and auth constants from user-management into shared-components backend user packages.
- Refactored user-management services, repositories, controllers and tests to consume the shared user package contracts.
- Changed product-gateway login and logout to call user-service through an internal auth client instead of handling auth locally.
- Added user-service auth routes for internal login and logout calls from other backend services.
- Added product-gateway `USER_SERVICE_URL` deployment configuration for service-to-service auth calls.
- Improved product-gateway auth request-id propagation so login and logout internal user-service calls keep the original request id.
- Aligned product-management login mocks and functional tests with the seeded `admin@test.com` user-management account.

## v1.11.0 - (11 Changes)
- Added product-management as a new domain-driven monorepo use case for online store catalog administration.
- Added product-management product-gateway and product-service backend applications with product CRUD routes.
- Added product domain models, services, repository interfaces and Postgres persistence for catalog products.
- Added Liquibase product table creation, seed and app permission scripts for default catalog data.
- Requires a foundation or all-layer deployment once so Liquibase can apply the product table and seed data to existing environments.
- Added product-admin-web-app for authenticated admin users to view, edit and manage products.
- Added product-management Helm deployment charts and test values for backend services and product admin web app.
- Added product-management build, image import, Helm deploy, rollout and useful link wiring to the deployment pipeline.
- Changed product-admin-web-app deployment wiring to use same-origin `/api` ingress routing and avoid browser mixed-content blocking.
- Added in-process and deployed product-gateway functional tests for login, list products, get product, create product, edit product and delete product flows.
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
