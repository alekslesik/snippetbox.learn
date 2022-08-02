package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Помощник serverError записывает сообщение об ошибке в errorLog и
// затем отправляет пользователю ответ 500 "Внутренняя ошибка сервера".
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Помощник clientError отправляет определенный код состояния и соответствующее описание
// пользователю. Мы будем использовать это в следующий уроках, чтобы отправлять ответы вроде 400 "Bad
// Request", когда есть проблема с пользовательским запросом.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// Мы также реализуем помощник notFound. Это просто
// удобная оболочка вокруг clientError, которая отправляет пользователю ответ "404 Страница не найдена".
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	// extract pattern depending "name"
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("Pattern %s not exist!", name))
		return
	}

	// initialize a new buffer
	buf := new(bytes.Buffer)

	// write template to the buffer, instead straight to http.ResponseWriter
	err := ts.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// write buffer to http.ResponseWriter
	buf.WriteTo(w)

	// rendering pattern files passing dynamic data from td variable
	// err = ts.Execute(w, td)
	// if err != nil {
	// 	app.serverError(w, err)
	// }

}
