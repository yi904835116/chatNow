'use strict';

const mongodb = require('mongodb');

class MessageStore {
    constructor(db, colName) {
        this.collection = db.collection(colName);
    }

    // insert() creates a new message in MongoDB.
    insert(message) {
        message._id = new mongodb.ObjectID();
        return this.collection.insertOne(message).then(() => message);
    }

    // get() retrieves one message from MongoDB for a given message ID.
    get(id) {
        return this.collection.findOne({ _id: id });
    }

    // getAll() retrieves up to 100 messages from MongoDB for a given channel ID.
    getAll(channelID) {
        return this.collection
            .find({ channelID: channelID })
            .limit(100)
            .toArray();
    }

    // update() updates a message for a given message ID.
    // It returns the updated message.
    update(id, updates) {
        let updateDoc = {
            $set: updates
        };
        return this.collection
            .findOneAndUpdate({ _id: id }, updateDoc, { returnOriginal: false })
            .then(result => {
                return result.value;
            });
    }

    // delete() deletes a message for a given message ID.
    delete(id) {
        return this.collection.deleteOne({ _id: id });
    }

    // deleteAll() deletes all messages for a given channel ID.
    deleteAll(channelID) {
        return this.collection.deleteMany({ channelID: channelID });
    }
}

module.exports = MessageStore;