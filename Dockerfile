FROM smm-scrapper:latest

WORKDIR /app

COPY . .


RUN apt-get update && apt-get install -y wget gnupg2
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add -
RUN echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' > /etc/apt/sources.list.d/google-chrome.list
RUN apt-get update && apt-get install -y google-chrome-stable

EXPOSE 8080

CMD ["go", "run", "./cmd/default/main.go"]