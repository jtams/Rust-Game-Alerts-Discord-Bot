require("dotenv").config();
require("mongoose").connect(process.env.DB, { useUnifiedTopology: true, useNewUrlParser: true, useCreateIndex: true });

module.exports.add = require("./add.js");
module.exports.remove = require("./remove.js");
