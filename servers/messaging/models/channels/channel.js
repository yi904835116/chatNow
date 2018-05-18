'use strict';

class Channel {
    constructor(name, description, privateChannel, members, creator) {

        this.name = name;

        this.description = description;

        this.privateChannel = privateChannel;

        this.members = members;

        this.creator = creator;

        this.createdAt = Date.now();

        this.editedAt = Date.now();

        // Note: channel object has another property id, which will be created
        // when we insert it to MongoDB.
    }
}

module.exports = Channel;