const Discord = require("discord.js");
const actionParser = require("./src/actionParser");
const Loop = require("./src/loop");
const pending = require("./src/schemas/pending");
require("dotenv").config();

const client = new Discord.Client();

client.on("ready", () => {
    console.log("Bot online");
});

const keywords = ["online", "offline", "delete", "modify", "add", "remove", "last", "on", "off"];

client.on("message", (message) => {
    if (message.author.bot) return; //Return if bot

    for (var i = 0; i < keywords.length; i++) {
        //Message must include keyword and be 5 characters longer
        if (message.content.includes(keywords[i]) && message.content.length > keywords[i].length + 5) {
            actionParser(message);
            break;
        }
    }

    if (message.content.toLowerCase() == "start") {
        const loop = new Loop({
            config: require("./src/config.json"),
            db: {
                Users: require("./src/schemas/user"),
                Pending: require("./src/schemas/pending"),
            },
        });
        console.log("Starting...");

        loop.updatePendingList();
    }

    if (message.content == "test") {
        message.channel.send(JSON.stringify(prompt, null, 2));
    }
});

client.login(process.env.DISCORD_BOT_TOKEN);
