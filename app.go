package main

/*
 * Project: Monte Curriculum
 * Author: Erasmo Cardoso - Software Engineer | Electronics Techniciant
 */

import (
	"context"
	"fmt"
	"mont_curriculum/internal/database"
	"mont_curriculum/internal/exporter"
	"mont_curriculum/internal/models"
	"mont_curriculum/internal/parser"
	"os"
	"encoding/base64"
	"encoding/json"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	err := database.InitDB("curriculum.db")
	if err != nil {
		fmt.Printf("Erro ao inicializar banco de dados: %v\n", err)
	}

	// Garante que o perfil com ID 1 exista para evitar erros de chaves estrangeiras
	_, _ = database.DB.Exec("INSERT OR IGNORE INTO profiles (id, first_name) VALUES (1, 'Novo Perfil')")
}

// SaveProfile salva ou atualiza o perfil principal (ID=1)
func (a *App) SaveProfile(p models.Profile) error {
	p.ID = 1
	p.Objective = parser.Sanitize(p.Objective)
	query := `INSERT INTO profiles (id, first_name, last_name, email, phone, address, age, photo, objective, linkedin, github, website)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			  ON CONFLICT(id) DO UPDATE SET
			  first_name=excluded.first_name, last_name=excluded.last_name, email=excluded.email,
			  phone=excluded.phone, address=excluded.address, age=excluded.age, photo=excluded.photo,
			  objective=excluded.objective, linkedin=excluded.linkedin, github=excluded.github, website=excluded.website`
	
	_, err := database.DB.Exec(query, p.ID, p.FirstName, p.LastName, p.Email, p.Phone, p.Address, p.Age, p.Photo, p.Objective, p.LinkedIn, p.GitHub, p.Website)
	return err
}

// GetProfile recupera o perfil principal
func (a *App) GetProfile() (models.Profile, error) {
	var p models.Profile
	err := database.DB.QueryRow("SELECT id, first_name, last_name, email, phone, address, age, photo, objective, linkedin, github, website, created_at FROM profiles WHERE id = 1").
		Scan(&p.ID, &p.FirstName, &p.LastName, &p.Email, &p.Phone, &p.Address, &p.Age, &p.Photo, &p.Objective, &p.LinkedIn, &p.GitHub, &p.Website, &p.CreatedAt)
	
	if err != nil {
		return models.Profile{}, err
	}
	return p, nil
}

// SaveEducation salva ou atualiza uma formação
func (a *App) SaveEducation(edu models.Education) error {
	var query string
	if edu.ID > 0 {
		query = `UPDATE education SET institution=?, course=?, start_date=?, end_date=?, description=? WHERE id=? AND profile_id=1`
		_, err := database.DB.Exec(query, edu.Institution, edu.Course, edu.StartDate, edu.EndDate, parser.Sanitize(edu.Description), edu.ID)
		return err
	}
	query = `INSERT INTO education (profile_id, institution, course, start_date, end_date, description) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, 1, edu.Institution, edu.Course, edu.StartDate, edu.EndDate, parser.Sanitize(edu.Description))
	return err
}

func (a *App) DeleteEducation(id int) error {
	_, err := database.DB.Exec("DELETE FROM education WHERE id=? AND profile_id=1", id)
	return err
}

// GetEducation retorna todas as formações do perfil
func (a *App) GetEducation() ([]models.Education, error) {
	rows, err := database.DB.Query("SELECT id, institution, course, start_date, end_date, description FROM education WHERE profile_id = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var edus []models.Education
	for rows.Next() {
		var e models.Education
		if err := rows.Scan(&e.ID, &e.Institution, &e.Course, &e.StartDate, &e.EndDate, &e.Description); err != nil {
			return nil, err
		}
		edus = append(edus, e)
	}
	return edus, nil
}

// SaveExperience salva ou atualiza uma experiência
func (a *App) SaveExperience(exp models.Experience) error {
	var query string
	if exp.ID > 0 {
		query = `UPDATE experiences SET company=?, position=?, start_date=?, end_date=?, description=? WHERE id=? AND profile_id=1`
		_, err := database.DB.Exec(query, exp.Company, exp.Position, exp.StartDate, exp.EndDate, parser.Sanitize(exp.Description), exp.ID)
		return err
	}
	query = `INSERT INTO experiences (profile_id, company, position, start_date, end_date, description) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, 1, exp.Company, exp.Position, exp.StartDate, exp.EndDate, parser.Sanitize(exp.Description))
	return err
}

