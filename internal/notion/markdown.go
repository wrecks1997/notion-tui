package notion

import (
	"fmt"
	"strings"
)

// BlocksToMarkdown converts a slice of Notion blocks to a markdown string.
func BlocksToMarkdown(blocks []Block) string {
	var sb strings.Builder
	numberedIdx := 0

	for _, block := range blocks {
		// Track numbered list index
		if block.Type == "numbered_list_item" {
			numberedIdx++
		} else {
			numberedIdx = 0
		}

		line := blockToMarkdown(block, numberedIdx)
		if line != "" {
			sb.WriteString(line)
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

func blockToMarkdown(block Block, numberedIdx int) string {
	switch block.Type {
	case "paragraph":
		if block.Paragraph == nil {
			return ""
		}
		text := richTextToMarkdown(block.Paragraph.RichText)
		return text + "\n"

	case "heading_1":
		if block.Heading1 == nil {
			return ""
		}
		return "# " + richTextToMarkdown(block.Heading1.RichText) + "\n"

	case "heading_2":
		if block.Heading2 == nil {
			return ""
		}
		return "## " + richTextToMarkdown(block.Heading2.RichText) + "\n"

	case "heading_3":
		if block.Heading3 == nil {
			return ""
		}
		return "### " + richTextToMarkdown(block.Heading3.RichText) + "\n"

	case "bulleted_list_item":
		if block.BulletedListItem == nil {
			return ""
		}
		return "- " + richTextToMarkdown(block.BulletedListItem.RichText)

	case "numbered_list_item":
		if block.NumberedListItem == nil {
			return ""
		}
		return fmt.Sprintf("%d. %s", numberedIdx, richTextToMarkdown(block.NumberedListItem.RichText))

	case "to_do":
		if block.ToDo == nil {
			return ""
		}
		check := " "
		if block.ToDo.Checked {
			check = "x"
		}
		return fmt.Sprintf("- [%s] %s", check, richTextToMarkdown(block.ToDo.RichText))

	case "toggle":
		if block.Toggle == nil {
			return ""
		}
		return "▸ " + richTextToMarkdown(block.Toggle.RichText)

	case "code":
		if block.Code == nil {
			return ""
		}
		lang := block.Code.Language
		if lang == "plain text" {
			lang = ""
		}
		return fmt.Sprintf("```%s\n%s\n```\n", lang, richTextToPlain(block.Code.RichText))

	case "quote":
		if block.Quote == nil {
			return ""
		}
		text := richTextToMarkdown(block.Quote.RichText)
		lines := strings.Split(text, "\n")
		var quoted []string
		for _, line := range lines {
			quoted = append(quoted, "> "+line)
		}
		return strings.Join(quoted, "\n") + "\n"

	case "callout":
		if block.Callout == nil {
			return ""
		}
		icon := ""
		if block.Callout.Icon != nil && block.Callout.Icon.Emoji != "" {
			icon = block.Callout.Icon.Emoji + " "
		}
		return "> " + icon + richTextToMarkdown(block.Callout.RichText) + "\n"

	case "divider":
		return "---\n"

	case "image":
		if block.Image == nil {
			return ""
		}
		url := ""
		if block.Image.File != nil {
			url = block.Image.File.URL
		} else if block.Image.External != nil {
			url = block.Image.External.URL
		}
		caption := richTextToPlain(block.Image.Caption)
		if caption == "" {
			caption = "image"
		}
		return fmt.Sprintf("![%s](%s)\n", caption, url)

	case "bookmark":
		if block.Bookmark == nil {
			return ""
		}
		caption := richTextToPlain(block.Bookmark.Caption)
		if caption == "" {
			caption = block.Bookmark.URL
		}
		return fmt.Sprintf("[%s](%s)\n", caption, block.Bookmark.URL)

	case "child_page":
		if block.ChildPage == nil {
			return ""
		}
		return fmt.Sprintf("📄 **%s**\n", block.ChildPage.Title)

	case "child_database":
		if block.ChildDatabase == nil {
			return ""
		}
		return fmt.Sprintf("🗃 **%s**\n", block.ChildDatabase.Title)

	case "table_of_contents":
		return "_Table of Contents_\n"

	default:
		return ""
	}
}

func richTextToMarkdown(texts []RichText) string {
	var parts []string
	for _, t := range texts {
		text := t.PlainText
		if t.Annotations != nil {
			if t.Annotations.Code {
				text = "`" + text + "`"
			}
			if t.Annotations.Bold {
				text = "**" + text + "**"
			}
			if t.Annotations.Italic {
				text = "_" + text + "_"
			}
			if t.Annotations.Strikethrough {
				text = "~~" + text + "~~"
			}
		}
		if t.Href != "" {
			text = fmt.Sprintf("[%s](%s)", text, t.Href)
		}
		parts = append(parts, text)
	}
	return strings.Join(parts, "")
}

func richTextToPlain(texts []RichText) string {
	var parts []string
	for _, t := range texts {
		parts = append(parts, t.PlainText)
	}
	return strings.Join(parts, "")
}
