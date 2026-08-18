package texts

import (
	"github.com/aid297/aid/v2/texts/buffer"
	"github.com/aid297/aid/v2/texts/html"
	"github.com/aid297/aid/v2/texts/markdown"
	"github.com/aid297/aid/v2/texts/rand"
	"github.com/aid297/aid/v2/texts/template"
	"github.com/aid297/aid/v2/texts/timer"
	"github.com/aid297/aid/v2/texts/transfer"
	"github.com/aid297/aid/v2/texts/volumer"
)

var (
	Buffer   buffer.Buffer           = (*buffer.BufferImpl)(nil)
	HTML     html.HTMLWriter         = (*html.HTMLWriterImpl)(nil)
	Markdown markdown.MarkdownWriter = (*markdown.MarkdownWriterImpl)(nil)
	Template template.Templater      = (*template.TemplaterImpl)(nil)
	Transfer transfer.Transfer       = (*transfer.TransferImpl)(nil)
	Rand     rand.Random             = (*rand.RandomImpl)(nil)
	Timer    timer.Timer             = (*timer.Impl)(nil)
	Volumer  volumer.Volumer         = (*volumer.VolumerImpl)(nil)
)
