package api

import (
	"fmt"
	"net/http"

	"github.com/LJlamar/go_final_project/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	searchStr := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, searchStr) // maximum value of given in task's advice: "show from 10 to 50 entries"
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Println(err)
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	fmt.Println(tasks)
	writeJson(w, TasksResp{
		Tasks: tasks,
	})
}
