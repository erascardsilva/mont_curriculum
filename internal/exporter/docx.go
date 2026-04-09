package exporter

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"mont_curriculum/internal/models"
	"mont_curriculum/internal/parser"
	"strings"
)

// GenerateDOCX cria um arquivo Word (.docx) compatível com o Microsoft Word a partir dos dados do currículo.
func GenerateDOCX(profile models.Profile, edus []models.Education, exps []models.Experience, projs []models.Project) ([]byte, error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// 1. [Content_Types].xml
	err := addFile(w, "[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
	<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
	<Default Extension="xml" ContentType="application/xml"/>
	<Default Extension="png" ContentType="image/png"/>
	<Default Extension="jpeg" ContentType="image/jpeg"/>
	<Default Extension="jpg" ContentType="image/jpeg"/>
	<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`)
	if err != nil {
		return nil, err
	}

	// 2. _rels/.rels
	err = addFile(w, "_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
	<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`)
	if err != nil {
		return nil, err
	}

	// 3. word/_rels/document.xml.rels (Necessário para a imagem)
	var relsContent strings.Builder
	relsContent.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	relsContent.WriteString(`<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">`)
	if profile.Photo != "" {
		relsContent.WriteString(`<Relationship Id="rIdPhoto" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image" Target="media/photo.img"/>`)
	}
	relsContent.WriteString(`</Relationships>`)
	err = addFile(w, "word/_rels/document.xml.rels", relsContent.String())
	if err != nil {
		return nil, err
	}

	// 4. Mídia (A foto real) e Cálculo de proporção
	imgWidth, imgHeight := 1440000, 1440000 // Default square EMUs
	if profile.Photo != "" {
		parts := strings.Split(profile.Photo, ",")
		if len(parts) > 1 {
			imageData, err := base64.StdEncoding.DecodeString(parts[1])
			if err == nil {
				imgFile, _ := w.Create("word/media/photo.img")
				imgFile.Write(imageData)

				// Detecta dimensões reais para evitar achatamento
				m, _, err := image.Decode(bytes.NewReader(imageData))
				if err == nil {
					bounds := m.Bounds()
					w_rect := bounds.Dx()
					h_rect := bounds.Dy()
					
					// Base: 1.5 polegadas (1371600 EMUs)
					baseSize := 1371600
					if w_rect > h_rect {
						imgWidth = baseSize
						imgHeight = int(float64(baseSize) * (float64(h_rect) / float64(w_rect)))
					} else {
						imgHeight = baseSize
						imgWidth = int(float64(baseSize) * (float64(w_rect) / float64(h_rect)))
					}
				}
			}
		}
	}

	// 5. word/document.xml
	docContent := buildDocumentXML(profile, edus, exps, projs, imgWidth, imgHeight)
	err = addFile(w, "word/document.xml", docContent)
	if err != nil {
		return nil, err
	}

	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func addFile(w *zip.Writer, name string, content string) error {
	f, err := w.Create(name)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(content))
	return err
}

