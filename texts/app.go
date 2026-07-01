package texts

import (
	`github.com/aid297/aid/v2/texts/buffer`
	`github.com/aid297/aid/v2/texts/html`
	`github.com/aid297/aid/v2/texts/markdown`
	`github.com/aid297/aid/v2/texts/template`
	`github.com/aid297/aid/v2/texts/transfer`
)

var (
	Buffer   buffer.Buffer
	HTML     html.HTMLWriter
	Markdown markdown.MarkdownWriter
	Template template.Templater
	Transfer transfer.Transfer
)
