FROM nginx:1.27-alpine

COPY user-management/frontend/app/admin-web-app/dist/admin-web-app/browser /usr/share/nginx/html
COPY user-management/frontend/app/admin-web-app/storybook-static /usr/share/nginx/html/storybook
RUN mkdir -p /usr/share/nginx/html/assets

EXPOSE 80
