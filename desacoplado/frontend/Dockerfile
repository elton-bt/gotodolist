FROM nginxinc/nginx-unprivileged:1.27-alpine

COPY desacoplado/frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY desacoplado/frontend/index.html /usr/share/nginx/html/index.html
COPY desacoplado/frontend/styles.css /usr/share/nginx/html/styles.css
COPY desacoplado/frontend/app.js /usr/share/nginx/html/app.js

EXPOSE 8080
