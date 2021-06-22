# Rust Alerts

**A bot that sends user's online status for the game Rust**

Rust (the game not the language) alerts is a discord bot made in JavaScript/Node using the Discord.js library. My friends and I like to play Rust and we make enemies. We would periodically check Steam to see if they were online or not. This bot changed that. It periodically checks battlemetrics.com to see who's on and updates our discord server.

## Usage

I currently don't have the steps to set up yourself. As it relies on a private API from OpenAI but I'll remove that dependency soon.

## Features

Using the setup command will prompt you to enter the Rust Battlemetrics Server ID. Once you paste that the bot will start monitoring the players. You can add and remove players from the tracking list. Once the bot detects that the person is online, it'll edit it's message and turn their name green. It'll also keep track of all the players' IDs so if they change their username it'll update that too.
