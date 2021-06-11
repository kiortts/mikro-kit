/*
В пакете описаны общие для сервиса публичные типы и функции, не предоставляемые в публичном API.
Правила записи данных в http ответы.
*/

package internal

import (
	"encoding/json"
	"net/http"
)

var EmptyJSON = []byte(`{"length": null}`)

// Запись http ответа в виде json с нужными заголовками.
func WriteJSONResponse(w http.ResponseWriter, status int, data []byte) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(status)
	w.Write(data)
}

// Запись структуры в виде json. Обертка для WriteJSONResponse.
func WriteItemAsJSON(w http.ResponseWriter, item interface{}) {

	data, err := json.MarshalIndent(item, "", "  ")
	if err != nil {
		WriteJSONResponse(w, http.StatusInternalServerError, EmptyJSON)
		return
	}

	WriteJSONResponse(w, http.StatusOK, data)
}

// w.Header().Set("Connection", "close")
// w.Header().Set("Access-Control-Allow-Origin", "*")
// w.Header().Set("Access-Control-Allow-Credentials", "true")
// w.Header().Set("Content-Length", strconv.Itoa(len(data)))
// w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
