package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/liuyh73/dailyhub.service/db"
	"github.com/liuyh73/dailyhub.service/model"
)

var Api = []string{}

type RespData struct {
	Status bool        `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

func writeResp(status bool, msg string, data interface{}) []byte {
	RespData := RespData{}
	RespData.Status = status
	RespData.Msg = msg
	RespData.Data = data
	respose, err := json.Marshal(RespData)
	if err != nil {
		log.Fatalln(err)
	}
	return respose
}

func jsonDecode(r io.Reader, val interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	return dec.Decode(val)
}

// 注册
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Print("Error in request")
		w.Write(writeResp(false, "Error in request", Token{}))
		return
	}
	has, err, _ := db.GetUserProfile(r.Form.Get("username"))
	checkErr(err)
	if !has && err == nil {
		token, err := createToken([]byte(SecretKey), Issuer, r.Form.Get("username"))
		fmt.Println(token)
		checkErr(err)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println("Error marshal the token")
			w.Write(writeResp(false, "Error marshal the token", Token{}))
		}
		w.WriteHeader(http.StatusOK)
		rows, err := db.InsertUserTokenItem(r.Form.Get("username"), token.DH_TOKEN)
		log.Println(rows)
		checkErr(err)
		rows, err = db.InsertUserProfile(model.Profile{
			Username: r.Form.Get("username"),
			Password: r.Form.Get("password"),
		})
		log.Println(rows)
		checkErr(err)
		w.Write(writeResp(true, "Succeed to register", token))
	} else {
		w.Write(writeResp(false, "Failed to register, the username has been occupied", Token{}))
	}
}

// 登录
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	checkErr(err)
	fmt.Println(r.Form.Get("username"))
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Print("Error in request")
		w.Write(writeResp(false, "Error in request", Token{}))
		return
	}
	user := &model.Profile{
		Username: r.Form.Get("username"),
		Password: r.Form.Get("password"),
	}
	fmt.Printf("%#v\n", user)
	has, err := db.Engine.Get(user)
	log.Println(has, err)
	if !has || err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Println("Error logging in")
		w.Write(writeResp(false, "Error logging in", Token{}))
		return
	}

	token, err := createToken([]byte(SecretKey), Issuer, user.Username)
	checkErr(err)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error marshal the token")
		w.Write(writeResp(false, "Error marshal the token", Token{}))
		return
	}
	w.WriteHeader(http.StatusOK)

	rows, err := db.DeleteUserTokenItem(user.Username)
	log.Println(rows)
	checkErr(err)
	rows, err = db.InsertUserTokenItem(user.Username, token.DH_TOKEN)
	log.Println(rows)
	checkErr(err)
	w.Write(writeResp(true, "Succeed to login", token))
}

// 退出
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	_, err := db.DeleteUserTokenItem(r.Context().Value("username").(string))
	checkErr(err)
	w.Write(writeResp(true, "Succeed to logout", Token{}))
}

// 用户信息
func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	has, err, profile := db.GetUserProfile(r.Context().Value("username").(string))
	checkErr(err)
	if has {
		w.Write(writeResp(true, "Succeed to get profile", profile))
	} else {
		w.Write(writeResp(false, "There is no profile about the user", profile))
	}
}

// 用户所有习惯
func GetHabitsHandler(w http.ResponseWriter, r *http.Request) {
	has, err, habits := db.GetUserHabits(r.Context().Value("username").(string))
	checkErr(err)
	if has {
		w.Write(writeResp(true, "Succeed to get user's habits", habits))
	} else {
		w.Write(writeResp(false, "There is no habits about the user", habits))
	}
}

// 根据id获取用户具体习惯详情
func GetHabitHandler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.RequestURI, "/")
	id := params[3]
	log.Println(id)
	has, err, habit := db.GetUserHabit(r.Context().Value("username").(string), id)
	checkErr(err)
	if has {
		w.Write(writeResp(true, "Succeed to get the habit", habit))
	} else {
		w.Write(writeResp(false, "There is no habit with the id about the user", habit))
	}
}

// 根据monthid获取用户该月份打卡详情
func GetMonthHandler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.RequestURI, "/")
	habitId := params[3]
	monthId := params[4]
	has, err, month := db.GetUserHabitMonth(r.Context().Value("username").(string), habitId, monthId)
	checkErr(err)
	if has {
		w.Write(writeResp(true, "Succeed to get the month info", month))
	} else {
		w.Write(writeResp(false, "Fail to get month info", month))
	}
}

// 根据dayid获取用户该日打卡详情
func GetDayHandler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.RequestURI, "/")
	habitId := params[3]
	monthId := params[4]
	dayId := params[5]
	has, err, day := db.GetUserHabitMonthDay(r.Context().Value("username").(string), habitId, monthId, dayId)
	checkErr(err)
	if has {
		w.Write(writeResp(true, "Succeed to get the day info", day))
	} else {
		w.Write(writeResp(false, "Fail to get the day info", day))
	}
}

// 创建habit
func PostHabitsHandler(w http.ResponseWriter, r *http.Request) {
	habit := &model.Habit{}
	err := jsonDecode(r.Body, habit)
	checkErr(err)
	if err == nil {
		_, err, habitId := db.InsertUserHabit(r.Context().Value("username").(string), *habit)
		checkErr(err)
		if err == nil {
			habit.Id = habitId
			w.Write(writeResp(true, "Succeed to create the habit", habit))
		} else {
			w.Write(writeResp(false, "Fail to create the habit", habit))
		}
	} else {
		w.Write(writeResp(false, "The request body need to be habit type", habit))
	}
}

// 打卡
func PostDayHandler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.RequestURI, "/")
	habitId := params[3]
	monthId := params[4]
	dayId := params[5]
	day := &model.Day{}
	err := jsonDecode(r.Body, day)
	day.Id = dayId
	checkErr(err)
	if err == nil {
		_, err := db.InsertUserHabitMonthDay(r.Context().Value("username").(string), habitId, monthId, *day)
		checkErr(err)
		if err == nil {
			w.Write(writeResp(true, "Succeed to punch in the day", day))
		} else {
			w.Write(writeResp(false, "Fail to punch in the day", day))
		}
	} else {
		w.Write(writeResp(false, "Failed to punch in the day", day))
	}
}

// 编辑、归档/结束、重要性
func PutHabitHandler(w http.ResponseWriter, r *http.Request) {
	habit := &model.Habit{}
	err := jsonDecode(r.Body, habit)
	checkErr(err)
	if err == nil {
		_, err = db.UpdateUserHabit(r.Context().Value("username").(string), *habit)
		checkErr(err)
		if err == nil {
			w.Write(writeResp(true, "Succeed to update the habit", habit))
		} else {
			w.Write(writeResp(false, "Fail to update the habit", habit))
		}
	} else {
		w.Write(writeResp(false, "The request body need to be habit type", habit))
	}
}

// 更新打卡信息
func PutDayHandler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.RequestURI, "/")
	habitId := params[3]
	monthId := params[4]
	day := &model.Day{}
	err := jsonDecode(r.Body, day)
	checkErr(err)
	if err == nil {
		_, err := db.UpdateUserHabitMonthDay(r.Context().Value("username").(string), habitId, monthId, *day)
		checkErr(err)
		if err == nil {
			w.Write(writeResp(true, "Succeed to update the punch info", day))
		} else {
			w.Write(writeResp(false, "Fail to update the punch info", day))
		}
	} else {
		w.Write(writeResp(false, "The request body need to be habit type", day))
	}
}

// 删除习惯
func DeleteHabitHandler(w http.ResponseWriter, r *http.Request) {
	habitId := strings.Split(r.RequestURI, "/")[3]
	_, err := db.DeleteUserHabit(r.Context().Value("username").(string), habitId)
	checkErr(err)
	if err == nil {
		w.Write(writeResp(true, "Succeed to delete the habit", nil))
	} else {
		w.Write(writeResp(false, "Failed to delete the habit", nil))
	}
}

// 取消打卡
func DeleteDayHandler(w http.ResponseWriter, r *http.Request) {
	params := strings.Split(r.RequestURI, "/")
	habitId := params[3]
	monthId := params[4]
	dayId := params[5]
	_, err := db.DeleteUserHabitMonthDay(r.Context().Value("username").(string), habitId, monthId, dayId)
	checkErr(err)
	if err == nil {
		w.Write(writeResp(true, "Succeed to delete the punch", nil))
	} else {
		w.Write(writeResp(false, "Failed to delete the punch", nil))
	}
}
