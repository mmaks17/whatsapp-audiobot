## This chatbot download audio messages from whatsapp chat, translate it and reply text message


yourtoken - token from https://mcs.mail.ru/app/mcs2676925534/services/machinelearning/voice/access/

* enable in your phone dev feature multi  login 
* first run app and render qr code `echo 2@... | qrencode -t ansiutf8` in a terminal  and scan it 

### example run 


you can build image 
```
docker build -t wh-audio . 
docker run -it  -v ./examplestore.db:/app/examplestore.db  -e vktoken=<yourtoken> wh-audio  
```

any more dev  command 
``
go get -u
go mod tidy
``