FROM node:22-alpine

WORKDIR /app

COPY package*.json ./

RUN npm install

COPY . .

RUN npm run build --if-present

RUN npm prune --production

EXPOSE 3000

CMD ["npm", "start"]
