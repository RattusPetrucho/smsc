package smsc

import (
	"errors"
	"net/http"
	"regexp"
	"sync"
)

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Клиент для работы с smsc. Является потокобезопасным.
type SmscClient struct {
	client *http.Client

	mu sync.RWMutex // protect change settings.

	login        string // Login in smsc.ru.
	password     string // Пароль либо md5-hash пароля.
	sender       string // sender - Имя отправителя, отображаемое в телефоне получателя. Разрешены английские буквы, цифры, пробел и некоторые символы. Длина – 11 символов или 15 цифр. Все имена регистрируются в личном кабинете. Для отключения Sender ID по умолчанию необходимо в качестве имени передать пустую строку.
	sender_email string // Email отправителя для email рассылки.
	tinyurl      string // Автоматическое сокращение ссылок в сообщении. 0-оставить. 1-сокращать.
	charset      string // Кодировака. По умолчанию utf-8. Так же может принимать значения: windows-1251, koi8-r.
}

// Создание объекта клиента. Принимает login и пароль/md5-hash пароля от аккаунта smsc.ru
func NewClient(login, password string) (*SmscClient, error) {
	if login == "" || password == "" {
		return nil, errors.New("empty login or password")
	}

	sc := &SmscClient{
		login:    login,
		password: password,
		charset:  "utf-8",
		client:   &http.Client{},
	}

	return sc, nil
}

// Задаёт имя отправителя, которое видит абонент при получении sms.
func (sc *SmscClient) SetSenderName(name string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if name == "" {
		return errors.New("empty name")
	}

	sc.sender = name

	return nil
}

// Задаёт email отправителя.
func (sc *SmscClient) SetSenderEmail(email string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if email == "" {
		return errors.New("empty email")
	}

	if err := validateEmailFormat(email); err != nil {
		return err
	}

	sc.sender_email = email

	return nil
}

// Задаёт кодировку. По умолчанию utf-8. Принимает значения: utf-8, windows-1251, koi8-r.
func (sc *SmscClient) SetCharset(charset string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	switch charset {
	case "utf-8":
		sc.charset = charset
	case "windows-1251":
		sc.charset = charset
	case "koi8-r":
		sc.charset = charset
	default:
		return errors.New("unknown charset")
	}

	return nil
}

// Включение автоматического сокращения ссылок в сообщении
func (sc *SmscClient) EnableTinyUrl() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.tinyurl = "1"
}

// Отключение автоматического сокращения ссылок в сообщении
func (sc *SmscClient) DisableTinyUrl() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.tinyurl = "0"
}

// Валидация формата email
func validateEmailFormat(email string) error {
	if !emailRegexp.MatchString(email) {
		return errors.New("bad email format")
	}

	return nil
}
