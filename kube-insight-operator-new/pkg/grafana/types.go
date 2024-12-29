package grafana

type Options struct {
	Name          string
	Namespace     string
	Labels        map[string]string
	AdminPassword string
	Storage       string
	PrometheusURL string
}

type Grafana struct {
	opts Options
}

func New(opts Options) *Grafana {
	if opts.Labels == nil {
		opts.Labels = make(map[string]string)
	}
	return &Grafana{
		opts: opts,
	}
}
