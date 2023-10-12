package ant

type ants struct {
	colony []ant
}

func (a *ants) addAnt(ant ant) {
	a.colony = append(a.colony, ant)
}
