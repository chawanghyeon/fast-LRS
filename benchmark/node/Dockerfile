FROM node:20

WORKDIR /app

COPY server.js ./

RUN npm init -y && npm install express pg uuid

CMD ["node", "server.js"]
