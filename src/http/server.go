package http

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
)

const (
	UsualHTMLTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8" />
    <title>Список элементов</title>
    <link rel="stylesheet" href="/static/usual.css" />
</head>
<body>
    <h1>Список элементов</h1>
    <ul>
        {{range .}}
        <li>ID: {{.ID}}, Название: {{.Name}}</li>
        {{else}}
        <li>Элементы отсутствуют</li>
        {{end}}
    </ul>
</body>
</html>`

	JSONHTMLTemplate = `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8" />
    <title>JSON данные</title>
    <link rel="stylesheet" href="/static/json.css" />
</head>
<body>
    <h1>JSON данные</h1>
    <pre>{{ . }}</pre>
</body>
</html>`
)

type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var Items = []Item{
	{ID: 1, Name: "Элемент 1"},
	{ID: 2, Name: "Элемент 2"},
}

// API обработчик отдаёт JSON
func ItemsAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Items)
}

// Шаблон для /items с выводом списка
var ItemsTemplate *template.Template

func ItemsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := ItemsTemplate
	if tmpl == nil {
		var err error
		// Попытка загрузить из файла
		tmpl, err = template.ParseFiles("static/usual.html")
		if err != nil {
			// Fallback на встроенный шаблон для тестирования
			tmpl, err = template.New("usual").Parse(UsualHTMLTemplate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	if err := tmpl.Execute(w, Items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Шаблон для /itemsjson где выводится JSON в формате {{ . }}
var ItemsJSONTemplate *template.Template

// Обработчик /itemsjson — отдаёт HTML с JSON в шаблоне
func ItemsJSONHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := ItemsJSONTemplate
	if tmpl == nil {
		var err error
		// Попытка загрузить из файла
		tmpl, err = template.ParseFiles("static/json.html")
		if err != nil {
			// Fallback на встроенный шаблон для тестирования
			tmpl, err = template.New("json").Parse(JSONHTMLTemplate)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	jsonData, err := json.MarshalIndent(Items, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, string(jsonData)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Start() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.Handle("/", fs)
	http.HandleFunc("/api/items", ItemsAPIHandler)
	http.HandleFunc("/items", ItemsHandler)
	http.HandleFunc("/itemsjson", ItemsJSONHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
