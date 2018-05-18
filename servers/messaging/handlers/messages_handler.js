const express = require("express");

function MessageHandler(app, messageStore) {
    if (!messageStore) {
        throw new Error('no channel and/or message store found');
    }

    // Allow message creator to modify this message.
    app.patch('/v1/messages/:messageID', (req, res, next) => {
        const userJSON = req.get('X-User');
        const user = JSON.parse(userJSON);
        const messageID = new mongodb.ObjectID(req.params.messageID);

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
        messageStore
            .get(messageID)
            .then(message => {
                if (message.creator.id !== user.id) {
                    res.set('Content-Type', 'text/plain');
                    res
                        .status(403)
                        .send('only message creator can delete this message');
                }
                return;
            })
            .then(() => {
                return messageStore.delete(messageID);
            })
            .then(deletedMessage => {
                res.set('Content-Type', 'text/plain');
                res
                    .status(200)
                    .send('message deleted');
            })
            .catch(err => {
                console.log(err);
            });
    });

}

module.exports = MessageHandler;