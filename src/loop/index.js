// const User = require("../schemas/user");
// const Pending = require("../schemas/pending");
const request = require("request");

module.exports = class {
    constructor({
        config,
        serverId,
        updateInterval = 30,
        db = {
            Users,
            Pending,
        },
    } = {}) {
        this.config = config;
        this.serverId = serverId;
        this.updateInterval = updateInterval * 1000;

        if (config) {
            this.serverId = config.serverId;
            this.updateInterval = config.updateInterval * 1000;
        }
        console.log(this.updateInterval);

        this.db = db;

        this.pendingList = [];
        this.currentlyOnline = [];

        this.loop = setInterval(() => {
            console.log("Loop Start");
            this.checkPending().then(() => {
                console.log("Loop end");
            });
        }, this.updateInterval);
    }

    checkPending() {
        return new Promise((resolve, reject) => {
            this.updatePendingList().then(() => {
                this.updateOnlineList().then((list) => {
                    this.currentlyOnline = [...list];
                    var userCount = this.currentlyOnline.length;
                    this.currentlyOnline.forEach((user) => {
                        userCount -= 1;
                        var inside = this.pendingList.find((u) => {
                            if (u.username == user.username) return true;
                        });
                        if (inside != undefined && inside.username == user.username) {
                            this.db.Pending.deleteOne({ username: user.username }).then(() => {
                                let u = new this.db.Users({
                                    username: user.username,
                                    id: user.id,
                                    serverId: this.serverId,
                                    group: "Default",
                                    lastOnline: new Date(),
                                });
                                u.save().then(() => {
                                    resolve();
                                });
                            });
                        }
                    });
                    if (userCount == 0) {
                        resolve();
                    }
                });
            });
        });
    }

    updateOnlineList() {
        return new Promise((resolve, reject) => {
            var newOnlineList = [];
            var options = {
                method: "GET",
                url: `https://api.battlemetrics.com/servers/${this.serverId}?include=player`,
                headers: {
                    Authorization: process.env.BATTLE_METRICS_KEY,
                },
            };
            request(options, function (error, response) {
                if (error) throw new Error(error);
                // console.log(response.body);
                try {
                    let included = JSON.parse(response.body).included;
                    // console.log(JSON.stringify(included, null, 2));

                    included.forEach((player) => {
                        if (player.type == "player") {
                            // console.log("Player: " + player.attributes.name);
                            newOnlineList.push({ id: player.attributes.id, username: player.attributes.name });
                        }
                    });
                    // this.currentlyOnline = [...newOnlineList];
                    // console.log("New Player List: " + JSON.stringify(this.currentlyOnline, null, 2));

                    resolve(newOnlineList);
                } catch (e) {
                    reject(e);
                    return;
                }
            });
        });
    }

    updatePendingList() {
        return new Promise((resolve, reject) => {
            this.db.Pending.find({}).then((doc) => {
                if (doc) {
                    this.pendingList = [...doc];
                    resolve();
                }
            });
        });
    }
};
