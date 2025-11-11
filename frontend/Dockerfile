# Stage 1: Build the Svelte application using Node 22 (latest at the time of writing)
FROM node:22-alpine AS build

WORKDIR /app

# Copy package.json and install dependencies
COPY package*.json ./
RUN npm install

# Copy the rest of the source code and build the app
# The output directory of 'npm run build' is usually 'dist' or 'build'
COPY . .
RUN npm run build

# Stage 2: Serve the application with Nginx
# Use a minimal Nginx image
FROM nginx:alpine

# Copy the built static files from the 'build' stage to Nginx's serving directory
# Make sure to use the correct build output directory if it's not 'dist'
COPY --from=build /app/dist /usr/share/nginx/html

# Expose the port Nginx listens on (default 80)
EXPOSE 80

# Nginx starts automatically by default
CMD ["nginx", "-g", "daemon off;"]