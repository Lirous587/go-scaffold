package metrics

type NoOp struct {
}

func (n NoOp) Inc(key string, value int) {

}

func (n NoOp) ObserveDuration(key string, seconds float64) {
}