func (a *App) DeleteExperience(id int) error {
	_, err := database.DB.Exec("DELETE FROM experiences WHERE id=? AND profile_id=1", id)
	return err
}

// GetExperiences retorna todas as experiências do perfil
func (a *App) GetExperiences() ([]models.Experience, error) {
	rows, err := database.DB.Query("SELECT id, company, position, start_date, end_date, description FROM experiences WHERE profile_id = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exps []models.Experience
	for rows.Next() {
		var e models.Experience
		if err := rows.Scan(&e.ID, &e.Company, &e.Position, &e.StartDate, &e.EndDate, &e.Description); err != nil {
			return nil, err
		}
		exps = append(exps, e)
	}
	return exps, nil
}

// SaveProject salva ou atualiza um projeto
func (a *App) SaveProject(proj models.Project) error {
	var query string
	if proj.ID > 0 {
		query = `UPDATE projects SET name=?, description=?, url=? WHERE id=? AND profile_id=1`
		_, err := database.DB.Exec(query, proj.Name, parser.Sanitize(proj.Description), proj.URL, proj.ID)
		return err
	}
	query = `INSERT INTO projects (profile_id, name, description, url) VALUES (?, ?, ?, ?)`
	_, err := database.DB.Exec(query, 1, proj.Name, parser.Sanitize(proj.Description), proj.URL)
	return err
}

func (a *App) DeleteProject(id int) error {
	_, err := database.DB.Exec("DELETE FROM projects WHERE id=? AND profile_id=1", id)
	return err
}

// GetProjects retorna todos os projetos do perfil
func (a *App) GetProjects() ([]models.Project, error) {
	rows, err := database.DB.Query("SELECT id, name, description, url FROM projects WHERE profile_id = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projs []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.URL); err != nil {
			return nil, err
		}
		projs = append(projs, p)
	}
	return projs, nil
}

// SelectPhoto abre um diálogo para o usuário escolher uma foto e retorna em Base64
func (a *App) SelectPhoto() (string, error) {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Selecione sua Foto",
		Filters: []runtime.FileFilter{
			{DisplayName: "Imagens", Pattern: "*.jpg;*.png;*.jpeg;*.webp"},
		},
	})

	if err != nil || selection == "" {
		return "", err
	}

	data, err := os.ReadFile(selection)
	if err != nil {
		return "", err
	}

	mimeType := "image/jpeg"
	if selection[len(selection)-3:] == "png" {
		mimeType = "image/png"
	} else if selection[len(selection)-4:] == "webp" {
		mimeType = "image/webp"
	}

	encoded := base64.StdEncoding.EncodeToString(data)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}

// ExportPDFFromData gera o PDF usando o estado atual enviado pelo frontend (Stateless)
func (a *App) ExportPDFFromData(p models.Profile, edus []models.Education, exps []models.Experience, projs []models.Project, template string) error {
	fmt.Printf("Gerando PDF em memória com template: %s\n", template)
	
	pdfData, err := exporter.GeneratePDF(p, edus, exps, projs, template)
	if err != nil {
		fmt.Printf("Erro ao gerar PDF: %v\n", err)
		return fmt.Errorf("falha ao gerar PDF: %v", err)
	}

	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Salvar Currículo",
		DefaultFilename: fmt.Sprintf("Curriculo_%s.pdf", p.FirstName),
		Filters: []runtime.FileFilter{
			{DisplayName: "Arquivos PDF", Pattern: "*.pdf"},
		},
	})

	if err != nil || savePath == "" {
		return err
	}

	err = os.WriteFile(savePath, pdfData, 0644)
	if err != nil {
		fmt.Printf("Erro ao escrever arquivo: %v\n", err)
		return err
	}

	fmt.Printf("PDF (Memória) exportado com sucesso para: %s\n", savePath)
	return nil
}

