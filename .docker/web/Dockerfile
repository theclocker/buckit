FROM node:17-stretch-slim
RUN apt-get update || : && apt-get install python3 python3-pip -y
WORKDIR /web

ADD ./web/* ./
RUN npm install