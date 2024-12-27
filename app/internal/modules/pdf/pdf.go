package pdf

import (
	"IMP/app/internal/modules/pdf/mappers"
	"github.com/go-pdf/fpdf"
	"os"
)

type Builder struct {
	pdf         fpdf.Fpdf
	innerFolder *string
	fileName    string
}

// NewBuilder creates a new PDF builder based on fpdf file with default settings and font.
func NewBuilder(fileName string, innerFolder *string) *Builder {
	pdf := &Builder{
		pdf:         *fpdf.New("P", "mm", "A4", ""),
		innerFolder: innerFolder,
		fileName:    fileName,
	}
	pdf.pdf.AddPage()
	pdf.pdf.SetFont("Arial", "", 12)

	return pdf
}

// PrintLn prints a line of text to the PDF end of file.
func (p *Builder) PrintLn(text string) {
	p.pdf.Cell(40, 10, text)
	p.pdf.Ln(-1)
}

// PrintTable prints a table to the PDF end of file.
func (p *Builder) PrintTable(data []mappers.TableMapper) {
	for _, header := range data[0].Headers() {
		p.pdf.Cell(40, 10, header)
	}
	p.pdf.Ln(-1)

	// Write data
	for _, row := range data {
		for _, col := range row.ToTable() {
			p.pdf.Cell(40, 10, col)
		}
		p.pdf.Ln(-1)
	}
}

// Save saves the PDF file to the specified file name.
func (p *Builder) Save() error {
	path := "outputs/"
	if p.innerFolder != nil {
		if err := os.MkdirAll(path+*p.innerFolder, 0755); err != nil {
			return err
		}

		path += *p.innerFolder + "/"
	}
	err := p.pdf.OutputFileAndClose(path + p.fileName + ".pdf")

	return err
}
