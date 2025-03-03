package responses

type FoundPlayerResponse struct {
	Id            int    `json:"id"`
	FullNameEn    string `json:"fullNameEn"`
	FullNameLocal string `json:"fullNameLocal"`
	BirthDate     string `json:"birthDate"`
}
