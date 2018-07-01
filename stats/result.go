package stats

// Result is the type holding the results from parsing an individual file
// they are accumulated into a report after all files are complete

import "sync"

type result struct {
	// stores an individual file's results, to be aggregated into the whole
	charCounts  []float64
	tokenCounts []float64
	keywords    map[string]int
}

// handles writing and closing the channel based on number of files
type resultWriter struct {
	num  int
	lock *sync.Mutex
	c    chan result
}

func (w *resultWriter) init(size int) {
	w.num = size
	w.lock = &sync.Mutex{}
	w.c = make(chan result)
}

// write the response to the channel, and if it is the last response,
// close the channel.
// Note on mutex: This particular mutex is highly unlikely to become a bottleneck
// as it's only called once parsing the whole file is done. Thus, the only way it
// could become a significant bottleneck is if the data being analyzed is contained
// in a very large number of small files.
func (w *resultWriter) write(r result) {
	w.c <- r
	w.lock.Lock()
	defer w.lock.Unlock()
	w.num--
	if w.num == 0 {
		defer close(w.c)
	}
}
