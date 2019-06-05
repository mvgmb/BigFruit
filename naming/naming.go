package naming

import (
	"log"

	pb "github.com/mvgmb/BigFruit/proto"
)

var (
	servicesTable = make(map[string][]*pb.AOR)
	roundRobin    = make(map[string]int)
)

func bind(bindRequest *pb.NamingServiceBindRequest) *pb.NamingServiceBindResponse {
	serviceName := bindRequest.ServiceName

	servicesTable[serviceName] = append(servicesTable[serviceName], bindRequest.Aor)
	roundRobin[serviceName] = -1

	log.Println("New service registered: " + serviceName)

	return &pb.NamingServiceBindResponse{}
}

func lookup(lookupRequest *pb.NamingServiceLookupRequest) *pb.NamingServiceLookupResponse {
	serviceName := lookupRequest.ServiceName

	if _, ok := servicesTable[serviceName]; !ok {
		return &pb.NamingServiceLookupResponse{Error: "404 - Service not found"}
	}

	roundRobin[serviceName]++
	if roundRobin[serviceName] >= len(servicesTable[serviceName]) {
		roundRobin[serviceName] = 0
	}

	return &pb.NamingServiceLookupResponse{Aor: servicesTable[serviceName][roundRobin[serviceName]]}
}

func lookupMany(lookupManyRequest *pb.NamingServiceLookupManyRequest) *pb.NamingServiceLookupManyResponse {
	serviceName := lookupManyRequest.ServiceName
	numberOfAor := lookupManyRequest.NumberOfAor

	if _, ok := servicesTable[serviceName]; !ok {
		return &pb.NamingServiceLookupManyResponse{Error: "404 - Service not found"}
	}
	if uint32(len(servicesTable[serviceName])) < numberOfAor {
		return &pb.NamingServiceLookupManyResponse{Error: "There are less Aor than required"}
	}

	var AorList []*pb.AOR
	for i := uint32(0); i < numberOfAor; i++ {
		roundRobin[serviceName]++
		if roundRobin[serviceName] >= len(servicesTable[serviceName]) {
			roundRobin[serviceName] = 0
		}
		AorList = append(AorList, servicesTable[serviceName][roundRobin[serviceName]])
	}
	return &pb.NamingServiceLookupManyResponse{AorList: AorList}
}

func lookupAll(lookupAllRequest *pb.NamingServiceLookupAllRequest) *pb.NamingServiceLookupAllResponse {
	serviceName := lookupAllRequest.ServiceName

	if _, ok := servicesTable[serviceName]; !ok {
		return &pb.NamingServiceLookupAllResponse{Error: "404 - Service not found"}
	}

	return &pb.NamingServiceLookupAllResponse{AorList: servicesTable[serviceName]}
}
