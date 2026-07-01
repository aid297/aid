package texts

import (
	`github.com/aid297/aid/v2/texts/buffer`
	`github.com/aid297/aid/v2/texts/html`
	`github.com/aid297/aid/v2/texts/markdown`
	`github.com/aid297/aid/v2/texts/template`
	`github.com/aid297/aid/v2/texts/transfer`
)

var (
	Buffer         buffer.BufferImpl
	HTMLWriter     html.HTML
	MarkdownWriter markdown.Markdown
	Templater      template.Template
	Transfer       transfer.TransferImpl
)
