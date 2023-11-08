package helper

// var ErrUserNotFound = errors.New("record not found")

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}
