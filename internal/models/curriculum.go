package models

import "time"

// Profile representa os dados pessoais do usuário
type Profile struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Age       int       `json:"age"`
	Photo     string    `json:"photo"` // Base64 ou path local
	Objective string    `json:"objective"`
	LinkedIn  string    `json:"linkedin"`
	GitHub    string    `json:"github"`
	Website   string    `json:"website"`
	CreatedAt time.Time `json:"created_at"`
}

// Education representa o histórico acadêmico
type Education struct {
	ID          int64  `json:"id"`
	ProfileID   int64  `json:"profile_id"`
	Institution string `json:"institution"`
	Course      string `json:"course"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

// Experience representa a trajetória profissional
type Experience struct {
	ID          int64  `json:"id"`
	ProfileID   int64  `json:"profile_id"`
	Company     string `json:"company"`
	Position    string `json:"position"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Description string `json:"description"`
}

// Project representa projetos paralelos ou significativos
type Project struct {
	ID          int64  `json:"id"`
	ProfileID   int64  `json:"profile_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// Settings define as preferências de exportação e labels customizados
type Settings struct {
	ID        int64             `json:"id"`
	Language  string            `json:"language"` // "pt" ou "en"
	Template  string            `json:"template"`
	Labels    map[string]string `json:"labels"`     // Labels customizados para cada seção
	ShowPhoto bool              `json:"show_photo"`
}
