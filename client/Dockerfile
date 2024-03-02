FROM node:21.6.1

WORKDIR /usr/src/app

COPY ./client/package.json ./
COPY ./client/package-lock.json .

RUN npm install

COPY ./client .

RUN npm run build

EXPOSE 3000
CMD [ "npm", "run", "dev" ]