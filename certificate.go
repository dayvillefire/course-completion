package main

import (
	"log"

	"github.com/go-pdf/fpdf"
	"github.com/go-pdf/fpdf/contrib/gofpdi"
)

type generator struct {
	Template           string
	GlobalReplacements map[string]string
	Replacements       map[string]Replacement
}

func (g *generator) generate(fn string, repl map[string]string) error {
	log.Printf("generate")
	pdf := fpdf.New("L", "pt", "Letter", "")

	imp := gofpdi.NewImporter()
	tpl := imp.ImportPage(pdf, g.Template, 1, "/MediaBox")
	pageSizes := imp.GetPageSizes()

	// TODO: FIXME: Specify in configuration
	pdf.SetFontLocation("fonts")

	for _, r := range g.Replacements {
		log.Printf("Add font %s : %s", r.FontFamily, r.FontJson)
		pdf.AddFont(r.FontFamily, "", r.FontJson)
	}

	log.Printf("Add page")
	pdf.AddPage()

	log.Printf("Use imported template")
	imp.UseImportedTemplate(
		pdf,
		tpl,
		0, 0,
		pageSizes[0]["/MediaBox"]["w"], pageSizes[0]["/MediaBox"]["h"],
	)

	log.Printf("Iterating through replacements")
	for k, v := range g.Replacements {
		log.Printf("Replacement %s : %#v", k, v)
		text, ok := g.GlobalReplacements[k]
		if !ok {
			text, ok = repl[k]
			if !ok {
				log.Printf("WARN: no value for %s, skipping", k)
			}
		}
		pdf.SetFont(v.FontFamily, "", v.FontSize)
		wd := pdf.GetStringWidth(text) + 6
		pdf.SetX((pageSizes[0]["/MediaBox"]["w"] - wd) / 2)
		pdf.SetY(v.NameY)
		pdf.WriteAligned(0, v.FontSize, text, "C")
	}

	return pdf.OutputFileAndClose(fn)
}
