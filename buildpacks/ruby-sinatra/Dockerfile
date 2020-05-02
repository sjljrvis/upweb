# FROM nginx:1.17.9-alpine
FROM ruby:2.3-alpine

RUN apk update
RUN apk upgrade
RUN apk add curl wget bash nginx
RUN rm -rf /var/cache/apk/*

RUN rm /etc/nginx/nginx.conf

COPY nginx.conf /etc/nginx/nginx.conf
# COPY docker-entrypoint.sh /usr/local/bin/

# RUN ln -s /usr/local/bin/docker-entrypoint.sh / 

RUN mkdir /usr/app
WORKDIR /usr/app
COPY Gemfile /usr/app/
RUN gem install bundler -v "~>2.0"
RUN bundle install
COPY . /usr/app
CMD nginx && bundle exec rackup -o 0.0.0.0 -p 9292