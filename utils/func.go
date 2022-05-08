package utils

func TryExec(f func()) {
	if f != nil {
		f()
	}
}
