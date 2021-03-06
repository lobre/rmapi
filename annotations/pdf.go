package annotations

import (
	"fmt"
	"os"

	"github.com/jung-kurt/gofpdf"
	"github.com/juruen/rmapi/archive"
	"github.com/juruen/rmapi/encoding/rm"
)

const (
	RmX      = rm.Width
	RmY      = rm.Height
	A4X      = 210
	A4Y      = 297
	ratioA4X = float32(A4X) / float32(RmX)
	ratioA4Y = float32(A4Y) / float32(RmY)
)

type PdfGenerator struct {
	zipName        string
	outputFilePath string
	options        PdfGeneratorOptions
}

type PdfGeneratorOptions struct {
	AddPageNumbers bool
}

func CreatePdfGenerator(zipName, outputFilePath string, options PdfGeneratorOptions) PdfGenerator {
	return PdfGenerator{zipName: zipName, outputFilePath: outputFilePath, options: options}
}

func (p PdfGenerator) Generate() error {
	file, err := os.Open(p.zipName)
	if err != nil {
		return err
	}
	defer file.Close()

	zip := archive.NewZip()

	fi, err := file.Stat()
	if err != nil {
		return err
	}

	err = zip.Read(file, fi.Size())
	if err != nil {
		return err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.SetFont("Arial", "", 10)
	for i, page := range zip.Pages {
		if page.Data == nil {
			continue
		}

		pdf.AddPage()
		if p.options.AddPageNumbers {
			pdf.Cell(0, 0, fmt.Sprintf("%d", i+1))
		}

		for _, layer := range page.Data.Layers {
			for _, line := range layer.Lines {

				if len(line.Points) < 1 {
					continue
				}

				for i := 1; i < len(line.Points); i++ {
					s := line.Points[i-1]
					x1 := s.X * ratioA4X
					y1 := s.Y * ratioA4Y

					s = line.Points[i]
					x2 := s.X * ratioA4X
					y2 := s.Y * ratioA4Y

					pdf.Line(float64(x1), float64(y1), float64(x2), float64(y2))
				}
			}
		}
	}

	return pdf.OutputFileAndClose(p.outputFilePath)
}
