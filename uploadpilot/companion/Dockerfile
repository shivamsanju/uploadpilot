FROM node:18

WORKDIR /app

COPY ./package*.json ./

RUN npm install

COPY . .

EXPOSE 8082

# Run entrypoint script
CMD ["npm", "start"]
