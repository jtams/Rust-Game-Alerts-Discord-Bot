// const mongoose = require("mongoose");

const User = require("../schemas/user");
const Pending = require("../schemas/pending");
const config = require("../config.json");

module.exports = function (users, message, group = "Default") {
    if (!config.serverId) throw new Error("Missing server ID");
    const addedUsers = [];

    addToPending(users, group, addedUsers).then(() => {
        let str = "Added ";

        if (addedUsers.length > 2) {
            addedUsers.forEach((user) => {
                str += user + ", ";
            });
            str = str.trim();
            str = str.slice(0, -1);

            str2 = str.split(",");
            str2[str2.length - 2] += ", and" + str2[str2.length - 1]; //Adds "and last_username to second to last index"
            str2.pop(); //Removes last username
            str = str2.join();
        }

        if (addedUsers.length == 2) {
            str += addedUsers[0] + " and " + addedUsers[1];
        }

        if (addedUsers.length == 1) {
            str += addedUsers[0];
        }

        if (addedUsers.length == 0) {
            str = "No users added to pending.";
        } else {
            str += " to pending list. Waiting to see if they're online.";
        }

        message.channel.send(str);
    });

    // Creating message for bot to say
};

function addToPending(users, group, addedUsers) {
    return new Promise((resolve, reject) => {
        var usersPending = users.length;
        users.forEach((user, index) => {
            User.findOne({ username: user }).then((doc1) => {
                if (!doc1) {
                    Pending.findOne({ username: user })
                        .then((doc) => {
                            if (!doc) {
                                let u = new Pending({
                                    username: user,
                                    serverId: config.serverId,
                                    group: group,
                                });
                                u.save().then(() => {
                                    addedUsers.push(user);
                                    usersPending -= 1;
                                    if (usersPending == 0) {
                                        resolve();
                                    }
                                });
                            } else {
                                usersPending -= 1;
                            }
                            if (usersPending == 0) {
                                resolve();
                            }
                        })
                        .catch((e) => {
                            console.log(e);
                            usersPending -= 1;
                        });
                } else {
                    usersPending -= 1;
                    if (usersPending == 0) {
                        resolve();
                    }
                }
            });
        });
    });
}
