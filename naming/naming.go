package naming

import (
	"errors"
	"log"

	"github.com/mvgmb/BigFruit/util"
)

var (
	servicesTable = make(map[string][]util.AOR)
	roundRobin    = -1
)

func bind(aor *util.AOR) {
	servicesTable[aor.ID] = append(servicesTable[aor.ID], *aor)

	log.Println("New service registered: " + aor.ID)
	log.Println(servicesTable)
}

func lookup(fileName string) (*util.AOR, error) {
	if _, ok := servicesTable[fileName]; !ok {
		return nil, errors.New("404 - Service not found")
	}

	if roundRobin+1 == len(servicesTable[fileName]) {
		roundRobin = 0
	} else {
		roundRobin++
	}

	return &servicesTable[fileName][roundRobin], nil
}
