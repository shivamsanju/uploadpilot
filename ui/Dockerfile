# Stage 1: Build the React app with the SDK
FROM node:16-alpine as build

# Set the working directory in the container
WORKDIR /app

# Copy the package.json and package-lock.json files
COPY ./ui/package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application files
COPY ./ui .

# Copy the React SDK into the container (adjust path accordingly)
COPY ./sdk/react ./uploadpilot-react-sdk

# Install and build the React SDK
RUN cd /app/uploadpilot-react-sdk && npm install && npm run build && cd /app

# Ensure the SDK is linked or included properly
RUN npm install /app/uploadpilot-react-sdk

# Build the React app for production
RUN npm run build

# Stage 2: Serve the built app with Nginx
FROM nginx:alpine

# Copy the build files from the previous stage to the Nginx directory
COPY --from=build /app/build /usr/share/nginx/html

# Expose port 80 (Nginx default)
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
