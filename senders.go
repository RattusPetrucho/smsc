package smsc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

// Ответ от smsc
type Response struct {
	Count int     // Кол-во отправленных сообщений
	Cost  float64 // Стоимость рассылки
}

type smscResponse struct {
	Cnt  float64 `json:"cnt"`
	Cost string  `json:"cost"`
	Err  string  `json:"error"`
}

func (sc *Client) send(params *url.Values) (*Response, error) {
	req, err := http.NewRequest("POST", "http://smsc.ru/sys/send.php", bytes.NewBufferString(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	sr := new(smscResponse)

	err = json.NewDecoder(resp.Body).Decode(sr)
	if err != nil {
		return nil, err
	}

	if sr.Err != "" {
		return nil, errors.New(sr.Err)
	}

	ret := new(Response)
	ret.Count = int(sr.Cnt)
	ret.Cost, err = strconv.ParseFloat(sr.Cost, 64)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Отправка текстового sms сообщения.
// id - Идентификатор сообщения. Назначается Клиентом. Служит для дальнейшей идентификации сообщения. Если не указывать, то будет назначен автоматически. Не обязательно уникален.
// Идентификатор представляет собой 32-битное число в диапазоне от 1 до 2147483647, либо строку длиной до 40 символов, состоящую из латинских букв, цифр и символов ".-_".
// В случае 2-х одинаковых идентификаторов по запросу статуса будет возвращен статус последнего сообщения.
// message - тело сообщения.
// phones - Cписок номеров мобильных телефонов в международном формате, на которые отправляется сообщение
func (sc *Client) SendSms(id, message string, phones ...string) (*Response, error) {
	sc.mu.RLock()

	if len(phones) == 0 {
		sc.mu.RUnlock()
		return nil, errors.New("did not set phone numbers")
	}

	var phones_str = ""
	var sep = ""
	for _, phone := range phones {
		phones_str += sep + phone
		sep = ","
	}

	params := url.Values{}
	params.Add("login", sc.login)
	params.Add("psw", sc.password)
	params.Add("phones", phones_str)
	params.Add("mes", message)
	params.Add("charset", sc.charset)
	params.Add("cost", "2")
	params.Add("fmt", "3")
	params.Add("tinyurl", sc.tinyurl)
	if sc.sender != "" {
		params.Add("sender", sc.sender)
	}
	if id != "" {
		params.Add("id", id)
	}

	sc.mu.RUnlock()

	return sc.send(&params)
}

// Получение стоимости sms рассылки без реального отправления
func (sc *Client) GetSmsCost(message string, phones ...string) (*Response, error) {
	sc.mu.RLock()

	if len(phones) == 0 {
		sc.mu.RUnlock()
		return nil, errors.New("did not set phone numbers")
	}

	params := url.Values{}
	params.Add("login", sc.login)
	params.Add("psw", sc.password)
	params.Add("charset", sc.charset)
	params.Add("tinyurl", sc.tinyurl)
	if sc.sender != "" {
		params.Add("sender", sc.sender)
	}
	sc.mu.RUnlock()

	var phones_str = ""
	var sep = ""
	for _, phone := range phones {
		phones_str += sep + phone
		sep = ","
	}

	params.Add("phones", phones_str)
	params.Add("mes", message)
	params.Add("cost", "1")
	params.Add("fmt", "3")

	return sc.send(&params)
}

// Отправка голосового сообщения
// id - Идентификатор сообщения. Назначается Клиентом. Служит для дальнейшей идентификации сообщения. Если не указывать, то будет назначен автоматически. Не обязательно уникален.
// Идентификатор представляет собой 32-битное число в диапазоне от 1 до 2147483647, либо строку длиной до 40 символов, состоящую из латинских букв, цифр и символов ".-_".
// В случае 2-х одинаковых идентификаторов по запросу статуса будет возвращен статус последнего сообщения.
// message - тело сообщения.
// phones - Cписок номеров мобильных телефонов в международном формате, на которые отправляется сообщение
func (sc *Client) SendVoice(id, message string, phones ...string) (*Response, error) {
	sc.mu.RLock()

	if len(phones) == 0 {
		sc.mu.RUnlock()
		return nil, errors.New("did not set phone numbers")
	}

	params := url.Values{}
	params.Add("login", sc.login)
	params.Add("psw", sc.password)
	params.Add("voice", sc.voice)
	if sc.sender != "" {
		params.Add("sender", sc.sender)
	}
	sc.mu.RUnlock()

	var phones_str = ""
	var sep = ""
	for _, phone := range phones {
		phones_str += sep + phone
		sep = ","
	}

	params.Add("phones", phones_str)
	params.Add("mes", message)
	params.Add("call", "1")
	params.Add("cost", "2")
	params.Add("fmt", "3")
	if id != "" {
		params.Add("id", id)
	}

	return sc.send(&params)
}

// Получение стоимости голосовой рассылки без реального отправления
func (sc *Client) GetVoiceCost(message string, phones ...string) (*Response, error) {
	sc.mu.RLock()

	if len(phones) == 0 {
		sc.mu.RUnlock()
		return nil, errors.New("did not set phone numbers")
	}

	params := url.Values{}
	params.Add("login", sc.login)
	params.Add("psw", sc.password)
	params.Add("voice", sc.voice)
	if sc.sender != "" {
		params.Add("sender", sc.sender)
	}
	sc.mu.RUnlock()

	var phones_str = ""
	var sep = ""
	for _, phone := range phones {
		phones_str += sep + phone
		sep = ","
	}

	params.Add("phones", phones_str)
	params.Add("mes", message)
	params.Add("call", "1")
	params.Add("cost", "1")
	params.Add("fmt", "3")

	return sc.send(&params)
}

// Отправка email сообщения
// subject - Тема сообщения, обязательный параметр.
// message - Тело сообщения
// emails - список email адресов
func (sc *Client) SendEmail(subject, message string, emails ...string) (*Response, error) {
	sc.mu.RLock()

	if len(emails) == 0 {
		sc.mu.RUnlock()
		return nil, errors.New("did not set emails")
	}
	if subject == "" {
		sc.mu.RUnlock()
		return nil, errors.New("did not set subject")
	}
	if sc.sender_email == "" {
		sc.mu.RUnlock()
		return nil, errors.New("did not set sender email")
	}

	params := url.Values{}
	params.Add("login", sc.login)
	params.Add("psw", sc.password)
	params.Add("sender", sc.sender_email)
	params.Add("charset", sc.charset)
	sc.mu.RUnlock()

	emails_str := ""
	sep := ""
	for _, email := range emails {
		emails_str += sep + email
		sep = ","
	}

	params.Add("phones", emails_str)
	params.Add("mes", message)
	params.Add("subj", subject)
	params.Add("mail", "1")
	params.Add("cost", "2")
	params.Add("fmt", "3")

	return sc.send(&params)
}

// Получение стоимости email рассылки без реального отправления
func (sc *Client) GetEmailCost(subject, message string, emails ...string) (*Response, error) {
	sc.mu.RLock()

	if len(emails) == 0 {
		sc.mu.RUnlock()
		return nil, errors.New("did not set emails")
	}
	if subject == "" {
		sc.mu.RUnlock()
		return nil, errors.New("did not set subject")
	}
	if sc.sender_email == "" {
		sc.mu.RUnlock()
		return nil, errors.New("did not set sender email")
	}

	params := url.Values{}
	params.Add("login", sc.login)
	params.Add("psw", sc.password)
	params.Add("sender", sc.sender_email)
	params.Add("charset", sc.charset)
	sc.mu.RUnlock()

	emails_str := ""
	sep := ""
	for _, email := range emails {
		emails_str += sep + email
		sep = ","
	}

	params.Add("phones", emails_str)
	params.Add("mes", message)
	params.Add("subj", subject)
	params.Add("mail", "1")
	params.Add("cost", "1")
	params.Add("fmt", "3")

	return sc.send(&params)
}
