'use strict';

const mongodb = require('mongodb');

class ChannelStore {
    constructor(db, colName) {
        this.collection = db.collection(colName);
    }

    // insert() creates a new channel in MongoDB.
    insert(channel) {
        channel.id = new mongodb.ObjectID();
        return this.collection.insertOne(channel).then(() => channel);
    }

    // get() retrieves one channel from MongoDB for a given channel ID.
    get(id) {
        return this.collection.findOne({ id: id });
    }

    // getByName retrieves one channel from MongoDB for a given channel name.
    getByName(name) {
        return this.collection.findOne({ name: name });
    }

    // getAll() retrieves all channel from MongoDB.
    getAll() {
        return this.collection
            .find({})
            .limit(100)
            .toArray();
    }

    // update() updates a channel for a given channel ID.
    // It returns the updated channel.
    update(id, updates) {
        let updateDoc = {
            $set: updates
        };
        return this.collection
            .findOneAndUpdate({ id: id }, updateDoc, { returnOriginal: false })
            .then(result => {
                return result.value;
            });
    }

    // delete() deletes a channel for a given channel ID.
    delete(id) {
        return this.collection.deleteOne({ id: id });
    }
}

module.exports = ChannelStore;