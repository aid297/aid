package sonarqube

type CodeQualityWrong struct {
	Errors []struct {
		Msg string `json:"msg"`
	} `json:"errors"`
}
