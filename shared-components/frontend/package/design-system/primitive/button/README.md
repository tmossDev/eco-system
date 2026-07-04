# @primitive/button

`@primitive/button` is a primitive package. It provides the shared button affordance used by patterns and applications.

Keep this package generic and token-driven. Add constrained variants here when the whole system needs them; keep feature-specific labels, routing, and submit workflows in patterns or feature packages.

## Layer

- Type: primitive
- Depends on: `@foundation/tokens`
- Used by: patterns and applications

## Building

```bash
corepack pnpm build:button
```
