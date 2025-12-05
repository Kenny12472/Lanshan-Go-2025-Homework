package information

var Records = [1][2]any{}

type LoginFrame[A comparable, B comparable] struct {
	Username A
	Password B
}

func (a *LoginFrame[A, B]) GetInformation(username A, password B) {
	a.Username = username
	a.Password = password
}

func (a *LoginFrame[A, B]) MatchInformation() bool {
	if Records[0][0] == nil || Records[0][1] == nil {
		return false
	}
	return a.Username == Records[0][0].(A) && a.Password == Records[0][1].(B)
}

func (a *LoginFrame[A, B]) Register() {
	Records[0][0] = a.Username
	Records[0][1] = a.Password
}
