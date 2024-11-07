FROM node:22.8.0-slim

USER node

WORKDIR /home/node/app

CMD ["npm", "run", "dev"]

# FROM node:20.18.0 AS builder

# LABEL author="armando"

# WORKDIR /app

# COPY package.json package-lock.json ./

# RUN npm install --save-dev eslint eslint-config-next @typescript-eslint/eslint-plugin @typescript-eslint/parser --force


# COPY . .

# RUN npm run build


# FROM node:20.18.0 AS production


# WORKDIR /app


# COPY --from=builder /app ./


# EXPOSE 3000


# CMD ["npm", "start"]

