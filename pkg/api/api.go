package api

import "net/http"

func Init(mux *http.ServeMux) {
	mux.Handle("/", http.FileServer(http.Dir("./web/")))

	mux.HandleFunc("/api/signin", authenticationHandler)

	mux.HandleFunc("/api/nextdate", nextDayHandler)

	mux.HandleFunc("/api/task/done", auth(clearTaskHandler))
	mux.HandleFunc("/api/task", auth(taskHandler))

	mux.HandleFunc("/api/tasks", auth(tasksListHandler))
}
