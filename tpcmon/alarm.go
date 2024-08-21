package tpcmon

type Dimensions struct {
	App       string `json:"app"`
	Product   string `json:"product"`
	AlertName string `json:"alertName"`
	HostName  string `json:"hostName"`
	Level     int    `json:"level"`
	Message   string `json:"message"`
}

type Alarm struct {
	Product    string     `json:"product"`
	Service    string     `json:"service"`
	App        string     `json:"app"`
	Host       string     `json:"host"`
	Level      int        `json:"level"`
	AlertName  string     `json:"alertName"`
	MetricName string     `json:"metricName"`
	Dimensions Dimensions `json:"dimentions"`
	Subject    string     `json:"subject"`
	Source     string     `json:"source"`
}

func NewAlarm(product, service, app, host, alertName, message, source string) *Alarm {
	return &Alarm{
		Product:    product,
		Service:    service,
		App:        app,
		Host:       host,
		Level:      1,
		AlertName:  alertName,
		MetricName: alertName,
		Subject:    alertName,
		Dimensions: Dimensions{
			App:       app,
			Product:   product,
			AlertName: alertName,
			HostName:  host,
			Level:     1,
			Message:   message,
		},
		Source: source,
	}
}
