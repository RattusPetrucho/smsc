# smsc
Пакет smsc предназначен для отправки сообщений через https://smsc.ru.

##Instalation

```
go get github.com/RattusPetrucho/bashorg_reader
```
##Documentation

import "github.com/RattusPetrucho/smsc"


Создание объекта клиента. Принимает login и пароль/md5-hash пароля от аккаунта smsc.ru. Клиент безопасен для одновременного использования несколькими goroutines.
```
sc, err := smsc.New("login", "password")
if err != nil {
    log.Fatal(err)
}
```

Задаёт имя отправителя, которое видит абонент при получении sms.
```
err = sc.SetSenderName("sender")
if err != nil {
    log.Fatal(err)
}
```

Выбор голоса используемого для озвучивания текста. Для голосовых сообщений.
```
err = sc.SetVoice(smsc.VoiceFemail)
if err != nil {
    log.Fatal(err)
}
```

Задаёт email отправителя.
```
err = sc.SetSenderEmail("example@example.com")
if err != nil {
    log.Fatal(err)
}
```

Задаёт кодировку. По умолчанию utf-8. Принимает значения: utf-8, windows-1251, koi8-r.
```
err = sc.SetCharset("utf-8")
if err != nil {
    log.Fatal(err)
}
```

Включение автоматического сокращения ссылок в сообщении
```
sc.EnableTinyUrl()
```

Отключение автоматического сокращения ссылок в сообщении
```
sc.DisableTinyUrl()
```

Отправка текстового sms сообщения. id - Идентификатор сообщения. Назначается Клиентом. Служит для дальнейшей идентификации сообщения. Если не указывать, то будет назначен автоматически. Не обязательно уникален. Идентификатор представляет собой 32-битное число в диапазоне от 1 до 2147483647, либо строку длиной до 40 символов, состоящую из латинских букв, цифр и символов ".-_". В случае 2-х одинаковых идентификаторов по запросу статуса будет возвращен статус последнего сообщения. message - тело сообщения. phones - Cписок номеров мобильных телефонов в международном формате, на которые отправляется сообщение
```
resp, err := sc.SendSms("", "hello world!", "1113311567")
if err != nil {
    log.Fatal(err)
}
```

Отправка голосового сообщения id - Идентификатор сообщения. Назначается Клиентом. Служит для дальнейшей идентификации сообщения. Если не указывать, то будет назначен автоматически. Не обязательно уникален. Идентификатор представляет собой 32-битное число в диапазоне от 1 до 2147483647, либо строку длиной до 40 символов, состоящую из латинских букв, цифр и символов ".-_". В случае 2-х одинаковых идентификаторов по запросу статуса будет возвращен статус последнего сообщения. message - тело сообщения. phones - Cписок номеров мобильных телефонов в международном формате, на которые отправляется сообщение
```
resp, err := sc.SendVoice("", "hello world!", "1113311567")
if err != nil {
    log.Fatal(err)
}
```

Отправка email сообщения subject - Тема сообщения, обязательный параметр. message - Тело сообщения. emails - список email адресов
```
resp, err := sc.SendEmail("subject", "Hello world!", "example1@example.com", "example2@example.com")
if err != nil {
    log.Fatal(err)
}
```

Получение стоимости sms рассылки без реального отправления
```
resp, err := sc.GetSmsCost("Hell world!", "1113311567")
if err != nil {
    log.Fatal(err)
}
```

Получение стоимости голосовой рассылки без реального отправления
```
resp, err := sc.GetVoiceCost("Hell world!", "1113311567")
if err != nil {
    log.Fatal(err)
}
```

Получение стоимости email рассылки без реального отправления
```
resp, err := sc.GetEmailCost("subject", "Hello world!", "example1@example.com", "example2@example.com")
if err != nil {
    log.Fatal(err)
}
```

Ответ от smsc
```
type Response struct {
    Count int     // Кол-во отправленных сообщений
    Cost  float64 // Стоимость рассылки
}
```

##CONSTANTS

Возможные типы голоса для озвучивания текста в голосовых сообщениях.
Используются как параметры для функции SetVoice
* VoiceMale
* Voice2Male
* VoiceFemail
* Voice2Femail
* Voice3Femail
* Voice4Femail