const express = require("express");

function ChannelHandler(app, channelStore, messageStore) {

    app.get("/v1/channels", (req, res, next) => {
        channelStore
            .getAll()
            .then(channels => {
                res.json(channels);
            })
            .catch(err => {
                console.log(err);
            });
    });

    app.post("/v1/channels", (req, res, next) => {
        const name = req.body.name;
        if (!name) {
            res.set('Content-Type', 'text/plain');
            res
                .status(400)
                .send('please attach channel name in the request');
            return;
        }

        let description = '';
        if (req.body.description) {
            description = req.body.description;
        }

        const userJSON = req.get('X-User');
        const user = JSON.parse(userJSON);
        const members = [];
        if(req.body.members){
            members = req.body.members;
        }

        members.push(user.id);
        const channel = new Channel(name, description, true, members, user);
        channelStore
            .insert(channel)
            .then(channel => {
                res.json(channel);
            })
            .catch(err => {
                console.log(err);
            });

    });

    app.get("/v1/channels/:channelID", (req, res, next) => {
        const id = req.params.channelID;
        const channelID = new mongodb.ObjectID(id);

        channelStore
            .get(channelID)
            .then(channel => {
                if (channel.privateChannel && channel.members) {
                    let userJSON = req.get('X-User');
                    let user = JSON.parse(userJSON);
                    if (members.indexOf(user.id) == -1) {
                        res.set('Content-Type', 'text/plain');
                        res
                            .status(403)
                            .send('current user is not a member in this channel');
                        return;
                    }
                }
            })
            .catch(err => {
                console.log(err);
            });

        messageStore
            .getAll(channelID)
            .then(messages => {
                res.json(messages);
            })
            .catch(err => {
                console.log(err);
            });
    });

    app.post("/v1/channels/:channelID", (req, res, next) => {
        const id = req.params.channelID;
        const channelID = new mongodb.ObjectID(id);
        const userJSON = req.get('X-User');
        const user = JSON.parse(userJSON);
        const messageBody = req.body.body;
        const message = new Message(channelID, messageBody, user);

        channelStore
            .get(channelID)
            .then(channel => {
                if (channel.privateChannel && channel.members) {
                    if (members.indexOf(user.id) == -1) {
                        res.set('Content-Type', 'text/plain');
                        res
                            .status(403)
                            .send('current user is not a member in this channel');
                        return;
                    }
                }
            })
            .catch(err => {
                console.log(err);
            });

        messageStore
            .insert(message)
            .then(message => {
                res.status(201);
                res.json(message);
            })
            .catch(err => {
                console.log(err);
            });
    });

    app.patch("/v1/channels/:channelID", (req, res, next) => {
        const id = req.params.channelID;
        const userJSON = req.get('X-User');
        const user = JSON.parse(userJSON);
        const channelID = new mongodb.ObjectID(id);
        channelStore
            .get(channelID)
            .then(channel => {
                if (!channel) {
                    res.set('Content-Type', 'text/plain');
                    res
                        .status(400)
                        .send('no such channel found');
                    return;
                }
                // If the current user isn't the creator, respond with the status code 403
                // (Forbidden).
                if (!channel.creator || channel.creator.id !== user.id) {
                    res.set('Content-Type', 'text/plain');
                    res
                        .status(403)
                        .send('current user is not the creator of this channel');
                    return;
                }
            })
            .catch(err => {
                console.log(err);
            });

        const updates = {};
        if (req.body.name) {
            updates.name = req.body.name;
        }
        if (req.body.description) {
            updates.description = req.body.description;
        }
        updates.editedAt = Date.now();
        channelStore
            .update(channelID, updates)
            .then(updatedChannel => {
                res.json(updatedChannel);
            })
            .catch(err => {
                console.log(err);
            });

    });

    app.delete("/v1/channels/:channelID", (req, res, next) => {
        const id = req.params.channelID;
        const userJSON = req.get('X-User');
        const user = JSON.parse(userJSON);
        const channelID = new mongodb.ObjectID(id);

        channelStore
            .get(channelID)
            .then(channel => {
                if (!channel) {
                    res.set('Content-Type', 'text/plain');
                    res
                        .status(400)
                        .send('no such channel found');
                    return;
                }
                // If the current user isn't the creator, respond with the status code 403
                // (Forbidden).
                if (!channel.creator || channel.creator.id !== user.id) {
                    res.set('Content-Type', 'text/plain');
                    res
                        .status(403)
                        .send('current user is not the creator of this channel');
                    return;
                }
            })
            .catch(err => {
                console.log(err);
            });

        messageStore.deleteAll(channelID);
        channelStore.delete(channelID);
        res.set('Content-Type', 'text/plain');
        res
            .status(200)
            .send('channel deleted');
    });

    app.post("/v1/channels/:channelID/members", (req, res, next) => {
        const id = req.params.channelID;
        const channelID = new mongodb.ObjectID(id);
        const userJSON = req.get('X-User');
        const user = JSON.parse(userJSON);
        const body = req.body;
        let members = null;

        channelStore
            .get(channelID)
            .then(channel => {
                if (!channel) {
                    res.set('Content-Type', 'text/plain');
                    res
                        .status(400)
                        .send('no such channel found');
                    return;
                }
                // If the current user isn't the creator, respond with the status code 403
                // (Forbidden).
                if (!channel.creator || channel.creator.id !== user.id) {
                    res.set('Content-Type', 'text/plain');
                    res
                        .status(403)
                        .send('current user is not the creator of this channel');
                    return;
                }
                members = channel.members;
            })
            .catch(err => {
                console.log(err);
            });

        const updates = {};

        members.push(body.id);
        if (members) {
            updates.members = members;
        }

        updates.editedAt = Date.now();
        channelStore
            .update(channelID, updates)
            .then(updatedChannel => {
                res.set('Content-Type', 'text/plain');
                res
                    .status(201)
                    .send('member has been added to the channel');
            })
            .catch(err => {
                console.log(err);
            });

    });

    app.delete("/v1/channels/:channelID/members", (req, res, next) => {
        const id = req.params.channelID;
        const channelID = new mongodb.ObjectID(id);
        const userJSON = req.get('X-User');
        const user = JSON.parse(userJSON);
        const body = req.body;
        let members = null;

        channelStore
            .get(channelID)
            .then(channel => {
                if (!channel) {
                    res.set('Content-Type', 'text/plain');
                    res
                        .status(400)
                        .send('no such channel found');
                    return;
                }
                // If the current user isn't the creator, respond with the status code 403
                // (Forbidden).
                if (!channel.creator || channel.creator.id !== user.id) {
                    res.set('Content-Type', 'text/plain');
                    res
                        .status(403)
                        .send('current user is not the creator of this channel');
                    return;
                }
                members = channel.members;
            })
            .catch(err => {
                console.log(err);
            });

        const updates = {};

        let index = members.indexOf(body.id);
        if (index > -1) {
            members.splice(index, 1);
        }

        if (members) {
            updates.members = members;
        }

        updates.editedAt = Date.now();
        channelStore
            .update(channelID, updates)
            .then(updatedChannel => {
                res.set('Content-Type', 'text/plain');
                res
                    .status(200)
                    .send('member has been removed to the channel');
            })
            .catch(err => {
                console.log(err);
            });
    });

    // // error handler that will be called if any handler earlier in the chain throws
    // // an exception or passes an error to next()
    // app.use((err, req, res, next) => {
    //     //write a stack trace to standard out, which writes to the server's log
    //     console.error(err.stack)

    //     //but only report the error message to the client, with a 500 status code
    //     res.set("Content-Type", "text/plain");
    //     res
    //         .status(500)
    //         .send(err.message);
    // });

};

module.exports = ChannelHandler;