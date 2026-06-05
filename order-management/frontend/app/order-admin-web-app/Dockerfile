FROM nginx:1.27-alpine

COPY order-management/frontend/app/order-admin-web-app/dist/order-admin-web-app/browser /usr/share/nginx/html
COPY order-management/frontend/app/order-admin-web-app/storybook-static /usr/share/nginx/html/storybook
RUN mkdir -p /usr/share/nginx/html/assets

EXPOSE 80
