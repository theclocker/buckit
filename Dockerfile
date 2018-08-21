FROM golang:1.10

RUN ["apt-get", "update"]
RUN ["apt-get", "install", "-y", "zsh", "nano", "fonts-powerline", "fontconfig", "locales"]
RUN wget https://github.com/robbyrussell/oh-my-zsh/raw/master/tools/install.sh -O - | zsh || true
RUN sed -i 's/ZSH_THEME="robbyrussell"/ZSH_THEME="bira"/g' ~/.zshrc

WORKDIR /go/src/github.com/theclocker/buckit
ADD . /go/src/github.com/theclocker/buckit

RUN go get -d -v ./...
RUN go install -v ./...

RUN go get github.com/codegangsta/gin