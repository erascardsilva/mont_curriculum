package parser

import (
	"bytes"
	"fmt"
	"mont_curriculum/internal/models"
	"regexp"
	"strings"
	"unicode"

	"github.com/ledongthuc/pdf"
)

// ParseResume analisa um arquivo e tenta extrair informações completas
func ParseResume(filePath string) (*ExtractedData, error) {
	var text string
	var err error

	if strings.HasSuffix(filePath, ".pdf") {
		text, err = readPdf(filePath)
	} else {
		return nil, fmt.Errorf("formato não suportado: %s", filePath)
	}

	if err != nil {
		return nil, err
	}

	return extractDeepHeuristics(text), nil
}

func readPdf(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var buf bytes.Buffer
	b, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(b)
	
	cleanText := buf.String()
	// Tenta consertar PDFs com espaços intercalados (E r a s m o)
	if isExtremelySpaced(cleanText) {
		cleanText = fixSpacedOutText(cleanText)
	}
	
	return cleanText, nil
}

func fixSpacedOutText(text string) string {
	var result strings.Builder
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		// Heurística aprimorada: Só limpa se houver um padrão redundante de letra-espaço-letra
		// sem espaços maiores que sugeririam separação de palavras real.
		if isExtremelySpaced(line) {
			tmp := regexp.MustCompile(`\s{2,}`).ReplaceAllString(line, "|||")
			tmp = strings.ReplaceAll(tmp, " ", "")
			line = strings.ReplaceAll(tmp, "|||", " ")
		}
		result.WriteString(line + "\n")
	}
	return result.String()
}

func isExtremelySpaced(text string) bool {
	// Letra isolada repetida 4 vezes ou mais: "E r a s"
	re := regexp.MustCompile(`^([A-Z]\s){4,}`)
	return re.MatchString(strings.TrimSpace(text))
}

func normalize(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	tbl := map[string]string{
		"á": "a", "ã": "a", "â": "a", "à": "a",
		"é": "e", "ê": "e",
		"í": "i",
		"ó": "o", "õ": "o", "ô": "o",
		"ú": "u",
		"ç": "c",
	}
	for k, v := range tbl {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}

func extractDeepHeuristics(text string) *ExtractedData {
	data := &ExtractedData{
		Profile: models.Profile{ID: 1},
	}
	lines := strings.Split(text, "\n")

	// Identificação de Seções
	sections := map[string][]string{}
	currentSection := "header"
	
	keywords := map[string]string{
		"educacao": "education", "formacao": "education", "academico": "education", "qualificacoes": "education",
		"experiencia": "experience", "profissional": "experience", "carreira": "experience", "historico": "experience",
		"projetos": "projects", "portfolio": "projects", "github": "projects",
		"sobre": "objective", "resumo": "objective", "objetivo": "objective", "perfil": "objective",
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" { continue }
		
		norm := normalize(trimmed)
		foundHeader := false
		
		for kw, sec := range keywords {
			// Se a linha for exatamente ou contiver a palavra-chave de forma isolada
			if norm == kw || (len(trimmed) < 40 && strings.Contains(norm, kw)) {
				if len(trimmed) < 50 {
					currentSection = sec
					foundHeader = true
					break
				}
			}
		}
		
		if !foundHeader {
			sections[currentSection] = append(sections[currentSection], trimmed)
		}
	}

	extractProfile(data, sections["header"], text)
	
	if lines, ok := sections["objective"]; ok {
		data.Profile.Objective = strings.Join(lines, " ")
	}

	extractEducation(data, sections["education"])
	extractExperience(data, sections["experience"])
	extractProjects(data, sections["projects"])

	// Sanitização final dos campos de texto longo
	data.Profile.Objective = Sanitize(data.Profile.Objective)
	for i := range data.Experience {
		data.Experience[i].Description = Sanitize(data.Experience[i].Description)
	}
	for i := range data.Education {
		data.Education[i].Description = Sanitize(data.Education[i].Description)
	}

	return data
}

func extractProfile(data *ExtractedData, headerLines []string, fullText string) {
	reEmail := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
	data.Profile.Email = reEmail.FindString(fullText)
	
	rePhone := regexp.MustCompile(`\(\d{2}\)\s\d{4,5}-\d{4}|\d{2}\s\d{4,5}-\d{4}|\d{4,5}-\d{4}`)
	data.Profile.Phone = rePhone.FindString(fullText)

	// Nome: Busca nas primeiras linhas do header
	for _, line := range headerLines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 5 && !strings.Contains(trimmed, "@") && !strings.Contains(trimmed, "http") {
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 && unicode.IsUpper(rune(parts[0][0])) {
				data.Profile.FirstName = parts[0]
				data.Profile.LastName = strings.Join(parts[1:], " ")
				break
			}
		}
	}
}

