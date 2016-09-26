package main

import (
	"fmt"
	"github.com/sclevine/agouti"
	"github.com/sclevine/agouti/api"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func check(err error) bool {
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// LoginWithUsername
func LoginWithUsername() {
	var err error
	dir, err := os.Getwd()
	check(err)
	log.Println("dir:", dir)
	driver := agouti.ChromeDriver(agouti.Desired(agouti.Capabilities{
		"chromeOptions": map[string]interface{}{
			"args": []string{
				// fmt.Sprintf("--proxy-server=http://%s", proxy),
				"--no-sandbox",
				"--force-device-scale-factor=0.5",
				"--kiosk",
				"--fullscreen",
				fmt.Sprintf("load-and-launch-app=%s/app", dir),
				fmt.Sprintf("--user-data-dir=%s/profile", dir),
			},
		},
	}))

	if err := driver.Start(); err != nil {
		log.Fatalf("Failed to start driver:%v", err)
	} else {
		log.Println("Started the driver")
	}

	page, err = driver.NewPage(agouti.Browser("chrome"))

	var windows []*api.Window
	for {

		time.Sleep(1 * 1e9)
		windows, err = page.Session().GetWindows()
		if len(windows) == 2 {
			break
		}
	}
	check(err)
	// page.CloseWindow()
	var windowID string
	for i, window := range windows {
		log.Println(i, window.ID)
		windowID = window.ID
	}
	page.CloseWindow()
	log.Println("LineID is", windowID)
	page.Session().SetWindow(windows[1])
	check(err)

	for {
		err := page.FindByID("login_maximize").Click()
		if check(err) == true {
			break
		}
		time.Sleep(1 * 1e9)
	}
	for {
		err1 := page.FindByID("login_email").Clear()
		err2 := page.FindByID("login_email").SendKeys(os.Getenv("username"))
		if check(err1) == true && check(err2) == true {
			break
		}
		time.Sleep(1 * 1e9)
	}
	for {

		err = page.FindByID("login_pwd").SendKeys(os.Getenv("password"))
		if check(err) == true {
			break
		}
		time.Sleep(1 * 1e9)
	}

	for {
		err := page.FindByID("login_btn").Click()
		if check(err) == true {
			break
		}
		time.Sleep(1 * 1e9)
	}
	// driver.Stop(
}

var page *agouti.Page

func sendMessages(text string) {

	for {
		err := page.FindByID("_chat_room_input").Clear()
		check(err)
		text = strings.Replace(text, "\n", "\uE008\uE007\uE008", -1)
		err1 := page.FindByID("_chat_room_input").SendKeys(text)
		err2 := page.FindByID("_chat_room_input").SendKeys("\uE007")
		if check(err1) == true && check(err2) == true {
			break
		}
		time.Sleep(1 * 1e9)
	}

}

func findFriendList() (int, *agouti.MultiSelection) {
	var i = 0
	var err error
	var a *agouti.MultiSelection
	for {
		err = page.FindByClass("mdLFT07Friends").Click()
		if check(err) {
			break
		}
		time.Sleep(1 * 1e9)
	}
	a = page.All("#contact_mode_contact_list li.mdMN02Li")
	i, err = a.Count()
	check(err)
	if i > 0 {
		log.Println("i:", i)
	}
	return i, a

}

func findTalkList() (int, *agouti.MultiSelection) {
	var i = 0
	var err error
	var a *agouti.MultiSelection
	for {
		err = page.FindByClass("mdLFT07Chats").Click()
		if check(err) {
			break
		}
		time.Sleep(1 * 1e9)
	}
	a = page.All("#_chat_list_body li")
	i, err = a.Count()
	check(err)
	if i > 0 {
		log.Println("i:", i)
	}
	return i, a

}
func findGroupList() (int, *agouti.MultiSelection) {
	var i = 0
	var err error
	var a *agouti.MultiSelection
	for {
		err = page.FindByClass("mdLFT07Groups").Click()
		if check(err) {
			break
		}
		time.Sleep(1 * 1e9)
	}
	a = page.All("#joinedGroup li")
	i, err = a.Count()
	check(err)
	if i > 0 {
		log.Println("i:", i)
	}
	return i, a

}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.Index(a, b) > -1 {
			return true
		}
	}
	return false
}
func main() {
	LoginWithUsername()
	friendCount, friendList := findGroupList()
	// friendCount, friendList := findFriendList()
	// exportNameList := []string{"MMM", "巧克", "回報"}
	exportNameList := []string{"-ignore"}

	for i := 0; i < friendCount; i++ {
		friend := friendList.At(i)
		name, err := friend.Attribute("title")
		if stringInSlice(name, exportNameList) {
			continue
		}
		check(err)
		log.Println(name)
		check(err)
		time.Sleep(1 * 1e9)
		friend.Click()
		sendMessages(os.Getenv("context"))

		sec, err := strconv.Atoi(os.Getenv("dealysec"))
		time.Sleep(time.Duration(sec) * 1e9)

	}

	time.Sleep(99999999 * 1e9)

}
