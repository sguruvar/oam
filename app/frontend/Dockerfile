FROM node:14

WORKDIR /usr/src/app

COPY package.json package-lock.json ./
RUN npm ci --only=production

COPY . .

EXPOSE 3000
CMD [ "node", "server.js" ]