package catalog

func Describe(box string) string {
	switch box {
	case "inbound":
		return "待投递"
	case "sent":
		return "已投递"
	case "bounced":
		return "退信"
	case "deferred":
		return "延期"
	default:
		return box
	}
}
