# Stage 1: Build the Svelte application
FROM node:25-alpine AS builder

WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./

RUN npm ci

# Copy source code and build the application
COPY . .
RUN npm run build

# Stage 2: Serve static files
FROM node:25-alpine

WORKDIR /app

RUN npm install -g serve

COPY --from=builder /app/build ./build

EXPOSE 3000

CMD ["serve", "-s", "build", "-l", "3000"]
