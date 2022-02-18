## This chatbot download audio messages from whatsapp chat, translate it and reply text message


vktoken - token from https://mcs.mail.ru/app/mcs2676925534/services/machinelearning/voice/access/

yatoken - token from https://cloud.yandex.ru/docs/iam/concepts/authorization/api-key

* enable in your phone dev feature multi  login 
* first run app and render qr code `echo 2@... | qrencode -t ansiutf8` in a terminal  and scan it 

### example run 
```
export vktoken="<yor_mcs_token>"
export VOICE_MODEL="MAILRU" 
```

or 

```
export yatoken="<your_yandex_token>"
export VOICE_MODEL="YANDEX"
```

```
go mod tidy
go build
./wh-audiobot
```
you can build image 
```
docker build -t wh-audio . 
docker run -it  -v ./examplestore.db:/app/examplestore.db  -e vktoken=<yourtoken> wh-audio  
```
