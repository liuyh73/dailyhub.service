package db

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/liuyh73/dailyhub.service/model"
	yaml "gopkg.in/yaml.v2"
)

type config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Schema   string `yaml:"schema"`
}

var Engine *xorm.Engine
var timeFormat string = "2006-01-02 15:04"

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func init() {
	conf := &config{}
	confFile, err := ioutil.ReadFile("db/conf.yml")
	checkErr(err)
	// fmt.Println(string(confFile))
	err = yaml.Unmarshal(confFile, conf)
	checkErr(err)
	dataSourceName := conf.Username + ":" + conf.Password + "@tcp(" + conf.Host + ":" + conf.Port + ")/" + conf.Schema + "?charset=utf8"
	// fmt.Println(dataSourceName)
	Engine, err = xorm.NewEngine("mysql", dataSourceName)
	checkErr(err)
}

func GetUserTokenItem(username string) (bool, error, *model.TokenItem) {
	tokenItem := &model.TokenItem{
		Username: username,
	}
	ok, err := Engine.Get(tokenItem)
	return ok, err, tokenItem
}

func GetUserProfile(username string) (bool, error, *model.Profile) {
	profile := &model.Profile{
		Username: username,
	}
	ok, err := Engine.Get(profile)
	return ok, err, profile
}

func GetUserHabits(username string) (bool, error, []*model.Habit) {
	ok, err, profile := GetUserProfile(username)
	var habits []*model.Habit
	var habit *model.Habit
	for _, habitId := range profile.Habits {
		ok, err, habit = GetUserHabit(username, habitId)
		habits = append(habits, habit)
	}
	return ok, err, habits
}

func GetUserHabit(username, habitId string) (bool, error, *model.Habit) {
	tid := username + "-" + habitId
	habit := &model.Habit{
		Id: tid,
	}
	ok, err := Engine.Get(habit)
	habit.Id = habitId
	return ok, err, habit
}

func GetUserHabitMonth(username, habitId, monthId string) (bool, error, *model.Month) {
	tid := username + "-" + habitId + "-" + monthId
	month := &model.Month{
		Id: tid,
	}
	ok, err := Engine.Get(month)
	month.Id = monthId
	return ok, err, month
}

func GetUserHabitMonthDay(username, habitId, monthId, dayId string) (bool, error, *model.Day) {
	tid := username + "-" + habitId + "-" + monthId + "-" + dayId
	day := &model.Day{
		Id: tid,
	}
	ok, err := Engine.Get(day)
	day.Id = dayId
	return ok, err, day
}

func GetUserDailyCommits(username string) (error, []model.DailyCommit) {
	dailyCommits := make([]model.DailyCommit, 0)
	allDailyCommits := make([]model.DailyCommit, 0)
	err := Engine.Find(&allDailyCommits)
	for _, dailyCommit := range allDailyCommits {
		if strings.HasPrefix(dailyCommit.Id, username) {
			dailyCommit.Id = strings.TrimPrefix(dailyCommit.Id, username+"-")
			dailyCommits = append(dailyCommits, dailyCommit)
		}
	}
	return err, dailyCommits
}

func InsertUserTokenItem(username, dh_token string) (int64, error) {
	tokenItem := model.TokenItem{
		Username: username,
		DH_TOKEN: dh_token,
	}
	return Engine.Insert(tokenItem)
}

func InsertUserProfile(profile model.Profile) (int64, error) {
	return Engine.Insert(profile)
}

func InsertUserHabit(username string, habit model.Habit) (int64, error, string) {
	_, err, profile := GetUserProfile(username)
	checkErr(err)
	id := 0
	if len(profile.Habits) > 0 {
		id, _ = strconv.Atoi(profile.Habits[len(profile.Habits)-1])
	}
	habitId := strconv.Itoa(id + 1)
	// log.Println(habitId)
	profile.Habits = append(profile.Habits, habitId)
	_, err = UpdateUserProfile(*profile)
	checkErr(err)
	habit.Id = username + "-" + habitId
	rows, err := Engine.Insert(habit)
	return rows, err, habitId
}

