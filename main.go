package main

import (
	"log"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/liuyh73/dailyhub.service/service"
)

const defaultPort = "9090"

func NewServer() *negroni.Negroni {
	router := mux.NewRouter()
	initRoutes(router)

	n := negroni.Classic()
	n.UseHandler(router)
	return n
}

func initRoutes(router *mux.Router) {
	router.Use(service.JWTMiddleware)
	router.HandleFunc("/api", service.ApiHandler).Methods("GET")
	// 注册、登录、退出
	router.HandleFunc("/api/login", service.LoginHandler).Methods("POST")
	router.HandleFunc("/api/register", service.RegisterHandler).Methods("POST")
	router.HandleFunc("/api/logout", service.LogoutHandler).Methods("POST", "GET")
	// user相关
	router.HandleFunc("/api/users/{user}", service.GetProfileHandler).Methods("GET")

	// habits相关
	// GET
	sub := router.PathPrefix("/api/habits").Subrouter()
	sub.HandleFunc("", service.GetHabitsHandler).Methods("GET")
	sub.HandleFunc("/{habitId:[0-9]+}", service.GetHabitHandler).Methods("GET")
	sub.HandleFunc("/{habitId:[0-9]+}/{monthId:[0-9-]+}", service.GetMonthHandler).Methods("GET")
	sub.HandleFunc("/{habitId:[0-9]+}/{monthId:[0-9-]+}/{dayId:[0-9]+}", service.GetDayHandler).Methods("GET")

	// POST
	sub.HandleFunc("", service.PostHabitsHandler).Methods("POST")
	sub.HandleFunc("/{habitId:[0-9]+}/{monthId:[0-9-]+}/{dayId:[0-9]+}", service.PostDayHandler).Methods("POST")

	// PUT 修改habits信息，修改打卡日志信息
	sub.HandleFunc("/{habitId:[0-9]+}", service.PutHabitHandler).Methods("PUT")
	sub.HandleFunc("/{habitId:[0-9]+}/{monthId:[0-9-]+}/{dayId:[0-9]+}", service.PutDayHandler).Methods("PUT")

	// DELETE
	sub.HandleFunc("/{habitId:[0-9]+}", service.DeleteHabitHandler).Methods("DELETE")
	sub.HandleFunc("/{habitId:[0-9]+}/{monthId:[0-9-]+}/{dayId:[0-9]+}", service.DeleteDayHandler).Methods("DELETE")

	// dailyCommits相关
	sub = router.PathPrefix("/api/dailycommits").Subrouter()
	// GET
	sub.HandleFunc("", service.GetDailyCommitsHandler).Methods("GET")
	// POST
	sub.HandleFunc("", service.PostDailyCommitHandler).Methods("POST")
	// PUT
	sub.HandleFunc("/{dailyCommitId:[0-9]+}", service.PutDailyCommitHandler).Methods("PUT")
	// DELETE
	sub.HandleFunc("/{dailyCommitId:[0-9]+}", service.DeleteDailyCommitHandler).Methods("DELETE")
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	server := NewServer()
	server.Run(":" + port)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
}