func buildDocumentXML(profile models.Profile, edus []models.Education, exps []models.Experience, projs []models.Project, imgW, imgH int) string {
	var body bytes.Buffer

	// Se houver foto, adiciona o bloco de desenho (Drawing) com proporção correta
	if profile.Photo != "" {
		body.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:drawing><wp:inline distT="0" distB="0" distL="0" distR="0">
			<wp:extent cx="%d" cy="%d"/><wp:docPr id="1" name="Photo"/><wp:cNvGraphicFramePr><a:graphicFrameLocks xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" noChangeAspect="1"/></wp:cNvGraphicFramePr>
			<a:graphic xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"><a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture">
			<pic:pic xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture"><pic:nvPicPr><pic:cNvPr id="0" name="photo.img"/><pic:cNvPicPr/></pic:nvPicPr>
			<pic:blipFill><a:blip r:embed="rIdPhoto" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill>
			<pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr></pic:pic>
			</a:graphicData></a:graphic></wp:inline></w:drawing></w:r></w:p>`, imgW, imgH, imgW, imgH))
	}

	// Cabeçalho
	body.WriteString(`<w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="48"/></w:rPr><w:t>`)
	body.WriteString(profile.FirstName + " " + profile.LastName)
	body.WriteString(`</w:t></w:r></w:p>`)

	// Contatos e Redes Sociais
	contactInfo := profile.Email
	if profile.Phone != "" { contactInfo += " | " + profile.Phone }
	body.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:t>%s</w:t></w:r></w:p>`, contactInfo))
	
	if profile.Address != "" {
		body.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:t>%s</w:t></w:r></w:p>`, profile.Address))
	}

	// Links Sociais - Parágrafos individuais para melhor leitura
	if profile.LinkedIn != "" {
		body.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/><w:color w:val="38BDF8"/></w:rPr><w:t>LinkedIn: </w:t></w:r><w:r><w:rPr><w:color w:val="38BDF8"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, profile.LinkedIn))
	}
	if profile.GitHub != "" {
		body.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/><w:color w:val="38BDF8"/></w:rPr><w:t>GitHub: </w:t></w:r><w:r><w:rPr><w:color w:val="38BDF8"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, profile.GitHub))
	}
	if profile.Website != "" {
		body.WriteString(fmt.Sprintf(`<w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/><w:color w:val="38BDF8"/></w:rPr><w:t>Portfólio: </w:t></w:r><w:r><w:rPr><w:color w:val="38BDF8"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, profile.Website))
	}

	body.WriteString(`<w:p><w:pPr><w:pBdr><w:bottom w:val="single" w:sz="6" w:space="1" w:color="auto"/></w:pBdr></w:pPr></w:p>`)

	if profile.Objective != "" {
		body.WriteString(addSectionHeader("RESUMO PROFISSIONAL"))
		body.WriteString(fmt.Sprintf(`<w:p><w:r><w:t>%s</w:t></w:r></w:p>`, parser.Sanitize(profile.Objective)))
	}

	if len(exps) > 0 {
		body.WriteString(addSectionHeader("EXPERIÊNCIA PROFISSIONAL"))
		for _, exp := range exps {
			body.WriteString(fmt.Sprintf(`<w:p><w:r><w:rPr><w:b/></w:rPr><w:t>%s @ %s</w:t></w:r></w:p>`, exp.Position, exp.Company))
			body.WriteString(fmt.Sprintf(`<w:p><w:r><w:rPr><w:i/></w:rPr><w:t>%s - %s</w:t></w:r></w:p>`, exp.StartDate, exp.EndDate))
			if exp.Description != "" {
				body.WriteString(fmt.Sprintf(`<w:p><w:r><w:t>%s</w:t></w:r></w:p>`, parser.Sanitize(exp.Description)))
			}
			body.WriteString(`<w:p/>`)
		}
	}

	if len(edus) > 0 {
		body.WriteString(addSectionHeader("FORMAÇÃO ACADÊMICA"))
		for _, edu := range edus {
			body.WriteString(fmt.Sprintf(`<w:p><w:r><w:rPr><w:b/></w:rPr><w:t>%s</w:t></w:r></w:p>`, edu.Course))
			body.WriteString(fmt.Sprintf(`<w:p><w:r><w:t>%s (%s - %s)</w:t></w:r></w:p>`, edu.Institution, edu.StartDate, edu.EndDate))
		}
	}

	if len(projs) > 0 {
		body.WriteString(addSectionHeader("PROJETOS E PORTFÓLIO"))
		for _, proj := range projs {
			body.WriteString(fmt.Sprintf(`<w:p><w:r><w:rPr><w:b/></w:rPr><w:t>%s</w:t></w:r></w:p>`, proj.Name))
			if proj.URL != "" {
				body.WriteString(fmt.Sprintf(`<w:p><w:r><w:rPr><w:color w:val="38BDF8"/><w:u w:val="single"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, proj.URL))
			}
			if proj.Description != "" {
				body.WriteString(fmt.Sprintf(`<w:p><w:r><w:t>%s</w:t></w:r></w:p>`, proj.Description))
			}
			body.WriteString(`<w:p/>`)
		}
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" 
            xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships" 
            xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing" 
            xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main" 
            xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">
	<w:body>
		%s
	</w:body>
</w:document>`, body.String())
}

func addSectionHeader(title string) string {
	return fmt.Sprintf(`<w:p><w:pPr><w:spacing w:before="240" w:after="120"/><w:pBdr><w:bottom w:val="single" w:sz="4" w:space="1" w:color="E2E8F0"/></w:pBdr></w:pPr><w:r><w:rPr><w:b/><w:sz w:val="28"/><w:color w:val="38BDF8"/></w:rPr><w:t>%s</w:t></w:r></w:p>`, title)
}
