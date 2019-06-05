package naming

import (
	"log"

	naming "github.com/mvgmb/BigFruit/proto/naming"
)

var (
	servicesTable = make(map[string][]*naming.AOR)
	roundRobin    = make(map[string]int)
)

func bind(bindRequest *naming.BindRequest) *naming.BindResponse {
	serviceName := bindRequest.ServiceName

	servicesTable[serviceName] = append(servicesTable[serviceName], bindRequest.Aor)
	roundRobin[serviceName] = -1

	log.Println("New service registered: " + serviceName)

	return &naming.BindResponse{}
}

func lookup(lookupRequest *naming.LookupRequest) *naming.LookupResponse {
	serviceName := lookupRequest.ServiceName

	if _, ok := servicesTable[serviceName]; !ok {
		return &naming.LookupResponse{Error: "404 - Service not found"}
	}

	roundRobin[serviceName]++
	if roundRobin[serviceName] >= len(servicesTable[serviceName]) {
		roundRobin[serviceName] = 0
	}

	return &naming.LookupResponse{Aor: servicesTable[serviceName][roundRobin[serviceName]]}
}

func lookupMany(lookupManyRequest *naming.LookupManyRequest) *naming.LookupManyResponse {
	serviceName := lookupManyRequest.ServiceName
	numberOfAor := lookupManyRequest.NumberOfAor

	if _, ok := servicesTable[serviceName]; !ok {
		return &naming.LookupManyResponse{Error: "404 - Service not found"}
	}
	if uint32(len(servicesTable[serviceName])) < numberOfAor {
		return &naming.LookupManyResponse{Error: "There are less Aor than required"}
	}

	var AorList []*naming.AOR
	for i := uint32(0); i < numberOfAor; i++ {
		roundRobin[serviceName]++
		if roundRobin[serviceName] >= len(servicesTable[serviceName]) {
			roundRobin[serviceName] = 0
		}
		AorList = append(AorList, servicesTable[serviceName][roundRobin[serviceName]])
	}
	return &naming.LookupManyResponse{AorList: AorList}
}

func lookupAll(lookupAllRequest *naming.LookupAllRequest) *naming.LookupAllResponse {
	serviceName := lookupAllRequest.ServiceName

	if _, ok := servicesTable[serviceName]; !ok {
		return &naming.LookupAllResponse{Error: "404 - Service not found"}
	}

	return &naming.LookupAllResponse{AorList: servicesTable[serviceName]}
}
