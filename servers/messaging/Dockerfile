FROM node
# set the current working directory to /messaging
WORKDIR /messaging
# copy the package.json and package-lock.json files to the work directory
COPY package.json package-lock.json ./
# run npm install to install all dependent packages
RUN npm install
# copy your JavaScript source files
COPY . .
# declare that your server will listen on port 80
EXPOSE 80
# set the entrypoint command
ENTRYPOINT ["node", "index.js"]