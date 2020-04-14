FROM nginx:1.17.9-alpine

RUN apk update && apk add --update nodejs nodejs-npm
RUN rm /etc/nginx/nginx.conf

COPY nginx.conf /etc/nginx/nginx.conf
COPY docker-entrypoint.sh /usr/local/bin/

RUN ln -s /usr/local/bin/docker-entrypoint.sh / 

RUN mkdir -p ~/app
WORKDIR ~/app

ADD . .
RUN npm install
ENTRYPOINT ["sh", "docker-entrypoint.sh"]


