package metrics

type NoOp struct {
}

func (n NoOp) Inc(action, status string, value int) {

}

func (n NoOp) ObserveDuration(action, status string, seconds float64) {
}
