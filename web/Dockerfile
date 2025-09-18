# Build Stage
FROM node:22.19.0-alpine AS builder
# Set working directory
WORKDIR /app
# Copy package files and install dependencies
COPY web/package.json web/package-lock.json ./
# Install dependencies
RUN npm install
# Copy the rest of the application code
COPY web/ .
# Build the Next.js application
RUN npm run build
# Prune dev dependencies for a smaller production image
RUN npm prune --production

# Production Stage
FROM node:22.19.0-alpine AS runner
# Set working directory
WORKDIR /app
# Copy built application and production dependencies from the builder stage
COPY --from=builder /app ./
# Expose the port
EXPOSE 3000
# Start the Next.js application
CMD ["npm", "start"]