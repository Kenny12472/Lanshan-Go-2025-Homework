package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type User[A comparable, B comparable] interface {
	getInformation() (A, B)
	matchInformation() bool
}

type LoginFrame[A comparable, B comparable] struct {
	username A
	password B
}

var Records = [1][2]any{}

func (a *LoginFrame[A, B]) getInformation() (A, B) {
	return a.username, a.password
}

func (a *LoginFrame[A, B]) matchInformation() bool {
	if Records[0][0] == nil || Records[0][1] == nil {
		return false
	}
	return a.username == Records[0][0].(A) && a.password == Records[0][1].(B)
}

type registerFrame[A comparable, B comparable] struct {
	LoginFrame[A, B]
}

func (r *registerFrame[A, B]) getInformation() (A, B) {
	return r.username, r.password
}

func (r *registerFrame[A, B]) matchInformation() bool {
	if Records[0][0] == nil || Records[0][1] == nil {
		return false
	}
	return r.username == Records[0][0].(A) && r.password == Records[0][1].(B)
}

func (r *registerFrame[A, B]) register() bool {
	Records[0][0] = r.username
	Records[0][1] = r.password
	return true
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("注册登录界面")
	myWindow.Resize(fyne.NewSize(400, 250))

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("用户名")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("密码")

	loginButton := widget.NewButton("登录", func() {
		var login User[string, string] = &LoginFrame[string, string]{
			username: usernameEntry.Text,
			password: passwordEntry.Text,
		}
		username, password := login.getInformation()
		login = &LoginFrame[string, string]{username: username, password: password}

		if login.matchInformation() {
			widget.ShowPopUp(widget.NewLabel("登录成功"), myWindow.Canvas())
			MainWindow(myApp)
			myWindow.Close()
		} else {
			widget.ShowPopUp(widget.NewLabel("登录失败"), myWindow.Canvas())
		}
	})

	registerButton := widget.NewButton("注册", func() {
		var register User[string, string] = &registerFrame[string, string]{
			LoginFrame[string, string]{username: usernameEntry.Text, password: passwordEntry.Text},
		}
		username, password := register.getInformation()
		register = &registerFrame[string, string]{LoginFrame[string, string]{username: username, password: password}}

		register.(*registerFrame[string, string]).register()
		widget.ShowPopUp(widget.NewLabel("注册成功"), myWindow.Canvas())
	})

	buttons := container.NewHBox(loginButton, registerButton)
	content := container.NewVBox(
		usernameEntry,
		passwordEntry,
		buttons,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func MainWindow(a fyne.App) {
	mainWindow := a.NewWindow("主界面")
	mainWindow.Resize(fyne.NewSize(600, 400))

	label := widget.NewLabel("点击输入文本。")
	mainWindow.SetContent(container.NewCenter(label))

	mainWindow.Show()
}
