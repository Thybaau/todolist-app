FROM node:21.6.1 AS builder

WORKDIR /usr/src/app
COPY ./client .
WORKDIR /usr/src/app/client
RUN npm install
RUN npm run build

FROM nginx
RUN rm /etc/nginx/conf.d/default.conf
COPY --from=builder /usr/src/app/dist /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/nginx.conf

EXPOSE 80