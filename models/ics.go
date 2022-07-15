package models

import (
	"csu-import/utils"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type S string
type Courses []struct {
	Jc       int    `json:"jc"`
	Title    string `json:"title"`
	Xq       int    `json:"xq"`
	Jx0404id string `json:"jx0404Id"`
}
type Schedule struct {
	Name     string
	Week     []int
	Day      int
	Time     int
	Teacher  string
	Location string
}

func GetIcs(courseList []Schedule) (os.File, error) {
	timetable := utils.Timetable
	// 在这里更改开学日期，其实可以前端传参给后端，这样只改前端就行了
	t, _ := time.Parse("2006-01-02", "2022-08-29")
	f, err := os.Create("temp.ics")
	if err != nil {
		fmt.Println(err)
	}
	_, err = f.Write([]byte("BEGIN:VCALENDAR\nVERSION:2.0\n"))
	if err != nil {
		return os.File{}, nil
	}
	//fmt.Println(n)
	content := ""
	for _, course := range courseList {
		for _, week := range course.Week {
			date := t.Add(time.Hour * 24 * time.Duration(7*(week-1)+course.Day-1))
			dateStr := strings.ReplaceAll(date.Format("2006-01-02"), "-", "")
			content += fmt.Sprintf("BEGIN:VEVENT\nSUMMARY:%s\nDTSTART;TZID=\"UTC+08:00\";VALUE=DATE-TIME:%sT%s\nDTEND;TZID=\"UTC+08:00\";VALUE=DATE-TIME:%sT%s\nLOCATION:%s\nDESCRIPTION:由https://import.suink.cn 导入\nBEGIN:VALARM\nTRIGGER:-PT18M\nACTION:DISPLAY\nEND:VALARM\nEND:VEVENT\n",
				course.Name, dateStr, timetable[course.Time-1].Start, dateStr, timetable[course.Time].End, course.Location)
		}

	}
	//fmt.Println(content)
	_, err = f.Write([]byte(content))
	if err != nil {
		return os.File{}, nil
	}
	_, err = f.Write([]byte("END:VCALENDAR"))
	if err != nil {
		return os.File{}, nil
	}
	defer f.Close()
	return *f, nil
}
func CourseParser(s string) ([]Schedule, error) {
	var courses Courses
	if err := json.Unmarshal([]byte(s), &courses); err != nil {
		//fmt.Println(courses)
	}
	//fmt.Println(courses)
	var schedules []Schedule
	// 遍历所有课程
	for _, course := range courses {
		if course.Xq != 0 {
			// 课程名称
			re, _ := regexp.Compile("课程名称：(.*)\n")
			courseName := re.FindStringSubmatch(course.Title)[1]
			//fmt.Println(courseName)
			// 周次
			re, _ = regexp.Compile("周次：(.*)\n")
			week := S(re.FindStringSubmatch(course.Title)[1])
			weekList, _ := week.WeekParser()
			//fmt.Println(week.WeekParser())
			// 星期
			day := course.Xq - 1
			//fmt.Println("星期", day)
			// 节次
			re, _ = regexp.Compile("节次：(.*)\n")
			jc := re.FindStringSubmatch(course.Title)[1][0:2]
			session, _ := strconv.Atoi(jc)
			//fmt.Println(jc)
			// 教师
			re, _ = regexp.Compile("上课教师：(.*)\n")
			teacher := re.FindStringSubmatch(course.Title)[1]
			//fmt.Println(teacher)
			// 教室
			re, _ = regexp.Compile("上课地点：(.*)\n")
			location := re.FindStringSubmatch(course.Title)[1]
			//fmt.Println(location)
			schedule := Schedule{
				courseName,
				weekList,
				day,
				session,
				teacher,
				location,
			}
			schedules = append(schedules, schedule)
		}

	}
	//fmt.Println(schedules)
	return schedules, nil
}

// WeekParser 将周次从文本转换为数组
func (s *S) WeekParser() ([]int, error) {
	var weekList []int
	// 1-4,7-9,12(周),13-18(单周) 这个真离谱
	zcs := strings.Split(string(*s), ",")
	for _, zc := range zcs {
		//fmt.Println(zc)
		ds := 0
		if strings.Contains(zc, "单") { //单周
			ds = 1
		} else if strings.Contains(string(*s), "双") { //双周
			ds = 2
		}
		zc = strings.ReplaceAll(zc, "周", "")
		zc = strings.ReplaceAll(zc, "单", "")
		zc = strings.ReplaceAll(zc, "双", "")
		zc = strings.ReplaceAll(zc, "()", "")
		//fmt.Println(zc)
		if strings.Contains(zc, "-") {
			zcs1 := strings.Split(zc, "-")
			start, _ := strconv.Atoi(zcs1[0])
			//fmt.Println(start)
			end, _ := strconv.Atoi(zcs1[1])
			switch ds { // 单双
			case 0:
				for i := start; i <= end; i++ {
					weekList = append(weekList, i)
				}
			case 1: // 单周
				for i := start; i <= end; i++ {
					if i&1 == 1 {
						weekList = append(weekList, i)
					}
				}
			case 2: // 双周
				for i := start; i <= end && i&1 == 0; i++ {
					if i&1 == 0 {
						weekList = append(weekList, i)
					}
				}

			}
		} else { // 12(周)
			week, _ := strconv.Atoi(zc)
			weekList = append(weekList, week)
		}
		// fmt.Println(ds)
	}
	//fmt.Println(weekList)
	return weekList, nil
}
