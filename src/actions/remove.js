module.exports = function (users, message) {
    // Creating message for bot to say
    let str = "Removing ";

    if (users.length > 2) {
        users.forEach((user) => {
            str += user + ", ";
        });
        str = str.trim();
        str = str.slice(0, -1);

        str2 = str.split(",");
        str2[str2.length - 2] += ", and" + str2[str2.length - 1]; //Adds "and last_username to second to last index"
        str2.pop(); //Removes last username
        str = str2.join();
    }

    if (users.length == 2) {
        str += users[0] + " and " + users[1];
    }

    if (users.length == 1) {
        str += users[0];
    }

    message.channel.send(str);
};
