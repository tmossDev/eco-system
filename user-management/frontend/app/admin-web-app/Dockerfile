FROM node:24-alpine AS builder

WORKDIR /workspace

RUN corepack enable

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY shared-components/frontend/package/design-system ./shared-components/frontend/package/design-system
COPY user-management/frontend/app/admin-web-app ./user-management/frontend/app/admin-web-app

RUN pnpm install --frozen-lockfile
RUN pnpm --filter admin-web-app build

FROM nginx:1.27-alpine

COPY --from=builder /workspace/user-management/frontend/app/admin-web-app/dist/admin-web-app/browser /usr/share/nginx/html

EXPOSE 80
