package rtspsimpleserver

type SConf struct {
	Stream string `json:"stream"`
	Conf   Conf
	Id     int
}

// Conf - структура вида поля "conf" в ответе с сервера
type Conf struct {
	SourceProtocol             string `json:"sourceProtocol"`
	SourceOnDemandStartTimeout string `json:"sourceOnDemandStartTimeout"`
	SourceOnDemandCloseAfter   string `json:"sourceOnDemandCloseAfter"`
	ReadUser                   string `json:"readUser"`
	ReadPass                   string `json:"readPass"`
	RunOnDemandStartTimeout    string `json:"runOnDemandStartTimeout"`
	RunOnDemandCloseAfter      string `json:"runOnDemandCloseAfter"`
	RunOnReady                 string `json:"runOnReady"`
	Source                     string `json:"source"`
}
