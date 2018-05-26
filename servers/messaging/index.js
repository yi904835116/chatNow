'use strict';

const mongodb = require('mongodb');
const mongoAddr = process.env.DBADDR || '192.168.99.100:27017';
const mongoURL = `mongodb://${mongoAddr}`;
// const mongoURL = "mongodb://" + mongoAddr + "/info_344";
const ChannelStore = require('./models/channels/channel-store.js');
const MessageStore = require('./models/messages/message-store.js');

const express = require('express');
const app = express();
const morgan = require('morgan');

const Channel = require('./models/channels/channel');

const ChannelHandler = require('./handlers/channels_handler.js');
const MessageHandler = require('./handlers/messages_handler.js');

const addr = process.env.ADDR || 'localhost:4000';
const [host,
    port] = addr.split(':');
const portNum = parseInt(port);

(async() => {

    try {
        const client = await mongodb
            .MongoClient
            .connect(mongoURL);

        const db = client.db("info_344");

        //Add global middlewares.
        app.use(morgan("dev"));
        // Parses posted JSON and makes it available from req.body.
        app.use(express.json());

        // All of the following APIs require the user to be authenticated. If the user
        // is not authenticated, respond immediately with the status code 401
        // (Unauthorized).
        app.use((req, res, next) => {
            const userJSON = req.get('X-User');
            if (!userJSON) {
                res.set('Content-Type', 'text/plain');
                res
                    .status(401)
                    .send('no X-User header found in the request');
                // Stop continuing.
                return;
            }
            // Invoke next chained handler if the user is authenticated.
            next();
        });

        // Initialize Mongo stores.
        let channelStore = new ChannelStore(db, 'channels');
        let messageStore = new MessageStore(db, 'messages');

        const defaultChannel = new Channel('general', null, false, null, null);
        const fetchedChannel = await channelStore.getByName(defaultChannel.name);
        // Add the default channel if not found.
        if (!fetchedChannel) {
            const channel = await channelStore.insert(defaultChannel);
        }

        // API resource handlers.
        ChannelHandler(app, channelStore, messageStore);
        MessageHandler(app, messageStore);

        app.listen(portNum, host, () => {
            console.log(`server is listening at http://${addr}`);
        });
    } catch (err) {
        console.log(err)
    }

})();