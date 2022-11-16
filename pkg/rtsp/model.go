package rtsp

type Conf struct {
	Source                     string `json:"source"`
	SourceProtocol             string `json:"sourceProtocol"`
	SourceAnyPortEnable        bool   `json:"sourceAnyPortEnable"`
	SourceFingerprint          string `json:"sourceFingerprint"`
	SourceOnDemand             bool   `json:"sourceOnDemand"`
	SourceOnDemandStartTimeout string `json:"sourceOnDemandStartTimeout"`
	SourceOnDemandCloseAfter   string `json:"sourceOnDemandCloseAfter"`
	SourceRedirect             string `json:"sourceRedirect"`
	DisablePublisherOverride   bool   `json:"disablePublisherOverride"`
	Fallback                   string `json:"fallback"`
	PublishUser                string `json:"publishUser"`
	PublishPass                string `json:"publishPass"`
	ReadUser                   string `json:"readUser"`
	ReadPass                   string `json:"readPass"`
	RunOnInit                  string `json:"runOnInit"`
	RunOnInitRestart           bool   `json:"runOnInitRestart"`
	RunOnDemand                string `json:"runOnDemand"`
	RunOnDemandRestart         bool   `json:"runOnDemandRestart"`
	RunOnDemandStartTimeout    string `json:"runOnDemandStartTimeout"`
	RunOnDemandCloseAfter      string `json:"runOnDemandCloseAfter"`
	RunOnReady                 string `json:"runOnReady"`
	RunOnReadyRestart          bool   `json:"runOnReadyRestart"`
	RunOnRead                  string `json:"runOnRead"`
	RunOnReadRestart           bool   `json:"runOnReadRestart"`
}
