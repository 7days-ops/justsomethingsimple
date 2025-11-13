package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"text/template"

	myhttp "webserver/http"

	"github.com/stretchr/testify/assert"
)

// Разбор обычного API (json)
func TestItemsAPIHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/items", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()

	myhttp.ItemsAPIHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var data []myhttp.Item
	err = json.NewDecoder(resp.Body).Decode(&data)
	assert.NoError(t, err)
	assert.Len(t, data, 2)
}

// Разбор plain HTML-страницы
func TestItemsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/items", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()

	myhttp.ItemsHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyStr := string(body)

	// Проверяем наличие конкретных элементов/текста в HTML-странице
	assert.Contains(t, bodyStr, "Элемент 1")
	assert.Contains(t, bodyStr, "Элемент 2")
	// Можно проверить наличие тегов или классов, если хотите
	assert.Contains(t, bodyStr, "<ul>")
}

// Разбор JSON, вложенного внутрь HTML (например, <pre>{{ . }}</pre>)
func TestItemsJSONHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/itemsjson", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()

	myhttp.ItemsJSONHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	bodyStr := string(body)

	// Находим JSON внутри <pre>...</pre> (или другого тега)
	start := strings.Index(bodyStr, "<pre>")
	end := strings.Index(bodyStr, "</pre>")
	assert.Greater(t, end, start)
	jsonStr := bodyStr[start+len("<pre>") : end]

	var data []myhttp.Item
	err = json.Unmarshal([]byte(jsonStr), &data)
	assert.NoError(t, err)
	assert.Len(t, data, 2)
}

func TestItemsHandlerError(t *testing.T) {
	origTemplate := myhttp.ItemsTemplate
	// Создаем невалидный шаблон для тестирования ошибки
	badTemplate, _ := template.New("bad").Parse("{{range .}}{{.NonExistentField}}{{end}}")
	myhttp.ItemsTemplate = badTemplate

	req, err := http.NewRequest("GET", "/items", nil)
	assert.NoError(t, err)
	w := httptest.NewRecorder()

	myhttp.ItemsHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	myhttp.ItemsTemplate = origTemplate // Восстанавливаем шаблон
}
