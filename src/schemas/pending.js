const mongoose = require("mongoose");

var pendingSchema = mongoose.Schema({
    username: String,
    serverId: String,
    dateAdded: { type: Date, default: new Date() },
    group: { type: String, default: "Default" },
});

module.exports = mongoose.model("Pending", pendingSchema);
