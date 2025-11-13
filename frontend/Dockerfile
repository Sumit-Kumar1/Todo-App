# Stage 1: Build the Svelte application
FROM node:25-alpine AS builder

WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm ci

# Copy source code and build the application
COPY . .
RUN npm run build

# Stage 2: Serve with Nginx
FROM nginx:stable-alpine

# Copy the built static files from the 'builder' stage to Nginx's serving directory
COPY --from=builder /app/build /usr/share/nginx/html

# Copy custom nginx configuration if needed
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose the port Nginx listens on (default 80)
EXPOSE 80

# Nginx starts automatically by default
CMD ["nginx", "-g", "daemon off;"]