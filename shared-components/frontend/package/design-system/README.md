# DesignSystem

This project was generated using [Angular CLI](https://github.com/angular/angular-cli) version 21.2.10.

## Component Boundaries

The design system is split into four layers. New components should enter the lowest layer that can honestly own them.

### Foundation

Foundation packages define the visual contract. They do not render product UI.

- `@foundation/tokens`: colors, spacing, radius, type scale, motion, and shadows.

### Primitives

Primitive packages are small building blocks. They should be generic, token-driven, accessible, and free of product workflows, routing, data fetching, and business copy.

- `@primitive/icon`
- `@primitive/button`
- `@primitive/form`

Primitive components should expose constrained variants instead of arbitrary visual escape hatches. For example, prefer `variant="primary"` over passing custom colors.

### Patterns

Pattern packages compose primitives into reusable UI structures. They may know about common layout behavior, but they should still avoid owning application state or domain workflows.

- `@shared/navigation-bar`
- `@shared/login-form`

Patterns depend on primitives, never the other way around. If a component starts needing service calls, feature-specific routes, permission checks, or domain state, it belongs outside the design system.

### Features

Feature components live in application or feature packages such as `auth-features` and `admin-features`. They can compose design-system primitives and patterns while owning product behavior.

Examples: login pages, authenticated dashboard shells, route-aware pages, API-backed forms, and permission-specific navigation.

## Development server

To start a local development server, run:

```bash
ng serve
```

Once the server is running, open your browser and navigate to `http://localhost:4200/`. The application will automatically reload whenever you modify any of the source files.

## Code scaffolding

Angular CLI includes powerful code scaffolding tools. To generate a new component, run:

```bash
ng generate component component-name
```

For a complete list of available schematics (such as `components`, `directives`, or `pipes`), run:

```bash
ng generate --help
```

## Building

To build the project run:

```bash
corepack pnpm build
```

This will compile your project and store the build artifacts in the `dist/` directory. By default, the production build optimizes your application for performance and speed.

## Design Tokens

Design tokens live in `foundation/tokens/tokens` and are generated with Style Dictionary into `@foundation/tokens`.

```bash
corepack pnpm build:tokens
```

Components can use the generated SCSS module:

```scss
@use "@foundation/tokens/scss" as tokens;

.example {
  color: tokens.$color-text-base;
}
```

Applications that need runtime CSS custom properties can import:

```scss
@use "@foundation/tokens/scss/theme";
```

When publishing a component package outside this repo, publish/install `@foundation/tokens` with it so the component styles keep the same token contract.

## Running unit tests

To execute unit tests with the [Vitest](https://vitest.dev/) test runner, use the following command:

```bash
ng test
```

## Running end-to-end tests

For end-to-end (e2e) testing, run:

```bash
ng e2e
```

Angular CLI does not come with an end-to-end testing framework by default. You can choose one that suits your needs.

## Additional Resources

For more information on using the Angular CLI, including detailed command references, visit the [Angular CLI Overview and Command Reference](https://angular.dev/tools/cli) page.
