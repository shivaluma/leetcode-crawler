package entity

type RecentAcSubmission struct {
	ID        string
	Title     string
	TitleSlug string
	Timestamp string
	LangName  string
	Lang      string
	Runtime   string
	Memory    string
	Url       string
}

type Submission struct {
	Runtime             int     `json:"runtime"`
	RuntimeDisplay      string  `json:"runtimeDisplay"`
	RuntimePercentile   float64 `json:"runtimePercentile"`
	RuntimeDistribution string  `json:"runtimeDistribution"`
	Memory              int     `json:"memory"`
	MemoryDisplay       string  `json:"memoryDisplay"`
	MemoryPercentile    float64 `json:"memoryPercentile"`
	MemoryDistribution  string  `json:"memoryDistribution"`
	Code                string  `json:"code"`
	Timestamp           int     `json:"timestamp"`
	StatusCode          int     `json:"statusCode"`
	User                *struct {
		Username string `json:"username"`
		Profile  *struct {
			RealName   string `json:"realName"`
			UserAvatar string `json:"userAvatar"`
		} `json:"profile"`
	} `json:"user"`
	Lang *struct {
		Name        string `json:"name"`
		VerboseName string `json:"verboseName"`
	} `json:"lang"`
	Question *struct {
		QuestionId string `json:"questionId"`
	} `json:"question"`
	Notes string `json:"notes"`

	LastTestcase string `json:"lastTestcase"`
}

type QuestionContent struct {
	Content string
}