// ExportPDF gera o PDF (Versão Legada: busca no banco)
func (a *App) ExportPDF(template string) error {
	p, _ := a.GetProfile()
	edus, _ := a.GetEducation()
	exps, _ := a.GetExperiences()
	projs, _ := a.GetProjects()
	return a.ExportPDFFromData(p, edus, exps, projs, template)
}

// ImportResume permite selecionar um arquivo e extrair os dados completos
func (a *App) ImportResume() (*parser.ExtractedData, error) {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Selecione um Currículo para Importar",
		Filters: []runtime.FileFilter{
			{DisplayName: "Documentos", Pattern: "*.pdf"},
		},
	})

	if err != nil || selection == "" {
		return nil, err
	}

	data, err := parser.ParseResume(selection)
	if err != nil {
		return nil, err
	}

	// 1. Limpar dados antigos para garantir importação limpa
	database.DB.Exec("DELETE FROM education WHERE profile_id = 1")
	database.DB.Exec("DELETE FROM experiences WHERE profile_id = 1")
	database.DB.Exec("DELETE FROM projects WHERE profile_id = 1")

	// 2. Salvar Perfil
	data.Profile.ID = 1
	a.SaveProfile(data.Profile)

	// 3. Salvar Formações
	for _, edu := range data.Education {
		edu.ProfileID = 1
		a.SaveEducation(edu)
	}

	// 4. Salvar Experiências
	for _, exp := range data.Experience {
		exp.ProfileID = 1
		a.SaveExperience(exp)
	}

	// 5. Salvar Projetos
	for _, proj := range data.Projects {
		proj.ProfileID = 1
		a.SaveProject(proj)
	}
	
	return data, nil
}

// ExportDOCX gera o arquivo Word e solicita ao usuário onde salvar
func (a *App) ExportDOCX() error {
	p, _ := a.GetProfile()
	edus, _ := a.GetEducation()
	exps, _ := a.GetExperiences()
	projs, _ := a.GetProjects()

	fmt.Println("Gerando DOCX...")
	docxData, err := exporter.GenerateDOCX(p, edus, exps, projs)
	if err != nil {
		fmt.Printf("Erro ao gerar DOCX: %v\n", err)
		return err
	}

	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Salvar Currículo (Word)",
		DefaultFilename: fmt.Sprintf("Curriculo_%s.docx", p.FirstName),
		Filters: []runtime.FileFilter{
			{DisplayName: "Arquivos Word", Pattern: "*.docx"},
		},
	})

	if err != nil || savePath == "" {
		return err
	}

	err = os.WriteFile(savePath, docxData, 0644)
	if err != nil {
		fmt.Printf("Erro ao escrever arquivo: %v\n", err)
		return err
	}

	fmt.Printf("DOCX exportado com sucesso para: %s\n", savePath)
	return nil
}

// GetSettings recupera as configurações globais (ID=1)
func (a *App) GetSettings() (models.Settings, error) {
	var s models.Settings
	var labelsJSON []byte
	var showPhotoInt int

	err := database.DB.QueryRow("SELECT id, language, template, show_photo, labels FROM settings WHERE id = 1").
		Scan(&s.ID, &s.Language, &s.Template, &showPhotoInt, &labelsJSON)

	if err != nil {
		// Se não existir, retorna padrão
		return models.Settings{ID: 1, Language: "pt", Template: "default", ShowPhoto: true}, nil
	}

	s.ShowPhoto = showPhotoInt == 1
	if len(labelsJSON) > 0 {
		_ = json.Unmarshal(labelsJSON, &s.Labels)
	}
	return s, nil
}

// SaveSettings salva ou atualiza as configurações globais
func (a *App) SaveSettings(s models.Settings) error {
	s.ID = 1
	labelsJSON, _ := json.Marshal(s.Labels)
	showPhotoInt := 0
	if s.ShowPhoto { showPhotoInt = 1 }

	query := `INSERT INTO settings (id, language, template, show_photo, labels)
			  VALUES (?, ?, ?, ?, ?)
			  ON CONFLICT(id) DO UPDATE SET
			  language=excluded.language, template=excluded.template,
			  show_photo=excluded.show_photo, labels=excluded.labels`
	
	_, err := database.DB.Exec(query, s.ID, s.Language, s.Template, showPhotoInt, labelsJSON)
	return err
}
