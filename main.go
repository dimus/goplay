package main

import (
	"log"
	"os"
	"runtime/pprof"
	"sync"

	"github.com/dimus/goplay/grammar"
)

func main() {
	pprof.StartCPUProfile(os.Stdout)
	defer pprof.StopCPUProfile()
	var wg sync.WaitGroup
	wg.Add(12)
	in := make(chan string)
	go setInput(in)
	for i := 0; i < 12; i++ {
		go worker(in, &wg)
	}
	wg.Wait()
}

func setInput(in chan<- string) {
	str := []string{
		"Pseudocercospora",
		"Pseudocercospora Speg.",
		"Döringina Ihering 1929 (synonym)",
		"Pseudocercospora Speg., Francis Jack.-Drake.",
		"Aaaba de Laubenfels, 1936",
		"Abbottia F. von Mueller, 1875",
		"Abella von Heyden, 1826",
		"Micropleura v Linstow 1906",
		"Pseudocercospora Speg. 1910",
		"Pseudocercospora Spegazzini, 1910",
		"Platypus bicaudatulus Schedl (1935h)",
		"Platypus bicaudatulus Schedl (1935)",
		"Tridentella tangeroae Bruce, 198?",
		"Rhynchonellidae d'Orbigny 1847",
		"Ataladoris Iredale & O'Donoghue 1923",
		"Saxo-Fridericia R. H. Schomb.",
		"Anteplana le Renard 1995",
		"Candinia le Renard, Sabelli & Taviani 1996",
		"Polypodium le Sourdianum Fourn.",
		"Ca Dyar 1914",
		"Ea Distant 1911",
		"Ge Nicéville 1895",
		"Ia Thomas 1902",
		"Io Lea 1831",
		"Io Blanchard 1852",
		"Ix Bergroth 1916",
		"Lo Seale 1906",
		"Oa Girault 1929",
		"Ra Whitley 1931",
		"Ty Bory de St. Vincent 1827",
		"Ua Girault 1929",
		"Aa Baker 1940",
		"Ja Uéno 1955",
		"Zu Walters & Fitch 1960",
		"La Bleszynski 1966",
		"Qu Durkoop",
		"As Slipinski 1982",
		"Ba Solem 1983",
		"Poaceae subtrib. Scolochloinae Soreng",
		"Zygophyllaceae subfam. Tribuloideae D.M.Porter",
	}
	count := 0
	for count < 200000 {
		for _, v := range str {
			count++
			in <- v
		}
	}
	close(in)
}

func worker(in <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	gnp := &grammar.GNParser{}
	gnp.Init()
	for v := range in {
		gnp.Buffer = v
		gnp.Reset()
		err := gnp.Parse()
		if err != nil {
			log.Printf("No parse for '%s': %s\n", v, err)
		}
		gnp.PrintSyntaxTree()
	}
}
