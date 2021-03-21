package intitools

type Program struct {
	ProgramId               string `json:"programId"`
	Status                  int    `json:"status"`
	ConfidentialityLevel    int    `json:"confidentialityLevel"`
	CompanyHandle           string `json:"companyHandle"`
	CompanyName             string `json:"companyName"`
	CompanySustainable      bool   `json:"companySustainable"`
	Handle                  string `json:"handle"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	MinBounty               string `json:"minBounty"`
	MaxBounty               string `json:"maxBounty"`
	LogoId                  string `json:"logoId"`
	IdentityCheckedRequired bool   `json:"identityCheckedRequired"`
	AwardRep                bool   `json:"awardRep"`
	SkipTriage              bool   `json:"skipTriage"`
	View                    int    `json:"view"`
}
