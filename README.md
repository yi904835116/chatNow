A online chat web application that allows users to sign in/sign up, update profile, search users, create channels, and chat in real-time



## Core Technologies

### API Server 

* Go
* Node.js
* Docker
* MySQL
* MongoDB
* Reddis
* RabbitMQ
* Websocket
* Concurrent programming with Mutex and channels

## Features
* Implemented user login with HTTPS Auth, Middleware, credential/session encryption with bcrpt and hmac
* Modelled and trained a NPL chatbot using Wit.ai and created a Node.js microservice that handles user’s
chat-related questions
* Enabled Websocket for real-time notification with Concurrency using goroutine, go channel, and Mutex Lock
* Deployed the go and Node servers and web client using Docker and Digital Ocean
* Suggest-as-you-type searching (Trie implementation)
* Dynamic service discovery
* Allow users to react to a message using an emoji.
* Detect weak passwords during sign-up
* Automatically generates page summary if a message contains links

## Software Architecture

![Software Architecture](https://raw.githubusercontent.com/zicodeng/tahc-z/master/software-architecture.png "Software Architecture")
