FROM nodesource/node

# These steps ensure that npm install is only run when package.json changes
COPY ./package.json .
ENV NODE_ENV=""
RUN npm install
COPY . .
