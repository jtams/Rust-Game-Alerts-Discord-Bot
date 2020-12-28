const fs = require("fs");

if (!fs.existsSync("./configs")) fs.mkdirSync("./configs");

CONFIG = {
    setup: "not complete",
    active: false,
    channelId: "0",
    channelName: null,
    battlemetrics: 0000000,
    time: 10000,
    alert: "least",
};

if (!fs.existsSync("./configs/config.json")) {
    fs.writeFileSync("./configs/config.json", JSON.stringify(CONFIG));
} else {
    CONFIG = JSON.parse(fs.readFileSync("./configs/config.json"));
    CONFIG.active = false;
    if (CONFIG.setup == "complete") {
        User = require("./user");
    }
}

if (!fs.existsSync("./configs/users.json")) {
    fs.writeFileSync("./configs/users.json", "[]");
} else {
    userData = JSON.parse(fs.readFileSync("./configs/users.json"));
}

if (!fs.existsSync("./configs/bot.json")) {
    fs.writeFileSync("./configs/bot.json", JSON.stringify(BOT));
    console.log("Your bot.json file is missing data, please open ./configs/bot.json and fill in the required information");
    return;
} else {
    BOT = JSON.parse(fs.readFileSync("./configs/bot.json"));
}

const Discord = require("discord.js");
const axios = require("axios");
const qs = require("qs");
const { verify } = require("crypto");

let USERS = [];
let pending = [];

