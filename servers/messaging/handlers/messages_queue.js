// sendToMQ is a message publisher 
// Send the message to RabbitMQ channel.
// Gateway will listen to this channel,
// and then distributes the message to all connected clients.

function sendToMQ(req, message) {
    var mqChannel = req.app.get('mqChannel');
    var qName = req.app.get('qName');
    var stringifiedData =  JSON.stringify(message)
    mqChannel.sendToQueue(qName, Buffer.from(stringifiedData));
};

module.exports = sendToMQ ;