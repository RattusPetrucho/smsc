package smsc

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type SmscResponse struct {
	Cnt  int
	Cost float64
}

type smscResponse struct {
	Cnt  float64 `json:"cnt"`
	Cost string  `json:"cost"`
	Err  string  `json:"error"`
}

func (sc *SmscClient) send(params *url.Values) (*SmscResponse, error) {
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

	ret := new(SmscResponse)
	ret.Cnt = int(sr.Cnt)
	ret.Cost, err = strconv.ParseFloat(sr.Cost, 64)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Отправка текстового sms сообщения.
// id - Идентификатор сообщения. Назначается Клиентом. Служит для дальнейшей идентификации сообщения. Если не указывать, то будет назначен автоматически. Не обязательно уникален.
// Идентификатор представляет собой 32-битное число в диапазоне от 1 до 2147483647, либо строку длиной до 40 символов, состоящую из латинских букв, цифр и символов ".-_".
// Если не указывать, то будет назначен автоматически. Не обязательно уникален. В случае 2-х одинаковых идентификаторов по запросу статуса будет возвращен статус последнего сообщения.
// mess - тело сообщения.
// phones - Cписок номеров мобильных телефонов в международном формате, на которые отправляется сообщение
func (sc *SmscClient) SendSms(id, mess string, phones ...string) (*SmscResponse, error) {
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
	params.Add("mes", mess)
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
