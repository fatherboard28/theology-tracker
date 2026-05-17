package models

type TopicType string

const (
	TopicTypeBookOfBible      TopicType = "Book of Bible"
	TopicTypeTheologicalTheme TopicType = "Theological Theme"
	TopicTypeDoctrine         TopicType = "Doctrine"
	TopicTypeOther            TopicType = "Other"
)

type Topic struct {
	ID            int64     `db:"id"`
	Title         string    `db:"title"`
	Description   string    `db:"description"`
	Type          TopicType `db:"type"`
	ParentTopicID *int64    `db:"parent_topic_id"` // nil if no parent
	CreatedAt     string    `db:"created_at"`
}

// TopicSummary is a lightweight view used in tag inputs and listings
// where the full topic struct isn't needed.
type TopicSummary struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
	Type  string `db:"type"`
}
