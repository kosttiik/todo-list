package ds

const (
	DateFormat = "20060102"
)

type Task struct {
	ID      int    `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}
