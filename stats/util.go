package stats

// checks to see if err is nil
func ok(err error) bool {
	if err != nil {
		return false
	}
	return true
}

// incremenet the duplicate counter
func inc() {
	dupLock.Lock()
	defer dupLock.Unlock()
	dupCount++
}