var onlineCount;
var update = true;
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

    if (msg.includes("battlemetrics.com/servers/rust") && !CONFIG.active) {
        msg = msg.replace(/[^0-9]/g, "");

        let config = {
            method: "get",
            url: `https://api.battlemetrics.com/servers/${msg}`,
            headers: {
                Authorization: BOT.api_auth_key,
            },
            data: "",
            timeout: 10000,
        };

        msgData.channel.send(`Verifying Battlemetrics #${msg} (this won't take more than 10 seconds)...`);
        axios(config)
            .then(function (response) {
                CONFIG.battlemetrics = msg;
                msgData.channel.send(`I will use server ID: ${msg}`);
                msgData.channel.send(`SETUP COMPLETE`);
                CONFIG.setup = "complete";
                msgData.channel.send(
                    "```COMMANDS:\n    !add username1, username2, etc      =   Adds username(s) to alert\n\n    !remove [username]                  =   Removes username to alert. One at a time\n\n    !list                               =   Lists users that I alert you about\n\n    !stop                               =   Stops alerting\n\n    !start                              =   Starts alerting\n\n    !status                             =   Show update status\n\n    !time [milliseconds]                =   Sets how often online status updates. Default 10000 (10 seconds). No faster than 10 seconds.\n\n    !alert [most/online/offline/least]  =   Changes when Rust Sends discord alert. Most: Anytime someone gets online of offline. Online: Only alerts when user gets online. Offline: Only alerts when user gets offlines. Least: Only alerts when neccessary such as the intial alert status; bot will just edit that message. Default = least\n\n    !config                            =    Returns all your configurations```"
                );
                save();
                User = require("./user");
            })
            .catch(function (error) {
                msgData.channel.send(
                    `Can't connect to link. Are you sure it's correct? If you can't connnect to it, it's likely incorrect: https://www.battlemetrics.com/servers/rust/${msg}`
                );
                console.log(
                    "Setup link failed. Either it doens't exist or there's a network problem. Try connecting to it yourself: https://api.battlemetrics.com/servers/" +
                        msg
                );
            });
        return;
    }

    if (msg == "!help") {
        msgData.channel.send(
            "```COMMANDS:\n    !add username1, username2, etc       =   Adds username(s) to alert\n\n    !add ally username1, username2, etc  =   Adds username(s) to allies\n\n    !add enemy username1, username2, etc =   Adds username(s) to enemies\n\n    !add squad username1, username2, etc =   Adds username(s) to squad\n\n    !remove [username]                   =   Removes username to alert. One at a time\n\n    !list                                =   Lists users that I alert you about\n\n    !stop                                =   Stops alerting\n\n    !start                               =   Starts alerting\n\n    !status                              =   Show update status\n\n    !time [milliseconds]                 =   Sets how often online status updates. Default 10000 (10 seconds). No faster than 10 seconds.\n\n    !alert [most/online/offline/least]   =   Changes when Rust Sends discord alert. Most: Anytime someone gets online of offline. Online: Only alerts when user gets online. Offline: Only alerts when user gets offlines. Least: Only alerts when neccessary such as the intial alert status; bot will just edit that message. Default = least\n\n    !config                              =   Returns all your configurations```"
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
        save();
    }

    if (msg.substr(0, 4) == "!add") {
        let type;
        let names = msg;
        if (msg.substr(5, 4).toLowerCase() == "ally" || msg.substr(5, 6).toLowerCase() == "friend") {
            type = "ally";
            names = msg.split(",");
            names[0] = names[0].split(" ").splice(2).join(" ");
        } else if (msg.substr(5, 5).toLowerCase() == "enemy") {
            type = "enemy";
            names = msg.split(",");
            names[0] = names[0].split(" ").splice(2).join(" ");
        } else if (msg.substr(5, 5).toLowerCase() == "squad") {
            type = "squad";
            names = msg.split(",");
            names[0] = names[0].split(" ").splice(2).join(" ");
        } else {
            type = "other";
            names = msg.split(",");
            names[0] = names[0].substr(4);
        }
        names.forEach((name, index) => {
            name = name.trim().toLowerCase();
            names[index] = name;
            if (USERS.some((i) => i.name == name)) {
                names.splice(index, 1);
            }
        });
        if (names.length == 0) {
            msgData.channel.send(`No users added`);
            return;
        }

        added = [];
        addedPending = [];
        verifyUsers = new Promise((resolve, reject) => {
            names.forEach((name, index) => {
                let user = new User();
                user.name = name;
                user.type = type;
                user.updateID()
                    .then((results) => {
                        if (user.id == 0) {
                            user.status = "1";
                            pending.push(user);
                            addedPending.push(user.name);
                        } else {
                            USERS.push(user);
                            added.push(user.name);
                            save();
                        }
                        if (added.length + addedPending.length == names.length) {
                            resolve();
                        }
                    })
                    .catch((err) => {
                        msgData.channel.send(`There was a problem adding user`);
                        needClear = true;
                        reject();
                    });
            });
        });
        verifyUsers.then(() => {
            if (added.length > 0) {
                msgData.channel.send(`Added ${added} to ${type}`);
            }
            if (addedPending.length > 0) {
                msgData.channel.send(
                    `Added ${addedPending} to pending (I can't identify them on the server currently. This means they are offline or their name is spelt incorrectly. Once they get online; they will be verified and moved from pending.)`
                );
            }
            needClear = true;
            if (CONFIG.active) checker();
            return;
        });
    }

    if (msg == "!list") {
        let returnMsg = "";
        if (USERS.length == 0) {
            msgData.channel.send("No users added. Use: !add [username], [username], [username]");
            return;
        }
        USERS.forEach((user) => {
            returnMsg += user.name + ", ";
        });
        msgData.channel.send(returnMsg.slice(0, returnMsg.length - 2));
        needClear = true;
    }

    if (msg.substr(0, 7) == "!remove") {
        var removed = "";
        if (USERS.length == 0) {
            msgData.channel.send("No users have been added yet. Use: !add [username], [username], [username]");
            return;
        }
        if (msg.substr(8) == "all") {
            USERS = [];
            save();
            msgData.channel.send("Removed all users");
            return;
        }

        let names = msg.split(",");
        names[0] = names[0].substr(7);
        let removal = new Promise((resolve, reject) => {
            names.forEach((name, index) => {
                name = name.trim().toLowerCase();
                USERS.forEach((user, i) => {
                    if (user.name == name) {
                        USERS.splice(i, 1);
                        removed += `${name}, `;
                    }
                });
                pending.forEach((user, j) => {
                    if (user.name == name) {
                        pending.splice(j, 1);
                        removed += `${name}, `;
                    }
                });
                if (index == names.length - 1) resolve();
            });
        });

        removal.then(() => {
            if (removed == "") {
                msgData.channel.send("No one removed");
            } else {
                msgData.channel.send(`Removed ${removed}`);
            }
            save();
            needClear = true;
            if (CONFIG.active) checker();
        });
    }

    if (msg == "!config") {
        msgData.channel.send(
            `\`\`\`CHANNEL_ID: ${CONFIG.channelId}\nCHANNEL_NAME: ${CONFIG.channelName}\nBATTLEMETRICS: ${CONFIG.battlemetrics}\nTIME: ${CONFIG.time}\nALERT: ${CONFIG.alert}\`\`\``
        );
        needClear = true;
    }

    if (msg == "!start") {
        if (USERS.length == 0 && pending.length == 0) {
            msgData.channel.send("No users added. Use: !add [username], [username], [username]");
            return;
        }
        needClear = false;
        clear(msgData);
        currentOnline = onlineCount;
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

async function userInitialization() {
    for (var i = 0; i < userData.length; i++) {
        user = userData[i];
        k = new User();
        k.name = user.name;
        k.id = user.id;
        k.online = user.online;
        k.type = user.type;
        k.status = user.status;
        await k.updateOnlineStatus().then((result) => {
            USERS.push(k);
            if (i == userData.length - 1) return;
        });
    }
}

console.log("Initializing Users");
userInitialization().then(() => {
    console.log("Users Initialized");
    client.login(BOT.token);
});
// --------------------------------------------FUNCTIONS---------------------------------------

var onlineCount = 0;
function checker() {
    checkPlayerCount().then((res, err) => {
        if (onlineCount != res) {
            update = true;
            onlineCount = res;
        }

        if (pending.length > 0) {
            pending.forEach((user, index) => {
                user.updateID().then((results) => {
                    if (user.id == 0) {
                        user.status = "1";
                    } else {
                        pending.splice(index, 1);
                        USERS.push(user);
                        save();
                        update = true;
                    }
                });
            });
        }

        if (needClear && updateMessage != "") {
            updateMessage.delete().then().catch(console.error);
            updateMessage = "";
            update = true;
            needClear = false;
        }

        USERS.forEach((user) => {
            user.updateOnlineStatus().then((result) => {
                if (result > 0) {
                    update = true;
                }
            });
        });

        if (update) {
            let date = new Date();
            let msg =
                "https://www.battlemetrics.com/servers/rust/" +
                CONFIG.battlemetrics +
                "\n```diff\nLAST UPDATE: " +
                new Date().toLocaleString("en-US", { day: "numeric", month: "numeric", hour: "numeric", minute: "numeric", hour12: true }) +
                "\n\nCURRENTLY ONLINE: " +
                onlineCount +
                "\n\n";

            let squadMsg = "";
            let enemiesMsg = "";
            let otherMsg = "";
            let alliesMsg = "";
            let pendingMsg = "";

            USERS.forEach((user) => {
                let addition;
                if (user.online) addition = "+ ";
                if (!user.online) addition = "- ";
                if (user.type == "other") otherMsg += `${addition}${user.name.toUpperCase()}\n`;
                if (user.type == "squad") squadMsg += `${addition}${user.name.toUpperCase()}\n`;
                if (user.type == "ally") alliesMsg += `${addition}${user.name.toUpperCase()}\n`;
                if (user.type == "enemy") enemiesMsg += `${addition}${user.name.toUpperCase()}\n`;
            });

            if (pending.length > 0) {
                pending.forEach((u) => {
                    pendingMsg += `- ${u.name}\n`;
                });
            }

            if (squadMsg != "") msg += `SQUAD: \n${squadMsg}\n`;
            if (enemiesMsg != "") msg += `ENEMIES: \n${enemiesMsg}\n`;
            if (alliesMsg != "") msg += `ALLIES: \n${alliesMsg}\n`;
            if (otherMsg != "") msg += `OTHERS: \n${otherMsg}\n`;
            if (pendingMsg != "") msg += `PENDING (Verifying user on server. If you typed the name correctly; they're offline): \n${pendingMsg}\n`;
            msg += "\n```";

            if (updateMessage == "") {
                client.channels.cache
                    .get(CONFIG.channelId)
                    .send(msg)
                    .then((upMsg) => {
                        updateMessage = upMsg;
                    });
            } else {
                updateMessage.edit(msg);
            }

            if (CONFIG.alert == "most") {
                client.channels.cache
                    .get(CONFIG.channelId)
                    .send("UPDATE")
                    .then((upMsg) => upMsg.delete().catch(console.error));
            }
        }
        update = false;
        return;
    });
}

function checkPlayerCount(metrics = CONFIG.battlemetrics) {
    return new Promise((resolve, reject) => {
        var data = "";

        var config = {
            method: "get",
            url: `https://api.battlemetrics.com/servers/${metrics}?include=player`,
            headers: {
                Authorization: BOT.api_auth_key,
            },
            data: data,
        };

        axios(config)
            .then((response) => {
                resolve(response.data.included.length);
            })
            .catch((err) => {
                console.log("Error getting player count");
                reject(err);
            });
    });
}

async function clear(msg) {
    const fetched = await msg.channel.messages.fetch({ limit: 99 });
    msg.channel.bulkDelete(fetched).then().catch(console.error);
}
