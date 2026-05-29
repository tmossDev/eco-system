# DesignSystem

This project was generated using [Angular CLI](https://github.com/angular/angular-cli) version 21.2.10.

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

Design tokens live in `ds/tokens/tokens` and are generated with Style Dictionary into `@ds/tokens`.

```bash
corepack pnpm build:tokens
```

Components can use the generated SCSS module:

```scss
@use "@ds/tokens/scss" as tokens;

.example {
  color: tokens.$color-text-base;
}
```

Applications that need runtime CSS custom properties can import:

```scss
@use "@ds/tokens/scss/theme";
```

When publishing a component package outside this repo, publish/install `@ds/tokens` with it so the component styles keep the same token contract.

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
