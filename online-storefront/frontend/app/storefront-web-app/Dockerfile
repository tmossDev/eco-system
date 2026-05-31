FROM nginx:1.27-alpine

COPY online-storefront/frontend/app/storefront-web-app/dist/storefront-web-app/browser /usr/share/nginx/html
RUN mkdir -p /usr/share/nginx/html/assets

EXPOSE 80
