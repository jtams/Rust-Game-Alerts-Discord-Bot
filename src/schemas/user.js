const mongoose = require("mongoose");

const userSchema = mongoose.Schema({
    username: String,
    id: String,
    serverId: String,
    group: { type: String, default: "Default" },
    lastOnline: Date,
});

module.exports = mongoose.model("User", userSchema);
