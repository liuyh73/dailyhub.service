package model

var Apis = `{
	"method_get": {
		"api_url": "http://129.204.110.49:9090/api",
		"user_info_url": "http://129.204.110.49:9090/api/users/{username}",
		"habits_url": "http://129.204.110.49:9090/api/habits",
		"habit_url": "http://129.204.110.49:9090/api/habits/{habitId}",
		"month_punch_info_url": "http://129.204.110.49:9090/api/habits/{habitId}/{monthId}",
		"day_punch_info_url": "http://129.204.110.49:9090/api/habits/{habitId}/{monthId}/{dayId}",
		"commits_info_url": "http://129.204.110.49:9090/api/dailycommits",
		"logout_url": "http://129.204.110.49:9090/api/logout"
	},
	"method_post": {
		"register_url": "http://129.204.110.49:9090/api/register",
		"login_url": "http://129.204.110.49:9090/api/login",
		"create_habit_url": "http://129.204.110.49:9090/api/habits",
		"create_punch_url": "http://129.204.110.49:9090/api/habits/{habitId}/{monthId}/{dayId}",
		"create_commit_url": "http://129.204.110.49:9090/api/dailycommits"
	},
	"method_put": {
		"edit_habit_url": "http://129.204.110.49:9090/api/habits/{habitId}",
		"edit_punch_url": "http://129.204.110.49:9090/api/habits/{habitId}/{monthId}/{dayId}",
		"edit_commit_url": "http://129.204.110.49:9090/api/dailycommits/{dailycommitId}"
	},
	"method_delete": {
		"delete_habit_url": "http://129.204.110.49:9090/api/habits/{habitId}",
		"delete_punch_url": "http://129.204.110.49:9090/api/habits/{habitId}/{monthId}/{dayId}",
		"delete_commit_url": "http://129.204.110.49:9090/api/dailycommits/{dailycommitId}"
	}
}`
