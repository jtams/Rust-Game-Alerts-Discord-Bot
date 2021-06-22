const openai = require("openai-node");
const actions = require("./actions");

module.exports = function (message) {
    parseMessageWithAI(message.content.trim())
        .then((actionList) => {
            actionList.forEach((action) => {
                switch (action.action) {
                    case "add":
                        actions.add(action.users, message);
                        break;
                    case "remove":
                        actions.remove(action.users, message);
                        break;
                }
            });
        })
        .catch((e) => {
            console.log(e);
            message.channel.send("I'm not sure what you want me to do.");
        });
};

function parseMessageWithAI(msg) {
    return new Promise((resolve, reject) => {
        console.log("Request made to OpenAI");
        resolve([{ action: "add", users: ["The Higher Realm"] }]);
        return;
        openai.api_key = process.env.OPENAI_API;
        openai.Completion.create({
            engine: "curie",
            prompt: `Extract user names and actions from list and put them in a JavaScript object.\nActions: currently_online, last_online, add, remove, change\n\nQ: add noobmaster69 The Higher Realm 420noscope and raymond345\nA: [{\"action\": \"add\", \"users\": [\"noobmaster69\", \"The Higher Realm\", \"420noscope\", \"raymond345\"]}]\n###\nQ: When was SHAWSKEE155 last online?\nA: [{\"action\": \"last_online\", \"users\": [\"SHAWSKEE155\"]}]\n###\nQ: add wisps and Lucifer\nA: [{\"action\": \"add\", \"users\": [\"Lucifer\", \"wisps\"]}]\n###\nQ: remove dryguy, otherguy, and sombody123 and then add supeuser23\nA: [{\"action\": \"remove\", \"users\": [\"dryguy\", \"otherguy\", \"sombody123\"]}, {\"action\": \"add\", \"users\": [\"superuser23\"]}]\n###\nQ: remove otherdude and noobmaster33\nA: [{\"action\": \"remove\", \"users\": [\"otherdude \", \"noobmaster33\"]}]\n###\nQ:${msg}\nA:`,
            // prompt: `Extract user names and actions from list and put them in a JavaScript object.\nActions: currently_online, last_online, add, remove, change\n\nQ: add noobmaster69 The Higher Realm 420noscope and raymond345\nA: {\"action\": \"add\", \"users\": [\"noobmaster69\", \"The Higher Realm\", \"420noscope\", \"raymond345\"]}\n###\nQ: When was SHAWSKEE155 last online?\nA: {\"action\": \"last_online\", \"users\": [\"SHAWSKEE155\"]}\n###\nQ: add wisps and Lucifer\nA: {\"action\": \"add\", \"users\": [\"Lucifer\", \"wisps\"]}\n###\nQ: ${this.originalMessage.content}\nA:`,
            temperature: 0.7,
            max_tokens: 100,
            top_p: 1,
            frequency_penalty: 1,
            presence_penalty: 0,
            stop: ["###", "\n"],
        }).then((res) => {
            console.log(res);
            try {
                let actionList = JSON.parse(res.choices[0].text.trim());
                console.log(actionList);
                resolve(actionList);
            } catch (e) {
                console.log(e);
                reject(e);
            }
        });
    });
}