func getFebDays(year string) int {
	y, err := strconv.Atoi(year)
	checkErr(err)
	if (y%4 == 0 && y%100 != 0) || y%400 == 0 {
		return 29
	}
	return 28
}

func getDays(monthId string) int {
	year := strings.Split(monthId, "-")[0]
	month := strings.Split(monthId, "-")[1]
	// log.Println(month)
	switch month {
	case "01", "03", "05", "07", "08", "10", "12":
		return 31
	case "04", "06", "09", "11":
		return 30
	default:
		return getFebDays(year)
	}
}

func InsertUserHabitMonth(username, habitId, monthId, dayId string) (int64, error) {
	month := model.Month{
		Id:          username + "-" + habitId + "-" + monthId,
		PlanPunch:   getDays(monthId),
		ActualPunch: 1,
		MissPunch:   getDays(monthId) - 1,
		Days:        []string{dayId},
	}
	return Engine.Insert(month)
}

// 打卡
func InsertUserHabitMonthDay(username, habitId, monthId string, day model.Day) (int64, error) {
	has, err, month := GetUserHabitMonth(username, habitId, monthId)
	checkErr(err)
	if has && err == nil {
		for _, dayId := range month.Days {
			if dayId == day.Id {
				return 0, nil
			}
		}
		month.Days = append(month.Days, day.Id)
		month.Id = username + "-" + habitId + "-" + monthId
		month.ActualPunch = month.ActualPunch + 1
		month.MissPunch = month.MissPunch - 1
		Engine.Where("id = ?", month.Id).Update(month)
	} else {
		_, err = InsertUserHabitMonth(username, habitId, monthId, day.Id)
		checkErr(err)
	}

	has, err, habit := GetUserHabit(username, habitId)
	checkErr(err)
	if has && err == nil {
		habit.LastRecentPunchTime = habit.RecentPunchTime
		habit.RecentPunchTime = day.Time
		habit.TotalPunch = habit.TotalPunch + 1
		lastRecentPunchDate, err := time.Parse(timeFormat, habit.LastRecentPunchTime)
		checkErr(err)
		lastRecentPunchDate = lastRecentPunchDate.AddDate(0, 0, 1)
		lRPD := strings.Split(lastRecentPunchDate.String(), " ")[0]
		rPD := strings.Split(habit.RecentPunchTime, " ")[0]
		if lRPD == rPD {
			habit.CurrcPunch = habit.CurrcPunch + 1
		} else if lRPD < rPD {
			if habit.OncecPunch < habit.CurrcPunch {
				habit.OncecPunch = habit.CurrcPunch
			}
			habit.CurrcPunch = 1
		}
		_, err = UpdateUserHabit(username, *habit)
		checkErr(err)
	}
	day.Id = username + "-" + habitId + "-" + monthId + "-" + day.Id
	return Engine.Insert(day)
}

// 创建dailyCommit
func InsertUserDailyCommit(username string, dailyCommit model.DailyCommit) (int64, error, string) {
	err, dailyCommits := GetUserDailyCommits(username)
	checkErr(err)
	if err == nil {
		lastId := "0"
		if len(dailyCommits) > 0 {
			lastId = dailyCommits[len(dailyCommits)-1].Id
		}
		lastIdInt, err := strconv.Atoi(lastId)
		checkErr(err)
		if err == nil {
			dailyCommitId := strconv.Itoa(lastIdInt + 1)
			dailyCommit.Id = username + "-" + dailyCommitId
			rows, err := Engine.Insert(dailyCommit)
			return rows, err, dailyCommitId
		}
		return 0, err, "0"
	}
	return 0, err, "0"
}

func UpdateUserTokenItem(username, dh_token string) (int64, error) {
	tokenItem := model.TokenItem{
		Username: username,
		DH_TOKEN: dh_token,
	}
	return Engine.Where("username = ?", username).Update(tokenItem)
}

