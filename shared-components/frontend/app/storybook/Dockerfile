FROM node:24-alpine AS builder

WORKDIR /workspace

RUN corepack enable

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY shared-components/frontend/package/design-system ./shared-components/frontend/package/design-system
COPY shared-components/frontend/app/storybook ./shared-components/frontend/app/storybook

RUN pnpm install --frozen-lockfile
RUN pnpm --filter storybook build

FROM nginx:1.27-alpine

COPY --from=builder /workspace/shared-components/frontend/app/storybook/storybook-static /usr/share/nginx/html

EXPOSE 80
