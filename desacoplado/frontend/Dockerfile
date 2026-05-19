FROM nginxinc/nginx-unprivileged:1.27-alpine

USER root

RUN apk add --no-cache gettext \
	&& chown -R 101:0 /usr/share/nginx/html

COPY desacoplado/frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY desacoplado/frontend/index.html /usr/share/nginx/html/index.html
COPY desacoplado/frontend/styles.css /usr/share/nginx/html/styles.css
COPY desacoplado/frontend/app.js /usr/share/nginx/html/app.js
COPY desacoplado/frontend/config.js /usr/share/nginx/html/config.js
COPY desacoplado/frontend/config.js.template /etc/gotodolist/config.js.template

USER 101

EXPOSE 8080

CMD ["/bin/sh", "-c", "export GOTODOLIST_API_BASE_URL=\"${GOTODOLIST_API_BASE_URL:-${FRONTEND_API_BASE_URL:-}}\" && envsubst '$GOTODOLIST_API_BASE_URL' < /etc/gotodolist/config.js.template > /usr/share/nginx/html/config.js && exec nginx -g 'daemon off;'"]
