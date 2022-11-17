package rtsp

type Conf struct {
	Stream                     string `json:"stream"`
	SourceProtocol             string `json:"sourceProtocol"`
	SourceOnDemandStartTimeout string `json:"sourceOnDemandStartTimeout"`
	SourceOnDemandCloseAfter   string `json:"sourceOnDemandCloseAfter"`
	ReadUser                   string `json:"readUser"`
	ReadPass                   string `json:"readPass"`
	RunOnDemandStartTimeout    string `json:"runOnDemandStartTimeout"`
	RunOnDemandCloseAfter      string `json:"runOnDemandCloseAfter"`
	RunOnReady                 string `json:"runOnReady"`
}
