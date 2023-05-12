package musescore

type MusescoreStruct struct {
	Context      string `json:"@context"`
	Type         string `json:"@type"`
	URL          string `json:"url"`
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Text         string `json:"text"`
	Keywords     string `json:"keywords"`
	Composer     struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"composer"`
	Arranger struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"arranger"`
	MusicalKey    string `json:"musicalKey"`
	DatePublished string `json:"datePublished"`
	DateModified  string `json:"dateModified"`
	CommentCount  string `json:"commentCount"`
	DiscussionURL string `json:"discussionUrl"`
}
