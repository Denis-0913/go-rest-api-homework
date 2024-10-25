package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Ниже напишите обработчики для каждого эндпоинта
// Обработчик для получения всех задач
// Конечная точка /tasks.
// Метод GET.
// При успешном запросе сервер должен вернуть статус 200 OK.
// При ошибке сервер должен вернуть статус 500 Internal Server Error.
func getTasks(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса tasks
	resp, err := json.Marshal(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

// Обработчик для отправки задачи на сервер
// Обработчик должен принимать задачу в теле запроса и сохранять ее в мапе.
// Конечная точка /tasks.
// Метод POST.
// При успешном запросе сервер должен вернуть статус 201 Created.
// При ошибке сервер должен вернуть статус 400 Bad Request.
func postTasks(w http.ResponseWriter, r *http.Request) {

	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

}

// Обработчик для получения задачи по ID
// Обработчик должен вернуть задачу с указанным в запросе пути ID, если такая есть в мапе.
// В мапе ключами являются ID задач. Вспомните, как проверить, есть ли ключ в мапе. Если такого ID нет, верните соответствующий статус.
// Конечная точка /tasks/{id}.
// Метод GET.
// При успешном выполнении запроса сервер должен вернуть статус 200 OK.
// В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
func getTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Обработчик удаления задачи по ID
// Обработчик должен удалить задачу из мапы по её ID. Здесь так же нужно сначала проверить, есть ли задача с таким ID в мапе, если нет вернуть соответствующий статус.
// Конечная точка /tasks/{id}.
// Метод DELETE.
// При успешном выполнении запроса сервер должен вернуть статус 200 OK.
// В случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
func deleteTask(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	//В случае отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request.
	_, ok := tasks[id]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	// регистрируем в роутере эндпоинт `/tasks` с методом GET, для которого используется обработчик `getTasks`
	r.Get("/tasks", getTasks)
	// регистрируем в роутере эндпоинт `/tasks` с методом POST, для которого используется обработчик `postTasks`
	r.Post("/tasks", postTasks)
	// регистрируем в роутере эндпоинт `artist/{id}` с методом GET, для которого используется обработчик `getTask`
	r.Get("/tasks/{id}", getTask)
	// регистрируем в роутере эндпоинт `artist/{id}` с методом DELETE, для которого используется обработчик `deleteTasks`
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
