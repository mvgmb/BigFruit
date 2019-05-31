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

func lookupMany(fileName string, numberOfNodes int) (*[]util.AOR, error) {
	var listOfNodes []util.AOR
	if _, ok := servicesTable[fileName]; !ok {
		return nil, errors.New("404 - Service not found")
	}
	if len(servicesTable[fileName]) < numberOfNodes {
		return nil, errors.New("there are less nodes than required")
	}
	for i := 0; i < numberOfNodes; i++ {
		if roundRobin+1 == len(servicesTable[fileName]) {
			roundRobin = 0
		} else {
			roundRobin++
		}
		listOfNodes = append(listOfNodes, servicesTable[fileName][roundRobin])
	}
	return &listOfNodes, nil
}

func lookupAll(fileName string) (*[]util.AOR, error) {
	if _, ok := servicesTable[fileName]; !ok {
		return nil, errors.New("404 - Service not found")
	}
	var listOfNodes = servicesTable[fileName]

	return &listOfNodes, nil
}
