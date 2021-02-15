package graph

func GetTransactionType(methodName string) string {
	switch methodName {
	case "Send":
		return "DEAL"
	case "ApplyRewards":
		return "REWARD"
	default:
		return "NETWORK_FEE"
	}
	// return ""
}
