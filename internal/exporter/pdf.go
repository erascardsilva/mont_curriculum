package exporter

/*
 * Project: Monte Curriculum
 * Author: Erasmo Cardoso - Software Engineer | Electronics Techniciant
 */

import (
	"fmt"
	"mont_curriculum/internal/models"
	"mont_curriculum/internal/parser"
	"strings"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

// GeneratePDF cria o arquivo PDF do currículo bando no modelo selecionado
func GeneratePDF(profile models.Profile, edus []models.Education, exps []models.Experience, projs []models.Project, template string) ([]byte, error) {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(15, 12, 15) // Margens aumentadas

	// Configuração de cores baseado no template
	primaryColor := color.Color{Red: 56, Green: 189, Blue: 248} // Default Blue
	if template == "professional" {
		primaryColor = color.Color{Red: 100, Green: 116, Blue: 139} // Slate
	} else if template == "vibrant" {
		primaryColor = color.Color{Red: 251, Green: 146, Blue: 60} // Orange
	}

	// Header com Foto e Nome
	m.Row(45, func() {
		m.Col(3, func() {
			if profile.Photo != "" {
				parts := strings.Split(profile.Photo, ",")
				if len(parts) > 1 {
					ext := consts.Jpg
					if strings.Contains(parts[0], "png") { ext = consts.Png }
					m.Base64Image(parts[1], ext, props.Rect{
						Center:  true,
						Percent: 95,
					})
				}
			}
		})
		m.Col(9, func() {
			m.Text(strings.ToUpper(fmt.Sprintf("%s %s", profile.FirstName, profile.LastName)), props.Text{
				Size:  22,
				Style: consts.Bold,
				Align: consts.Left,
				Top:   4,
			})
			
			contactLine := profile.Email
			if profile.Phone != "" { contactLine += "  |  " + profile.Phone }
			m.Text(contactLine, props.Text{
				Size: 10,
				Top:  16,
			})
			
			m.Text(profile.Address, props.Text{
				Size: 10,
				Top:  22,
			})

			// Links Sociais - Refatorado para melhor separação
			socialTop := float64(28)
			if profile.LinkedIn != "" {
				m.Text("LinkedIn: ", props.Text{Size: 8, Style: consts.Bold, Top: socialTop})
				m.Text(profile.LinkedIn, props.Text{Size: 8, Color: primaryColor, Top: socialTop, Left: 15})
				socialTop += 4
			}
			if profile.GitHub != "" {
				m.Text("GitHub: ", props.Text{Size: 8, Style: consts.Bold, Top: socialTop})
				m.Text(profile.GitHub, props.Text{Size: 8, Color: primaryColor, Top: socialTop, Left: 15})
				socialTop += 4
			}
			if profile.Website != "" {
				m.Text("Portfólio: ", props.Text{Size: 8, Style: consts.Bold, Top: socialTop})
				m.Text(profile.Website, props.Text{Size: 8, Color: primaryColor, Top: socialTop, Left: 15})
			}
		})
	})

	m.Line(1)

	// Objetivo / Resumo Profissional
	if profile.Objective != "" {
		addSectionTitle(m, "RESUMO PROFISSIONAL", primaryColor)
		obj := parser.Sanitize(profile.Objective)
		addRichText(m, obj, 10, color.Color{Red: 30, Green: 41, Blue: 59})
	}

	// Experiência (Ordem: Empresa -> Cargo -> Data)
	if len(exps) > 0 {
		addSectionTitle(m, "EXPERIÊNCIA PROFISSIONAL", primaryColor)
		for _, exp := range exps {
			// 1. Empresa (Bold)
			m.Row(7, func() {
				m.Col(12, func() {
					m.Text(strings.ToUpper(exp.Company), props.Text{
						Size:  11,
						Style: consts.Bold,
					})
				})
			})
			// 2. Cargo
			m.Row(6, func() {
				m.Col(12, func() {
					m.Text(exp.Position, props.Text{
						Size: 10,
						Top:  1,
					})
				})
			})
			// 3. Datas (Italic / Small)
			m.Row(5, func() {
				m.Col(12, func() {
					dateRange := exp.StartDate
					if exp.EndDate != "" { dateRange += " — " + exp.EndDate }
					m.Text(dateRange, props.Text{
						Size:  9,
						Style: consts.Italic,
						Color: color.Color{Red: 100, Green: 116, Blue: 139},
					})
				})
			})
			// 4. Descrição
			if exp.Description != "" && exp.Description != "-" {
				cleanDesc := parser.Sanitize(exp.Description)
				addRichText(m, cleanDesc, 9.5, color.Color{Red: 30, Green: 41, Blue: 59})
			}
			m.Row(6, func() {}) // Espaço entre itens
		}
	}

	// Formação (Ordem: Instituição -> Curso -> Data)
	if len(edus) > 0 {
		addSectionTitle(m, "FORMAÇÃO ACADÊMICA", primaryColor)
		for _, edu := range edus {
			m.Row(7, func() {
				m.Col(12, func() {
					m.Text(strings.ToUpper(edu.Institution), props.Text{
						Size:  11,
						Style: consts.Bold,
					})
				})
			})
			m.Row(6, func() {
				m.Col(12, func() {
					m.Text(edu.Course, props.Text{
						Size: 10,
					})
				})
			})
			m.Row(5, func() {
				m.Col(12, func() {
					dateRange := edu.StartDate
					if edu.EndDate != "" { dateRange += " — " + edu.EndDate }
					m.Text(dateRange, props.Text{
						Size:  9,
						Style: consts.Italic,
						Color: color.Color{Red: 100, Green: 116, Blue: 139},
					})
				})
			})
			m.Row(5, func(){})
		}
	}

	// Projetos
	if len(projs) > 0 {
		addSectionTitle(m, "PROJETOS E PORTFÓLIO", primaryColor)
		for _, proj := range projs {
			m.Row(7, func() {
				m.Col(12, func() {
					m.Text(proj.Name, props.Text{
						Size:  11,
						Style: consts.Bold,
					})
				})
			})
			if proj.URL != "" {
				m.Row(6, func() {
					m.Col(12, func() {
						m.Text(proj.URL, props.Text{
							Size:  9,
							Color: primaryColor,
						})
					})
				})
			}
			if proj.Description != "" {
				cleanDesc := parser.Sanitize(proj.Description)
				addRichText(m, cleanDesc, 9.5, color.Color{Red: 30, Green: 41, Blue: 59})
			}
			m.Row(5, func() {})
		}
	}

	buffer, err := m.Output()
	if err != nil { return nil, err }
	return buffer.Bytes(), nil
}

func addSectionTitle(m pdf.Maroto, title string, clr color.Color) {
	m.Row(14, func() {
		m.Col(12, func() {
			m.Text(title, props.Text{
				Size:  12,
				Style: consts.Bold,
				Color: clr,
				Top:   8,
			})
		})
	})
	m.Line(1)
	m.Row(8, func(){}) // Espaçamento aumentado para evitar encavalamento
}

// addRichText renderiza texto processando tópicos (bullets) e garantindo espaçamento compacto
func addRichText(m pdf.Maroto, text string, textSize float64, clr color.Color) {
	// Se não houver marcadores de lista, renderiza como bloco único para evitar gaps
	if !strings.Contains(text, "- ") && !strings.Contains(text, "* ") {
		content := strings.ReplaceAll(text, "\n", " ")
		charLimit := 85
		rowsNeeded := float64(len(content)/charLimit + 1)
		rowHeight := rowsNeeded * 4.6 // Aumento sutil no line-height
		
		m.Row(rowHeight, func() {
			m.Col(12, func() {
				m.Text(content, props.Text{
					Size:  textSize,
					Align: consts.Left,
					Color: clr,
					Top:   2, // Padding vertical interno
				})
			})
		})
		return
	}

	// Se houver lista, processa linha a linha de forma compacta
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" { continue }

		isBullet := strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*")
		content := line
		if isBullet {
			content = strings.TrimSpace(line[1:])
		}

		charLimit := 85
		rowsNeeded := float64(len(content)/charLimit + 1)
		rowHeight := rowsNeeded * 4.4 // Mais compacto para listas

		m.Row(rowHeight, func() {
			if isBullet {
				m.Col(1, func() {
					m.Text("•", props.Text{
						Size:  textSize + 1,
						Style: consts.Bold,
						Align: consts.Right,
						Color: clr,
						Top:   1,
					})
				})
				m.Col(11, func() {
					m.Text(content, props.Text{
						Size:  textSize,
						Align: consts.Left,
						Color: clr,
						Top:   1,
					})
				})
			} else {
				m.Col(12, func() {
					m.Text(content, props.Text{
						Size:  textSize,
						Align: consts.Left,
						Color: clr,
					})
				})
			}
		})
	}
}
