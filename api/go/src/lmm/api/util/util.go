package util

// Retry retries f for count times, retries infinitly if count is less than zero
func Retry(count int, f func() error) error {
	var err error
	for {
		err = f()
		if err == nil || count == 0 {
			break
		}
		if count > 0 {
			count--
		}
	}
	return err
}