func extractEducation(data *ExtractedData, lines []string) {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) < 5 || trimmed == "-" { continue }
		
		edu := models.Education{ProfileID: 1}
		if parts := strings.FieldsFunc(trimmed, func(r rune) bool { return r == '-' || r == '—' || r == '|' }); len(parts) >= 2 {
			edu.Course = strings.TrimSpace(parts[0])
			edu.Institution = strings.TrimSpace(parts[1])
		} else {
			edu.Course = trimmed
		}
		data.Education = append(data.Education, edu)
	}
}

func extractExperience(data *ExtractedData, lines []string) {
	var currentExp *models.Experience
	reDate := regexp.MustCompile(`(\d{2}/\d{2}/\d{4}|\d{2}/\d{4}|\d{4})`)
	reFullDate := regexp.MustCompile(`(?i)(jan|fev|mar|abr|mai|jun|jul|ago|set|out|nov|dez|atual|presente|\d{4})`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed == "-" { continue }

		// Detecta se a linha tem cara de Nova Experiência (Datas ou Empresa @ Cargo)
		hasDate := reFullDate.MatchString(trimmed)
		hasBusiness := strings.Contains(trimmed, "@") || strings.Contains(trimmed, "-")

		if hasDate && (hasBusiness || len(trimmed) < 100) {
			if currentExp != nil { data.Experience = append(data.Experience, *currentExp) }
			currentExp = &models.Experience{ProfileID: 1}
			
			// Extração de datas
			dates := reDate.FindAllString(trimmed, -1)
			if len(dates) > 0 {
				currentExp.StartDate = dates[0]
				if len(dates) > 1 { currentExp.EndDate = dates[1] }
			}
			if strings.Contains(strings.ToLower(trimmed), "atual") || strings.Contains(strings.ToLower(trimmed), "presente") {
				currentExp.EndDate = "Atual"
			}

			// Extração de Cargo/Empresa
			if strings.Contains(trimmed, "@") {
				parts := strings.Split(trimmed, "@")
				currentExp.Position = strings.TrimSpace(reDate.ReplaceAllString(parts[0], ""))
				currentExp.Company = strings.TrimSpace(reDate.ReplaceAllString(parts[1], ""))
			} else {
				currentExp.Position = strings.TrimSpace(reDate.ReplaceAllString(trimmed, ""))
			}
		} else if currentExp != nil {
			if currentExp.Company == "" && !hasDate {
				currentExp.Company = trimmed
			} else {
				// Filtra marcadores de lista inúteis e concatena
				desc := strings.TrimPrefix(trimmed, "-")
				desc = strings.TrimSpace(desc)
				if desc != "" {
					if currentExp.Description != "" { currentExp.Description += " " }
					currentExp.Description += desc
				}
			}
		}
	}
	if currentExp != nil { data.Experience = append(data.Experience, *currentExp) }
}

func extractProjects(data *ExtractedData, lines []string) {
	var currentProj *models.Project
	reURL := regexp.MustCompile(`(https?://[^\s]+|www\.[^\s]+)`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || trimmed == "-" { continue }

		if len(trimmed) < 80 && !reURL.MatchString(trimmed) && (currentProj == nil || len(trimmed) < 40) {
			if currentProj != nil { data.Projects = append(data.Projects, *currentProj) }
			currentProj = &models.Project{ProfileID: 1, Name: trimmed}
		} else if currentProj != nil {
			if reURL.MatchString(trimmed) {
				currentProj.URL = reURL.FindString(trimmed)
			} else {
				desc := strings.TrimPrefix(trimmed, "-")
				desc = strings.TrimSpace(desc)
				if desc != "" {
					if currentProj.Description != "" { currentProj.Description += " " }
					currentProj.Description += desc
				}
			}
		}
	}
	if currentProj != nil { data.Projects = append(data.Projects, *currentProj) }
}
