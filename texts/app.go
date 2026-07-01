package texts

import (
	`github.com/aid297/aid/v2/texts/buffer`
	`github.com/aid297/aid/v2/texts/html`
	`github.com/aid297/aid/v2/texts/markdown`
	`github.com/aid297/aid/v2/texts/template`
	`github.com/aid297/aid/v2/texts/transfer`
)

var (
	Buffer         buffer.Buffer           = (*buffer.BufferImpl)(nil)
	HTMLWriter     html.HTMLWriter         = (*html.HTML)(nil)
	MarkdownWriter markdown.MarkdownWriter = (*markdown.Markdown)(nil)
	Templater      template.Templater      = (*template.Template)(nil)
	Transfer       transfer.Transfer       = (*transfer.TransferImpl)(nil)
)
