const express = require("express");
const mongodb = require('mongodb');

const Channel = require('./../models/channels/channel');
const Message = require('./../models/messages/message');

const sendToMQ = require('./messages_queue');

function MessageHandler(app, channelStore, messageStore) {
    if (!messageStore) {
        throw new Error('no channel and/or message store found');
    }

    // Allow message creator to modify this message.
    app.patch('/v1/messages/:messageID', (req, res, next) => {
        const userJSON = req.get('X-User');
        const user = JSON.parse(userJSON);
        const messageID = new mongodb.ObjectID(req.params.messageID);
        var channelID;
        var channel;
        messageStore
            .get(messageID)
            .then(message => {
                if (message.creator.id !== user.id) {
                    res.set('Content-Type', 'text/plain');
                    res
                        .status(403)
                        .send('only message creator can modify this message');
                    return;
                }
                channelID = message.channelID;
            })
            .then(() => {
                channel = channelStore.get(channelID)
            })
            .then(() => {
                const updates = {
                    body: req.body.body,
                    editedAt: Date.now()
                };
                return messageStore.update(messageID, updates);
            })
            .then(updatedMessage => {
                res.json(updatedMessage);
                let data = {
                    type: 'message-update',
                    message: updatedMessage,
                    userIDs: channel.members
                };
                sendToMQ(req, data);

            })
            .catch(err => {
                console.log(err);
            });

    });

    // Allow message creator to delete this message.
    app.delete('/v1/messages/:messageID', (req, res, next) => {
        const userJSON = req.get('X-User');
        const user = JSON.parse(userJSON);
        const messageID = new mongodb.ObjectID(req.params.messageID);
        var channelID;
        var channel;


        messageStore
            .get(messageID)
            .then(message => {
                if (message.creator.id !== user.id) {
                    res.set('Content-Type', 'text/plain');
                    res
                        .status(403)
                        .send('only message creator can delete this message');
                    return;
                }
                channelID = message.channelID;
            })
            .then(() =>{
                channel= channelStore.get(channelID);
            })
            .then(() => {
                return messageStore.delete(messageID);
            })
            .then(deletedMessage => {
                res.set('Content-Type', 'text/plain');
                res
                    .status(200)
                    .send('message deleted');
                let data = {
                    type: 'message-delete',
                    messageID: messageID,
                    userIDs: channel.members
                };
                sendToMQ(req, data);
            })
            .catch(err => {
                console.log(err);
            });
    });

}

module.exports = MessageHandler;