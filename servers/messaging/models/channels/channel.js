'use strict';

class Channel {
    constructor(name, description, private, members, creator) {

        this.name = name;

        this.description = description;

        this.private = private;

        this.members = members;

        this.creator = creator;

        this.createdAt = Date.now();

        this.editedAt = Date.now();

        // Note: channel object has another property id, which will be created
        // when we insert it to MongoDB.
    }
}

module.exports = Channel;