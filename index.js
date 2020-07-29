const fs = require("fs");
const Discord = require("discord.js");
const axios = require("axios");

let CONFIG = JSON.parse(fs.readFileSync("./configs/config.json"));
let BOT = JSON.parse(fs.readFileSync("./configs/bot.json"));
let USERS = JSON.parse(fs.readFileSync("./configs/users.json"));

const ONLINE = [];
var update = false;
var updateMessage = "";
var needClear = false;
var loop;

function save() {
    fs.writeFileSync("./configs/config.json", JSON.stringify(CONFIG));
    fs.writeFileSync("./configs/users.json", JSON.stringify(USERS));
}

client = new Discord.Client();

client.on("ready", () => {
    console.log("Rust Alert Bot Online");
});

default_config = {
    setup: "not complete",
    active: false,
    channelId: "",
    channelName: "",
    battlemetrics: "",
    time: 10000,
    alert: "least",
};

//msg.channel.id or .name

client.on("message", (msg) => {
    if (msg.author.bot) return;
    var msgData = msg;
    msg = msg.content.toLowerCase();
    if ((CONFIG.setup == "waiting" && msg != "yes") || (CONFIG.setup == "waiting" && msg == "no")) {
        CONFIG.setup = "complete";
        save();
    }
    if (CONFIG.setup == "waiting" && msg == "yes") {
        msg = "!setup";
        CONFIG.setup = "not complete";
        fs.writeFileSync("./configs/config.json", JSON.stringify(default_config));
        CONFIG = JSON.parse(fs.readFileSync("./configs/config.json"));
    }

    if (msg === "!setup") {
        if (CONFIG.active == true) {
            clearInterval(loop);
            CONFIG.active = false;
            save();
            msgData.channel.send(`STATUS: STOPPED`);
        }
        if (CONFIG.setup == "complete") {
            msgData.channel.send("You've already completed the setup. This will reset your configuration. Are you sure? Yes or No?");
            CONFIG.setup = "waiting";
            return;
        }
        CONFIG.channelId = msgData.channel.id;
        CONFIG.channelName = msgData.channel.name;
        CONFIG.setup = "active";
        save();
        msgData.channel.send("I will use this channel for alerts from now on.");
        msgData.channel.send("Please paste the link to the battlemetrics.com page for desired server. EX: https://www.battlemetrics.com/servers/rust/99999999");
    }

    if (msg.includes("battlemetrics.com/servers/rust")) {
        msg = msg.replace(/[^0-9]/g, "");
        msg = parseInt(msg);
        if (getMetrics(msg) == false) {
            msgData.channel.send(`Can't connect to link. Are you sure it's correct? Please type it again.`);
        } else {
            CONFIG.battlemetrics = msg;
            msgData.channel.send(`I will use server ID: ${msg}`);
            msgData.channel.send(`SETUP COMPLETE`);
            CONFIG.setup = "complete";
            msgData.channel.send(
                "```COMMANDS:\n    !add username1, username2, etc      =   Adds username(s) to alert\n\n    !remove [username]                  =   Removes username to alert. One at a time\n\n    !list                               =   Lists users that I alert you about\n\n    !stop                               =   Stops alerting\n\n    !start                              =   Starts alerting\n\n    !status                             =   Show update status\n\n    !time [milliseconds]                =   Sets how often online status updates. Default 10000 (10 seconds). No faster than 10 seconds.\n\n    !alert [most/online/offline/least]  =   Changes when Rust Sends discord alert. Most: Anytime someone gets online of offline. Online: Only alerts when user gets online. Offline: Only alerts when user gets offlines. Least: Only alerts when neccessary such as the intial alert status; bot will just edit that message. Default = least\n\n    !config                            =    Returns all your configurations```"
            );
            save();
        }
        return;
    }

    if (msg == "!help") {
        msgData.channel.send(
            "```COMMANDS:\n    !add username1, username2, etc      =   Adds username(s) to alert\n\n    !remove [username]                  =   Removes username to alert. One at a time\n\n    !list                               =   Lists users that I alert you about\n\n    !stop                               =   Stops alerting\n\n    !start                              =   Starts alerting\n\n    !status                             =   Show update status\n\n    !time [milliseconds]                =   Sets how often online status updates. Default 10000 (10 seconds). No faster than 10 seconds.\n\n    !alert [most/online/offline/least]  =   Changes when Rust Sends discord alert. Most: Anytime someone gets online of offline. Online: Only alerts when user gets online. Offline: Only alerts when user gets offlines. Least: Only alerts when neccessary such as the intial alert status; bot will just edit that message. Default = least\n\n    !config                             =   Returns all your configurations```"
        );
        needClear = true;
    }

    if (msg == "!status") {
        msgData.channel.send("STATUS: " + CONFIG.active);
        needClear = true;
    }

    if (msg.substr(0, 5) == "!time") {
        if (msg == "!time") {
            msgData.channel.send(`Time between data pulls is set to **${CONFIG.time}ms**`);
            needClear = true;
            return;
        }
        if (typeof parseInt(msg.substr(5)) == "number") {
            CONFIG.time = parseInt(msg.substr(5));
            save();
            msgData.channel.send(`Updated time to ${CONFIG.time}`);
            if (CONFIG.active) {
                updateMessage = "";
                clearInterval(loop);
                CONFIG.active = false;
                save();
                msgData.channel.send(`STATUS: STOPPED`);
                return;
            }
        } else {
            msgData.channel.send(`Error. Time not updated. Time = ${CONFIG.time}`);
        }
        if (loop != undefined) {
            clearInterval(loop);
            update = true;
            checker();
            loop = setInterval(checker, CONFIG.time);
        }
        needClear = true;
    }

    if (msg.substr(0, 6) == "!alert") {
        if (msg == "!alert") {
            msgData.channel.send(`Alert status set to **${CONFIG.alert}**`);
            needClear = true;
            return;
        }
        if (msg.substr(7) === CONFIG.alert) {
            msgData.channel.send(`Alert already set to ${msg.substr(7)}`);
            needClear = true;
            return;
        }
        if (msg.substr(7) === "most") {
            CONFIG.alert = "most";
        }
        if (msg.substr(7) === "least") {
            CONFIG.alert = "least";
        }
        if (msg.substr(7) === "online") {
            CONFIG.alert = "online";
        }
        if (msg.substr(7) === "offline") {
            CONFIG.alert = "offline";
        }
        msgData.channel.send(`Alert changed to ${msg.substr(7)}`);
        needClear = true;
    }

    if (msg.substr(0, 4) == "!add") {
        let names = msg.split(",");
        names[0] = names[0].substr(4);
        names.forEach((name, index) => {
            name = name.trim();
            names[index] = name;
            if (USERS.includes(name)) {
                names.splice(index, 1);
            }
        });
        if (names.length == 0) {
            msgData.channel.send(`No users added`);
            return;
        }
        USERS.push(...names);
        save();
        msgData.channel.send(`Added ${names}`);
        needClear = true;
    }

    if (msg == "!list") {
        let returnMsg = "";
        if (USERS.length == 0) {
            msgData.channel.send("No users added. Use: !add [username], [username], [username]");
            return;
        }
        USERS.forEach((name) => {
            returnMsg += name + ", ";
        });
        msgData.channel.send(returnMsg);
        needClear = true;
    }

    if (msg.substr(0, 7) == "!remove") {
        if (USERS.length == 0) {
            msgData.channel.send("No users added. Use: !add [username], [username], [username]");
            return;
        }
        if (msg.substr(8) == "all") {
            USERS = [];
            save();
            msgData.channel.send("Removed all users");
            return;
        }
        let names = msg.split(",");
        let removed = [];
        names[0] = names[0].substr(7);
        names.forEach((name) => {
            name = name.trim();
            USERS.forEach((username, index) => {
                if (username.includes(name)) {
                    USERS.splice(index, 1);
                    removed.push(name);
                }
            });
        });
        save();
        if (removed.length == 0) {
            msgData.channel.send("No one removed");
        } else {
            msgData.channel.send(`Removed ${removed}`);
        }
        needClear = true;
    }

    if (msg == "!config") {
        msgData.channel.send(
            `\`\`\`CHANNEL_ID: ${CONFIG.channelId}\nCHANNEL_NAME: ${CONFIG.channelName}\nBATTLEMETRICS: ${CONFIG.battlemetrics}\nTIME: ${CONFIG.time}\nALERT: ${CONFIG.alert}\`\`\``
        );
        needClear = true;
    }

    if (msg == "!start") {
        if (USERS.length == 0) {
            msgData.channel.send("No users added. Use: !add [username], [username], [username]");
            return;
        }
        needClear = false;
        clear(msgData);
        CONFIG.active = true;
        save();
        msgData.channel.send(`STATUS: ACTIVE`);
        update = true;
        checker();
        loop = setInterval(checker, CONFIG.time);
    }
    if (msg == "!stop") {
        if (CONFIG.active == false) {
            msgData.channel.send("Already stopped");
            return;
        }
        updateMessage = "";
        clear(msgData);
        clearInterval(loop);
        CONFIG.active = false;
        save();
        msgData.channel.send(`STATUS: STOPPED`);
    }

    if (msg == "!clear") {
        clear(msgData);
    }
});

