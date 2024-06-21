package comms

type RoutingTable struct{}

func (rt *RoutingTable) userHubs(userID string) ([]string, error) {
	return []string{}, nil
}

func (rt *RoutingTable) userFDs(userID, hubID string) ([]int, error) {
	return []int{}, nil
}
