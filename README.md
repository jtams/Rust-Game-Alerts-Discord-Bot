# Rust Game Discord Alerts

**A bot that sends user's online status for the game Rust**

Rust (the game not the language) alerts is a Discord bot that will send a message to a channel with information about what players are online. It uses the Battlemetrics API to get the player list and then checks if the user is online.

![Tracker Message](./screenshot1.png)

## Usage

```
git clone https://github.com/jtams/Rust-Game-Alerts-Discord-Bot
cd Rust-Game-Alerts-Discord-Bot
go build
```

Create a `.env` file in the root directory or add environment variables with the following information:

```env
BOT_TOKEN=[your Discord bot token] [required]
GUILD_ID=[your Discord server ID] [required
LOG_LEVEL=[log level] [optional] [default: INFO] [options: DEBUG, INFO, WARN, ERROR]
SAVE_FILE=[file to save player data] [optional] [default: save_data.json]
```

Run the bot with:
`./playertrackerbot`

I recommend running this as a service so it can run in the background and can auto restart in case of a crash.

## Features

- Automatically updates player's name if they change it. (/info user <username> to view all previous names)
- Seperate users into groups for easier tracking.
- Add notes and locations to groups for additional information.
- Designed to be non-intrusive and will edit the message instead of spamming the channel.

## Commands

### /start <battle_metrics_server_id, optional>

Starts the tracker. Battle Metrics Server ID is requried to run the first time. If it is not provided, the bot will use the last server ID that was used.

### /stop

Stops the tracker.

### /users add <username(s)> <group>

Adds user(s) to the tracker. Seperate multiple names with a comma. The name doesn't need to be exact and is not case sensitive. The search algorithm prioritizes exact matches first, then names that start with the inputted name, and finally names that contain the inputted name.

### /users add-by-id <battle_metrics_player_id> <group>

Similiar to the add command, but uses the Battle Metrics player ID instead of the player name. Useful when there are players with the same name.

### /users remove <username> <group, optional>

Removes user(s) from the tracker. If a group is provided, it will only remove from that group.

### /users move <username> <group>

Moves a user to a different group.

### /groups add <group_name>

Adds a group to the tracker.

### /groups remove <group_name>

Removes a group from the tracker.

### /groups location <location, optional>

Sets the location of the group. This gets added next to the group name in the tracker. Useful for noting the base location of a group. Running this command without location will remove the location.

### /groups notes <notes, optional>

Sets the notes of the group. This gets added below the group name in the tracker. Useful for adding additional information about the group. Running this command without notes will remove the note.

### /info user <username>

Gets information about a user, including their BattleMetrics profile which will show the times they were on.

### /info group <group_name>

Gets information about a group.

## To Do

- [ ] Add a command to get server data from BattleMetrics
- [ ] Option to send a message when a user logs on/off
- [ ] Connect with Rust+ app to message in chat. (Not sure how possible this is).
- [ ] Predict when a user will be online next.
