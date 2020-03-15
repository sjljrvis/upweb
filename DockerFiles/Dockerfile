FROM nginx:alpine

RUN rm /etc/nginx/conf.d/default.conf
COPY Default/index.html /usr/share/nginx/html
COPY Default/assets /usr/share/nginx/html/assets
COPY Default/default.conf /etc/nginx/conf.d/default.conf