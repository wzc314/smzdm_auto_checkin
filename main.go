package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const username string = "your_username"
const password string = "your_password"

var header http.Header = map[string][]string{
	"Content-Type": {"application/x-www-form-urlencoded"},
	"User-Agent":   {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.124 Safari/537.36"},
	"Referer":      {"http://www.smzdm.com/"},
	"Origin":       {"http://www.smzdm.com/"},
}

var client http.Client

func main() {
	jar, err := cookiejar.New(nil)
	check_error(err)
	client.Jar = jar

	login()
	checkin()
	is_checkin()
	logout()
}

func login() {
	login_url := "https://zhiyou.smzdm.com/user/login/ajax_check"

	data := url.Values{}
	data.Set("username", username)
	data.Set("password", password)
	data.Set("rememberme", "on")
	data.Set("redirect_url", "http://www.smzdm.com")

	request, err := http.NewRequest("POST", login_url, strings.NewReader(data.Encode()))
	check_error(err)
	request.Header = header

	response, err := client.Do(request)
	defer response.Body.Close()
	check_error(err)
}

func checkin() {
	checkin_url := "http://zhiyou.smzdm.com/user/checkin/jsonp_checkin"

	request, err := http.NewRequest("POST", checkin_url, nil)
	check_error(err)
	request.Header = header

	response, err := client.Do(request)
	defer response.Body.Close()
	check_error(err)
}

func is_checkin() {
	is_checkin_url := "http://zhiyou.smzdm.com/user/info/jsonp_get_current?"

	request, err := http.NewRequest("GET", is_checkin_url, nil)
	check_error(err)
	request.Header = header

	response, err := client.Do(request)
	defer response.Body.Close()
	check_error(err)

	body, err := ioutil.ReadAll(response.Body)
	check_error(err)

	time_str := time.Now().Format("2006-01-02 15:04:05")
	var has_checkin string

	re := regexp.MustCompile(`"has_checkin":(\w+),`)
	params := re.FindSubmatch(body)
	if params == nil {
		fmt.Println("FindSubmatch returns nil.")
		return
	} else if string(params[1]) == "true" {
		has_checkin = "SUCCEED"
	} else {
		has_checkin = "FAILED"
	}

	has_checkin = strings.Join([]string{time_str, has_checkin, "\n"}, "\t")
	err = ioutil.WriteFile("checkin.log", []byte(has_checkin), 0644)
	check_error(err)
}

func logout() {
	logout_url := "http://zhiyou.smzdm.com/user/logout"

	request, err := http.NewRequest("POST", logout_url, nil)
	request.Header = header
	check_error(err)

	response, err := client.Do(request)
	defer response.Body.Close()
	check_error(err)
}

func check_error(e error) {
	if e != nil {
		panic(e)
	}
}
