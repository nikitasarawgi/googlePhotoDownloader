package main

type GoogleImageRequestDetails struct {
	Uri                string `default: "https://www.googleapis.com/customsearch/v1?"`
	ImageSearchKeyword string
	GoogleSearchKey    string
	SearchEngineId     string
}

type GoogleImage struct {
	ContextLink     string
	Height          int
	Width           int
	ByteSize        int
	ThumbnailLink   string
	ThumbnailHeight int
	ThumbnailWidth  int
}

type ImageItem struct {
	Kind        string
	Title       string
	HtmlTitle   string
	Link        string
	DisplayLink string
	Snippet     string
	HtmlSnippet string
	Mime        string
	FileFormat  string
	Image       GoogleImage
}

type ResponsePageDetails struct {
	Title          string
	TotalResults   string
	SearchTerms    string
	Count          int
	StartIndex     int
	InputEncoding  string
	OutputEncoding string
	Safe           string
	Cx             string
	SearchType     string
}

type GoogleImageSearchResponse struct {
	Kind        string
	ResponseUrl struct {
		UrlType  string `json:"url"`
		Template string `json:"template"`
	} `json:"url"`
	Queries struct {
		PreviousPage []ResponsePageDetails `json:"previousPage,omitempty"`
		Request      []ResponsePageDetails `json:"request"`
		NextPage     []ResponsePageDetails `json:"nextPage"`
	} `json:"queries"`
	Context struct {
		Title string
	}
	SearchInformation struct {
		SearchTime            float64
		FormattedSearchTime   string
		TotalResults          string
		FormattedTotalResults string
	}
	Spelling struct {
		CorrectedQuery     string
		HtmlCorrectedQuery string
	}
	Items []ImageItem
}
