package texts

import (
	"github.com/aid297/aid/v3/texts/buffer"
	"github.com/aid297/aid/v3/texts/html"
	"github.com/aid297/aid/v3/texts/markdown"
	"github.com/aid297/aid/v3/texts/rand"
	"github.com/aid297/aid/v3/texts/template"
	"github.com/aid297/aid/v3/texts/timer"
	"github.com/aid297/aid/v3/texts/transfer"
	"github.com/aid297/aid/v3/texts/volumer"
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