client.login(BOT.token);
// --------------------------------------------FUNCTIONS---------------------------------------

function checker() {
    getOnlineStatus(CONFIG.url).then(() => {
        if (needClear && updateMessage != "") {
            updateMessage.delete().then().catch(console.error);
            updateMessage = "";
            update = true;
            needClear = false;
        }
        if (update) {
            let date = new Date();
            let msg = "```diff\nLAST UPDATE " + date.getHours() + ":" + date.getMinutes() + ":" + date.getSeconds() + "\n";
            USERS.forEach((name) => {
                if (ONLINE.includes(name)) {
                    msg += "+ " + name.toUpperCase() + "\n";
                    if (CONFIG.alert == "online") {
                        client.channels.cache
                            .get(CONFIG.channelId)
                            .send("ONLINE UPDATE")
                            .then((upMsg) => upMsg.delete().catch(console.error));
                    }
                } else {
                    msg += "- " + name.toUpperCase() + "\n";
                    if (CONFIG.alert == "offline") {
                        client.channels.cache
                            .get(CONFIG.channelId)
                            .send("OFFLINE UPDATE")
                            .then((upMsg) => upMsg.delete().catch(console.error));
                    }
                }
            });
            msg += "\n```";

            if (updateMessage == "") {
                client.channels.cache
                    .get(CONFIG.channelId)
                    .send(msg)
                    .then((upMsg) => (updateMessage = upMsg));
            } else {
                updateMessage.edit(msg);
            }
            update = false;

            if (CONFIG.alert == "most") {
                client.channels.cache
                    .get(CONFIG.channelId)
                    .send("UPDATE")
                    .then((upMsg) => upMsg.delete().catch(console.error));
            }
        }
        return;
    });
}

