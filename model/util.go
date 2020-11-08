package model

func (p *AllUser) OrderToInt() int {

	switch p.OrderDir {
	case "asc":
		return 1
	case "desc":
		return -1
	default:
		break
	}

	return 1
}
