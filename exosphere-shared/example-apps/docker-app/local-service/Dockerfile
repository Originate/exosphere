FROM nodesource/node

# These steps ensure that npm install is only run when package.json changes
COPY ./package.json .
RUN npm install --production
COPY . .
