package helper

// var ErrUserNotFound = errors.New("record not found")

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// func RecordNotFoundErr(err error, customErr error) error {
//     if gorm.IsRecordNotFoundError(err) {
//         return customErr
//     }
//     return err
// }