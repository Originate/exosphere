FROM node

# These steps ensure that yarn is only run when package.json changes
COPY ./package.json .
RUN yarn --production=true
COPY . .
