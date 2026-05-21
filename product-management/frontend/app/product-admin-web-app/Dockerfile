FROM nginx:1.27-alpine

COPY product-management/frontend/app/product-admin-web-app/dist/product-admin-web-app/browser /usr/share/nginx/html
COPY product-management/frontend/app/product-admin-web-app/storybook-static /usr/share/nginx/html/storybook
RUN mkdir -p /usr/share/nginx/html/assets

EXPOSE 80
