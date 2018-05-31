'use strict';

class Message {
    constructor(channelID, body, creator) {

        this.channelID = channelID;

        this.body = body;

        this.createdAt = Date.now();

        this.creator = creator;

        this.editedAt = Date.now();


        // Note: message object has another property id, which will be created
        // when we insert it to MongoDB.
    }
}

module.exports = Message;