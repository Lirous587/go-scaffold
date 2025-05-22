package metrics

type Client interface {
	Inc(action, status string, value int)
	ObserveDuration(action, status string, seconds float64)
}
