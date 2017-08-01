FROM node

# These steps ensure that npm install is only run when package.json changes
COPY ./package.json .
RUN curl -o- -L https://yarnpkg.com/install.sh | bash
RUN yarn install --production
COPY . .
