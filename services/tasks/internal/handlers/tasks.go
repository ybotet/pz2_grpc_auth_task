package handlers

import (
	"encoding/json"
	"net/http"
)

type Task struct {
    ID      string `json:"id"`
    Title   string `json:"title"`
    UserID  string `json:"user_id"`
}

type TaskHandler struct {
    tasks []Task
}

func NewTaskHandler() *TaskHandler {
    return &TaskHandler{
        tasks: []Task{
            {ID: "1", Title: "Tarea 1", UserID: "user1"},
            {ID: "2", Title: "Tarea 2", UserID: "user2"},
        },
    }
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
    user := r.Context().Value("user").(string)
    
    // Filtrar tareas por usuario
    var userTasks []Task
    for _, task := range h.tasks {
        if task.UserID == user {
            userTasks = append(userTasks, task)
        }
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(userTasks)
}