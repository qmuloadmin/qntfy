package stats

import "testing"

func BenchmarkUpdateDuplicate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		inc()
	}
	dupCount = 0
}

func BenchmarkCheckKeyword(b *testing.B) {
	keywords := new(set)
	keywords.init()
	keywords.add("lorum")
	keywords.add("ipsum")
	keywords.add("foo")
	keywords.add("bar")
	keywords.add("baz")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		keywords.contains("lorum")
	}
}

func runE2E(files []string) {
	ProcessFiles("output.tsv", "../example_data/testkeys.txt", files)
}

/*
BenchmarkE2ENormalFiles tests end-to-end benchmark, including
calculating final stats and writing to file, while using
normal (not long) data files, and 4 files at once (on a 4 core system)
*/
func BenchmarkE2ENormalFiles(b *testing.B) {
	files := []string{
		"../example_data/test1.txt",
		"../example_data/test2.txt",
		"../example_data/test3.txt",
		"../example_data/test4.txt",
	}
	for i := 0; i < b.N; i++ {
		runE2E(files)
	}
}

/*
BenchmarkE2EHugeFiles tests end-to-end benchmark of the entire process,
including calculating final stats and writing to file, while using
larger files, and less concurrency. 2 files at once.
*/
func BenchmarkE2ELargeFiles(b *testing.B) {
	files := []string{
		"../example_data/huge1.txt",
		"../example_data/huge2.txt",
	}
	for i := 0; i < b.N; i++ {
		runE2E(files)
	}
}

func BenchmarkE2EMoreFilesThanCores(b *testing.B) {
	files := []string{
		"../example_data/test1.txt",
		"../example_data/test2.txt",
		"../example_data/test3.txt",
		"../example_data/test4.txt",
		"../example_data/test5.txt",
		"../example_data/test6.txt",
	}
	for i := 0; i < b.N; i++ {
		runE2E(files)
	}
}

func BenchmarkE2EParagraphs(b *testing.B) {
	files := []string{
		"../example_data/pgraph1.txt",
		"../example_data/pgraph2.txt",
		"../example_data/pgraph3.txt",
	}
	for i := 0; i < b.N; i++ {
		ProcessFiles("output.tsv", "../example_data/testkeys.txt", files)
	}
}

func runProcessFile(name string) (r result) {
	seenLines := new(set)
	seenLines.init()
	kw := *new(set)
	kw.init()
	kw.add("ipsum")
	kw.add("sed")
	kw.add("et")
	w := new(resultWriter)
	w.init(1)
	go processFile(name, seenLines, kw, w)
	for r = range w.c {
	}
	return
}

func BenchmarkNormalFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := runProcessFile("../example_data/test1.txt")
		if r.charCounts[0] != 56.0 {
			b.Fail()
		} else if r.keywords["ipsum"] != 29 {
			b.Fail()
		}
	}
}

func BenchmarkLargeFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := runProcessFile("../example_data/huge1.txt")
		if r.charCounts[0] != 56.0 {
			b.Fail()
		} else if r.keywords["ipsum"] != 143 {
			b.Fail()
		}
	}
}

func BenchmarkParagraphFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r := runProcessFile("../example_data/pgraph2.txt")
		if r.charCounts[0] != 890 {
			b.Fail()
		} else if r.keywords["ipsum"] != 21 {
			b.Fail()
		}
	}
}

func BenchmarkProcessNormalLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		kws := *new(set)
		kws.init()
		kws.add("ipsum")
		counter := make(map[string]int)
		seenLines := new(set)
		seenLines.init()
		b.StartTimer() // restting the timer makes this take a long time, but is accurate
		ln, count := processLine("Lorem ipsum dolor sit amet, consectetur adipiscing elit.", seenLines, kws, counter)
		b.StopTimer()
		if ln != 56 {
			b.Fail()
		} else if count != 8 {
			b.Fail()
		}
	}
}

func BenchmarkProcessParagraphLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		kws := *new(set)
		kws.init()
		kws.add("ipsum")
		counter := make(map[string]int)
		seenLines := new(set)
		seenLines.init()
		b.StartTimer() // restting the timer makes this take a long time, but is accurate
		ln, count := processLine("Nullam faucibus odio non lorem porta, sit amet ultricies nisl malesuada. Nullam rutrum, ante eu luctus rutrum, est diam pharetra nulla, vitae accumsan justo quam vitae quam. Proin at leo et enim finibus fermentum a at felis. Suspendisse euismod vitae odio nec tincidunt. Interdum et malesuada fames ac ante ipsum primis in faucibus. Integer a quam in enim pharetra facilisis et vel urna. Mauris feugiat tempus vestibulum. Maecenas purus orci, mattis nec accumsan quis, venenatis sit amet nunc. Phasellus tincidunt dapibus tortor sit amet pulvinar. Curabitur pharetra interdum massa vel auctor. Donec fringilla vestibulum quam, a pellentesque elit luctus ut. Etiam euismod ultricies urna, id facilisis diam iaculis ut. Aenean volutpat nibh purus, sit amet placerat nisi lobortis vitae. Aliquam pretium consectetur sagittis. Vestibulum consequat tortor placerat metus semper mattis.",
			seenLines, kws, counter)
		b.StopTimer()
		if ln != 880 {
			b.Fail()
		} else if count != 130 {
			b.Fail()
		}
	}
}
