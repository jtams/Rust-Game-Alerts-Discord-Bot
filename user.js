const fs = require("fs");
const axios = require("axios");
const qs = require("qs");
let CONFIG = JSON.parse(fs.readFileSync("./configs/config.json"));
let BOT = JSON.parse(fs.readFileSync("./configs/bot.json"));

class User {
    constructor(id = 0, name = "unknown", online = false, type = "other", status = "ok") {
        this.id = id;
        this.name = name.toLowerCase();
        this.online = online;
        this.type = type;
        this.status = status;
    }

    updateName() {
        var data = qs.stringify({});
        var config = {
            method: "get",
            url: `https://api.battlemetrics.com/players/${this.id}}`,
            headers: {
                Authorization: BOT.api_auth_key,
            },
            data: data,
        };

        axios(config)
            .then(function (response) {
                data = response.data.data;
                this.name = data.attributes.name;
                return 0;
            })
            .catch(function (error) {
                console.log("There was an error when retreiving a users name. Probably a 405. Probably shouldn't worry about it. Probably");
                return 1;
            });
    }

    updateID(serverID = CONFIG.battlemetrics) {
        var data = "";

        var config = {
            method: "get",
            url: `https://api.battlemetrics.com/servers/${serverID}?include=player`,
            headers: {
                Authorization: BOT.api_auth_key,
            },
            data: data,
        };

        return new Promise((resolve, reject) => {
            axios(config)
                .then((response) => {
                    data = response.data.included;
                    for (let i = 0; i < data.length; i++) {
                        var current = data[i].attributes;
                        if (current.name.toLowerCase() == this.name) {
                            this.id = parseInt(current.id);
                            this.online = true;
                            break;
                        } else {
                            this.online = false;
                            if (this.id == 0) {
                                this.status = "never found";
                                resolve("never found");
                            }
                            this.status = "not found";
                            resolve("not found");
                        }
                    }
                    if (this.id == 0) {
                        resolve();
                    } else {
                        resolve("found id");
                    }
                })
                .catch((err) => {
                    console.log("Error getting ID");
                    reject(err);
                });
        });
    }

    updateOnlineStatus(serverID = CONFIG.battlemetrics) {
        let data = "";
        var config = {
            method: "get",
            url: `https://api.battlemetrics.com/players/${this.id}/servers/${serverID}`,
            headers: {
                Authorization: BOT.api_auth_key,
            },
            data: data,
            timeout: 6000,
        };

        var current = this;

        return new Promise((resolve, reject) => {
            axios(config)
                .then(function (response) {
                    data = response.data.data;
                    if (data.attributes.online) {
                        if (current.online) {
                            resolve(0);
                        } else {
                            current.online = true;
                            current.updateName();
                            resolve(1);
                        }
                    } else {
                        if (current.online) {
                            current.online = false;
                            current.updateName();
                            resolve(2);
                        } else {
                            resolve(0);
                        }
                    }
                })
                .catch(function (err) {
                    console.log("There was an error getting update status by ID. Probably a 405. Probably shouldn't worry about it. Probably");
                    reject(err);
                });
        });
    }
}

module.exports = User;