func UpdateUserProfile(profile model.Profile) (int64, error) {
	return Engine.Where("username = ?", profile.Username).Cols("avatar,description,habits").Update(&profile)
}

// 编辑、重要性、归档
func UpdateUserHabit(username string, habit model.Habit) (int64, error) {
	habit.Id = username + "-" + habit.Id
	return Engine.Where("id=?", habit.Id).AllCols().Update(habit)
}

// 修改打卡信息
func UpdateUserHabitMonthDay(username, habitId, monthId string, day model.Day) (int64, error) {
	day.Id = username + "-" + habitId + "-" + monthId + "-" + day.Id
	return Engine.Where("id=?", day.Id).Update(day)
}

// 修改dailyCommit
func UpdateUserDailyCommit(username string, dailyCommit model.DailyCommit) (int64, error) {
	dailyCommit.Id = username + "-" + dailyCommit.Id
	return Engine.Where("id=?", dailyCommit.Id).AllCols().Update(dailyCommit)
}

// 删除token
func DeleteUserTokenItem(username string) (int64, error) {
	tokenItem := model.TokenItem{
		Username: username,
	}
	return Engine.Delete(tokenItem)
}

// 删除习惯
func DeleteUserHabit(username, habitId string) (int64, error) {
	tid := username + "-" + habitId
	// 删除habit
	habit := model.Habit{
		Id: tid,
	}
	rows, err := Engine.Delete(habit)
	months := make([]model.Month, 0)
	days := make([]model.Day, 0)
	// 删除与当前habit相关的打卡信息
	err = Engine.Find(&months)
	checkErr(err)
	err = Engine.Find(&days)
	checkErr(err)
	for _, month := range months {
		if strings.HasPrefix(month.Id, tid) {
			rows, err = Engine.Delete(month)
		}
	}
	for _, day := range days {
		if strings.HasPrefix(day.Id, tid) {
			rows, err = Engine.Delete(day)
		}
	}
	// 删除profile中的habit记录
	ok, err, profile := GetUserProfile(username)
	checkErr(err)
	if ok && err == nil {
		for i, hbt := range profile.Habits {
			if hbt == habitId {
				profile.Habits = append(profile.Habits[:i], profile.Habits[i+1:]...)
				break
			}
		}
		rows, err = UpdateUserProfile(*profile)
		// log.Println(rows)
		checkErr(err)
	}
	return rows, err
}

// 取消打卡
func DeleteUserHabitMonthDay(username, habitId, monthId, dayId string) (int64, error) {
	ok, err, month := GetUserHabitMonth(username, habitId, monthId)
	checkErr(err)
	if ok && err == nil {
		for i, day := range month.Days {
			if day == dayId {
				month.Days = append(month.Days[:i], month.Days[i+1:]...)
				break
			}
		}
		month.Id = username + "-" + habitId + "-" + monthId
		month.ActualPunch = month.ActualPunch - 1
		month.MissPunch = month.MissPunch + 1
		fmt.Printf("%#v\n", month)
		_, err := Engine.Where("id = ?", month.Id).Cols("actual_punch, days, miss_punch").Update(month)
		// log.Println(rows)
		checkErr(err)
	}

	has, err, habit := GetUserHabit(username, habitId)
	checkErr(err)
	if has && err == nil {
		habit.TotalPunch = habit.TotalPunch - 1
		habit.RecentPunchTime = habit.LastRecentPunchTime
		habit.CurrcPunch = habit.CurrcPunch - 1

		_, err = UpdateUserHabit(username, *habit)
		checkErr(err)
	}
	tid := username + "-" + habitId + "-" + monthId + "-" + dayId
	day := model.Day{
		Id: tid,
	}
	return Engine.Delete(day)
}

// 删除dailyCommit
func DeleteUserDailyCommit(username, dailyCommitId string) (int64, error) {
	dailyCommit := model.DailyCommit{
		Id: username + "-" + dailyCommitId,
	}
	return Engine.Delete(dailyCommit)
}
