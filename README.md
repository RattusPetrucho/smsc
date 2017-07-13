# smsc
Пакет smsc предназначен для отправки сообщений через https://smsc.ru

##Instalation

```
go get github.com/RattusPetrucho/bashorg_reader
```
##Documentation

import "github.com/RattusPetrucho/smsc"

CONSTANTS

const (
    VoiceMale    = "m"
    Voice2Male   = "m2"
    VoiceFemail  = "w"
    Voice2Femail = "w2"
    Voice3Femail = "w2"
    Voice4Femail = "w2"
)
    Возможные типы голоса для озвучивания текста в голосовых сообщениях.
    Используются как параметры для функции SetVoice

TYPES

type Client struct {
    // contains filtered or unexported fields
}
    Клиент для работы с smsc. Является потокобезопасным.

func New(login, password string) (*Client, error)
    Создание объекта клиента. Принимает login и пароль/md5-hash пароля от
    аккаунта smsc.ru

func (sc *Client) DisableTinyUrl()
    Отключение автоматического сокращения ссылок в сообщении

func (sc *Client) EnableTinyUrl()
    Включение автоматического сокращения ссылок в сообщении

func (sc *Client) GetEmailCost(subject, message string, emails ...string) (*Response, error)
    Получение стоимости email рассылки без реального отправления

func (sc *Client) GetSmsCost(message string, phones ...string) (*Response, error)
    Получение стоимости sms рассылки без реального отправления

func (sc *Client) GetVoiceCost(message string, phones ...string) (*Response, error)
    Получение стоимости голосовой рассылки без реального отправления

func (sc *Client) SendEmail(subject, message string, emails ...string) (*Response, error)
    Отправка email сообщения subkect - Тема сообщения, обязательный
    параметр. message - Тело сообщения emails - список email адресов

func (sc *Client) SendSms(id, message string, phones ...string) (*Response, error)
    Отправка текстового sms сообщения. id - Идентификатор сообщения.
    Назначается Клиентом. Служит для дальнейшей идентификации сообщения.
    Если не указывать, то будет назначен автоматически. Не обязательно
    уникален. Идентификатор представляет собой 32-битное число в диапазоне
    от 1 до 2147483647, либо строку длиной до 40 символов, состоящую из
    латинских букв, цифр и символов ".-_". В случае 2-х одинаковых
    идентификаторов по запросу статуса будет возвращен статус последнего
    сообщения. message - тело сообщения. phones - Cписок номеров мобильных
    телефонов в международном формате, на которые отправляется сообщение

func (sc *Client) SendVoice(id, message string, phones ...string) (*Response, error)
    Отправка голосового сообщения id - Идентификатор сообщения. Назначается
    Клиентом. Служит для дальнейшей идентификации сообщения. Если не
    указывать, то будет назначен автоматически. Не обязательно уникален.
    Идентификатор представляет собой 32-битное число в диапазоне от 1 до
    2147483647, либо строку длиной до 40 символов, состоящую из латинских
    букв, цифр и символов ".-_". В случае 2-х одинаковых идентификаторов по
    запросу статуса будет возвращен статус последнего сообщения. message -
    тело сообщения. phones - Cписок номеров мобильных телефонов в
    международном формате, на которые отправляется сообщение

func (sc *Client) SetCharset(charset string) error
    Задаёт кодировку. По умолчанию utf-8. Принимает значения: utf-8,
    windows-1251, koi8-r.

func (sc *Client) SetSenderEmail(email string) error
    Задаёт email отправителя.

func (sc *Client) SetSenderName(name string) error
    Задаёт имя отправителя, которое видит абонент при получении sms.

func (sc *Client) SetVoice(v string) error
    Выбор голоса используемого для озвучивания текста. Для голосовых
    сообщений.

type Response struct {
    Count int     // Кол-во отправленных сообщений
    Cost  float64 // Стоимость рассылки
}
    Ответ от smsc


