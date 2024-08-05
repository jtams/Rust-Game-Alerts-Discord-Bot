# docker build -t rust-game-alerts-discord-bot .
# or
# docker buildx build -t rust-game-alerts-discord-bot .

# docker run -d -v playertracker-data:/data -e BOT_TOKEN=your_bot_token -e GUILD_ID=your_guild_id -e LOG_LEVEL=INFO -e SAVE_FILE=/data/savefile.json --name discord-playertracker rust-game-alerts-discord-bot

# After testing and verifying it starts up, it's recommended to add --restart unless-stopped to the docker run command to ensure the bot restarts if the container crashes.

# Use the official Golang image as the base image
FROM golang:1.22-alpine

# Install bash and other dependencies
RUN apk update && apk add --no-cache bash

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

RUN go build .

# Command to run the shell script
CMD ["./playertrackerbot"]