async function clear(msg) {
    const fetched = await msg.channel.messages.fetch({ limit: 99 });
    msg.channel.bulkDelete(fetched).then().catch(console.error);
}

function getMetrics(serverID) {
    if (serverID == undefined) {
        if (CONFIG.battlemetrics == undefined) {
            return false;
        } else {
            serverID = CONFIG.battlemetrics;
        }
    }

    var data = "";

    var config = {
        method: "get",
        url: `https://api.battlemetrics.com/servers/${serverID}?include=player`,
        headers: {
            Authorization: BOT.api_auth_key,
        },
        data: data,
    };
    return new Promise((resolve) => {
        axios(config)
            .then((response) => {
                resolve(response.data);
            })
            .catch((err) => {
                return false;
            });
    });
}

function getOnlineStatus(serverID) {
    if (serverID == undefined) {
        if (CONFIG.battlemetrics == undefined) {
            return false;
        } else {
            serverID = CONFIG.battlemetrics;
        }
    }
    return new Promise((resolve) => {
        getMetrics(serverID)
            .then((data) => {
                data = data.included;
                USERS.forEach((name) => {
                    let nameFound = false;
                    for (let i = 0; i < data.length; i++) {
                        let users = data[i];
                        if (users.attributes.name.toLowerCase().includes(name)) {
                            if (!ONLINE.includes(name)) {
                                ONLINE.push(name);
                                nameFound = true;
                                update = true;
                                break;
                            } else {
                                nameFound = true;
                                break;
                            }
                        }
                    }
                    if (ONLINE.includes(name) && !nameFound) {
                        ONLINE.forEach((username, index) => {
                            if (name == username) {
                                ONLINE.splice(index, 1);
                                update = true;
                            }
                        });
                    }
                });
                resolve(true);
            })
            .catch((err) => {
                console.log(err);
            });
    });
}
