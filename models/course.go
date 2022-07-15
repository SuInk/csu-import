package models

import (
	"csu-import/utils"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GetCourse 获取课程信息
func GetCourse(client http.Client) (string, error) {
	// User ID from path `users/:id`
	// 在这里更改学期
	response, err := client.PostForm("https://csujwc.its.csu.edu.cn/jsxsd/kbxx/getKbxx.do", url.Values{"xnxq01id": {"2022-2023-1"}})
	if err != nil {
		return "", utils.ErrorJwc
	}
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Printf(string(body))
	return string(body), nil
}
