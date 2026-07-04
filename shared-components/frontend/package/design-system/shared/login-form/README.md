# @shared/login-form

`@shared/login-form` is a design-system pattern. It composes primitives from `@primitive/form` and `@primitive/button` into a reusable sign-in form shell.

Keep this package free of application behavior. Auth services, route redirects, permission decisions, and API-specific error handling belong in `auth-features` or the consuming app.

## Layer

- Type: pattern
- Depends on: `@primitive/button`, `@primitive/form`, `@foundation/tokens`
- Should not be imported by: primitive packages

## Building

```bash
corepack pnpm build:login-form
```
