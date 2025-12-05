package gui

import (
	"third/information"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Run() {
	myApp := app.New()
	loginWindow := NewLoginWindow(myApp)
	loginWindow.ShowAndRun()
}

func NewLoginWindow(a fyne.App) fyne.Window {
	loginFrame := &information.LoginFrame[string, string]{}

	myWindow := a.NewWindow("登录注册界面")
	myWindow.Resize(fyne.NewSize(400, 300))

	usernameEntry := widget.NewEntry()
	usernameEntry.SetPlaceHolder("用户名")
	passwordEntry := widget.NewPasswordEntry()
	passwordEntry.SetPlaceHolder("密码")

	loginButton := widget.NewButton("登录", func() {
		loginFrame.GetInformation(usernameEntry.Text, passwordEntry.Text)
		if loginFrame.MatchInformation() {
			MainWindow(a)
			myWindow.Close()
		} else {
			widget.ShowPopUp(widget.NewLabel("登录失败"), myWindow.Canvas())
		}
	})

	registerButton := widget.NewButton("注册", func() {
		loginFrame.GetInformation(usernameEntry.Text, passwordEntry.Text)
		loginFrame.Register()
		widget.ShowPopUp(widget.NewLabel("注册成功"), myWindow.Canvas())
	})

	buttons := container.NewHBox(loginButton, registerButton)
	content := container.NewVBox(usernameEntry, passwordEntry, buttons)
	myWindow.SetContent(content)

	return myWindow
}

func MainWindow(a fyne.App) {
	mainWindow := a.NewWindow("new")
	mainWindow.Resize(fyne.NewSize(600, 400))
	label := widget.NewLabel("点击输入文本")
	mainWindow.SetContent(container.NewCenter(label))
	mainWindow.Show()
}
