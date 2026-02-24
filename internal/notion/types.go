package notion

import "time"

// Common types

type RichText struct {
	Type        string       `json:"type"`
	PlainText   string       `json:"plain_text"`
	Href        string       `json:"href,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`
}

type Annotations struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

type Icon struct {
	Type  string `json:"type"`
	Emoji string `json:"emoji,omitempty"`
}

// Search types

type SearchRequest struct {
	Query string       `json:"query"`
	Sort  *SearchSort  `json:"sort,omitempty"`
	Filter *SearchFilter `json:"filter,omitempty"`
	StartCursor string `json:"start_cursor,omitempty"`
	PageSize    int    `json:"page_size,omitempty"`
}

type SearchSort struct {
	Direction string `json:"direction"`
	Timestamp string `json:"timestamp"`
}

type SearchFilter struct {
	Value    string `json:"value"`
	Property string `json:"property"`
}

type SearchResponse struct {
	Results    []SearchResult `json:"results"`
	HasMore    bool           `json:"has_more"`
	NextCursor string         `json:"next_cursor"`
}

type SearchResult struct {
	ID             string     `json:"id"`
	Object         string     `json:"object"` // "page" or "database"
	Icon           *Icon      `json:"icon,omitempty"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	Properties     Properties `json:"properties"`
	URL            string     `json:"url"`
}

type Properties map[string]Property

type Property struct {
	ID    string     `json:"id"`
	Type  string     `json:"type"`
	Title []RichText `json:"title,omitempty"`
}

// Page types

type Page struct {
	ID             string     `json:"id"`
	Object         string     `json:"object"`
	Icon           *Icon      `json:"icon,omitempty"`
	LastEditedTime time.Time  `json:"last_edited_time"`
	Properties     Properties `json:"properties"`
}

// Database types

type Database struct {
	ID             string     `json:"id"`
	Object         string     `json:"object"`
	Icon           *Icon      `json:"icon,omitempty"`
	Title          []RichText `json:"title"`
	LastEditedTime time.Time  `json:"last_edited_time"`
}

// Block types

type BlocksResponse struct {
	Results    []Block `json:"results"`
	HasMore    bool    `json:"has_more"`
	NextCursor string  `json:"next_cursor"`
}

type Block struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	HasChildren bool   `json:"has_children"`

	Paragraph        *ParagraphBlock   `json:"paragraph,omitempty"`
	Heading1         *HeadingBlock     `json:"heading_1,omitempty"`
	Heading2         *HeadingBlock     `json:"heading_2,omitempty"`
	Heading3         *HeadingBlock     `json:"heading_3,omitempty"`
	BulletedListItem *ListItemBlock    `json:"bulleted_list_item,omitempty"`
	NumberedListItem *ListItemBlock    `json:"numbered_list_item,omitempty"`
	Toggle           *ListItemBlock    `json:"toggle,omitempty"`
	ToDo             *ToDoBlock        `json:"to_do,omitempty"`
	Code             *CodeBlock        `json:"code,omitempty"`
	Quote            *QuoteBlock       `json:"quote,omitempty"`
	Callout          *CalloutBlock     `json:"callout,omitempty"`
	Divider          *struct{}         `json:"divider,omitempty"`
	Image            *FileBlock        `json:"image,omitempty"`
	Bookmark         *BookmarkBlock    `json:"bookmark,omitempty"`
	ChildPage        *ChildPageBlock   `json:"child_page,omitempty"`
	ChildDatabase    *ChildDBBlock     `json:"child_database,omitempty"`
	TableOfContents  *struct{}         `json:"table_of_contents,omitempty"`
}

type ParagraphBlock struct {
	RichText []RichText `json:"rich_text"`
}

type HeadingBlock struct {
	RichText []RichText `json:"rich_text"`
}

type ListItemBlock struct {
	RichText []RichText `json:"rich_text"`
}

type ToDoBlock struct {
	RichText []RichText `json:"rich_text"`
	Checked  bool       `json:"checked"`
}

type CodeBlock struct {
	RichText []RichText `json:"rich_text"`
	Language string     `json:"language"`
}

type QuoteBlock struct {
	RichText []RichText `json:"rich_text"`
}

type CalloutBlock struct {
	RichText []RichText `json:"rich_text"`
	Icon     *Icon      `json:"icon,omitempty"`
}

type FileBlock struct {
	Type     string    `json:"type"`
	File     *FileURL  `json:"file,omitempty"`
	External *FileURL  `json:"external,omitempty"`
	Caption  []RichText `json:"caption,omitempty"`
}

type FileURL struct {
	URL string `json:"url"`
}

type BookmarkBlock struct {
	URL     string     `json:"url"`
	Caption []RichText `json:"caption,omitempty"`
}

type ChildPageBlock struct {
	Title string `json:"title"`
}

type ChildDBBlock struct {
	Title string `json:"title"`
}
