package intitools

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
)

type Program struct {
	ProgramId               string                `json:"programId"`
	Status                  int                   `json:"status"`
	ConfidentialityLevel    int                   `json:"confidentialityLevel"`
	CompanyHandle           string                `json:"companyHandle"`
	CompanyName             string                `json:"companyName"`
	CompanySustainable      bool                  `json:"companySustainable"`
	Handle                  string                `json:"handle"`
	Name                    string                `json:"name"`
	Description             string                `json:"description"`
	MinBounty               string                `json:"minBounty"`
	MaxBounty               string                `json:"maxBounty"`
	LogoId                  string                `json:"logoId"`
	IdentityCheckedRequired bool                  `json:"identityCheckedRequired"`
	AwardRep                bool                  `json:"awardRep"`
	SkipTriage              bool                  `json:"skipTriage"`
	View                    int                   `json:"view"`
	OutScopes               []ProgramChanges      `json:"outOfScopes"`
	InScopes                []ProgramChanges      `json:"inScopes"`
	RulesOfEngagement       []ProgramRulesChanges `json:"rulesOfEngagements"`
	Faqs                    []ProgramChanges      `json:"faqs"`
	SeverityAssessments     []ProgramChanges      `json:"severityAssessments"`
	Domains                 []ProgramDomains      `json:"domains"`
}

type ProgramChanges struct {
	CreatedAt int64                 `json:"createdAt"`
	Content   ProgramChangesContent `json:"content"`
}

type ProgramChangesContent struct {
	Content string `json:"content"`
}

type ProgramRulesChanges struct {
	CreatedAt int64                      `json:"createdAt"`
	Content   ProgramRulesChangesContent `json:"content"`
}

type ProgramRulesChangesContent struct {
	Content ProgramChangesContentContent `json:"content"`
}

type ProgramChangesContentContent struct {
	Description string `json:"description"`
}

type ProgramDomains struct {
	CreatedAt int64                   `json:"createdAt"`
	Content   []ProgramDomainsContent `json:"content"`
}

type ProgramDomainsContent struct {
	Id           string `json:"id"`
	Type         int    `json:"type"`
	Endpoint     string `json:"endpoint"`
	BountyTierId int    `json:"bountyTierId"`
	Description  string `json:"description"`
}

func (c *Client) GetProgramContentDiff(a Activity, field string) string {

	ctx := c.HttpCtx
	apiURL := fmt.Sprintf("%s/core/researcher/programs/%s/%s", c.ApiURL,
		url.PathEscape(a.Companyhandle), url.PathEscape(a.Programhandle))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return ""
	}

	req = req.WithContext(ctx)

	res := Program{}

	if err := c.sendRequest(req, &res); err != nil {
		return ""
	}

	// Find content changes matching current activity

	var changes []ProgramChanges

	switch field {
	case "OutScopes":
		changes = res.OutScopes
	case "InScopes":
		changes = res.InScopes
	case "Faqs":
		changes = res.Faqs
	case "SeverityAssessments":
		changes = res.SeverityAssessments
	default:
		panic("Unknown field name")
	}

	activityIdx := 0
	activityCreated := a.CreatedAt / 1000 // Get rid of miliseconds

	for idx, chg := range changes {

		if (activityCreated) == chg.CreatedAt {
			activityIdx = idx
			break
		}

	}

	newContent := changes[activityIdx].Content.Content
	oldContent := ""

	newDate := time.Unix(int64(activityCreated), 0).String()
	oldDate := ""

	if activityIdx > 0 {
		oldContent = changes[activityIdx-1].Content.Content
		oldDate = time.Unix(int64(changes[activityIdx-1].CreatedAt), 0).String()
	}

	edits := myers.ComputeEdits(span.URIFromPath(oldDate), oldContent, newContent)
	if len(edits) == 0 {
		return "No changes"
	}
	content := fmt.Sprint(gotextdiff.ToUnified(newDate, oldDate, oldContent, edits))
	if len(content) > 1800 {
		content = "Message too long"
	}
	return content
}

func (c *Client) GetProgramRulesDiff(a Activity) string {

	ctx := c.HttpCtx
	apiURL := fmt.Sprintf("%s/core/researcher/programs/%s/%s", c.ApiURL,
		url.PathEscape(a.Companyhandle), url.PathEscape(a.Programhandle))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return ""
	}

	req = req.WithContext(ctx)

	res := Program{}

	if err := c.sendRequest(req, &res); err != nil {
		return ""
	}

	// Find content changes matching current activity

	changes := res.RulesOfEngagement

	activityIdx := 0
	activityCreated := a.CreatedAt / 1000 // Get rid of miliseconds

	for idx, chg := range changes {

		if (activityCreated) == chg.CreatedAt {
			activityIdx = idx
			break
		}

	}

	newContent := changes[activityIdx].Content.Content.Description
	oldContent := ""

	newDate := time.Unix(int64(activityCreated), 0).String()
	oldDate := ""

	if activityIdx > 0 {
		oldContent = changes[activityIdx-1].Content.Content.Description
		oldDate = time.Unix(int64(changes[activityIdx-1].CreatedAt), 0).String()
	}

	edits := myers.ComputeEdits(span.URIFromPath(oldDate), oldContent, newContent)
	if len(edits) == 0 {
		return "No changes"
	}
	content := fmt.Sprint(gotextdiff.ToUnified(newDate, oldDate, oldContent, edits))
	if len(content) > 1800 {
		content = "Message too long"
	}
	return content

}

func (c *Client) GetProgramDomainsDiff(a Activity) string {

	ctx := c.HttpCtx
	apiURL := fmt.Sprintf("%s/core/researcher/programs/%s/%s", c.ApiURL,
		url.PathEscape(a.Companyhandle), url.PathEscape(a.Programhandle))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return ""
	}

	req = req.WithContext(ctx)

	res := Program{}

	if err := c.sendRequest(req, &res); err != nil {
		return ""
	}

	// Find content changes matching current activity

	changes := res.Domains

	activityIdx := 0
	activityCreated := a.CreatedAt / 1000 // Get rid of miliseconds

	for idx, chg := range changes {

		if (activityCreated) == chg.CreatedAt {
			activityIdx = idx
			break
		}

	}

	prevProgramContent := changes[activityIdx-1].Content
	nextProgramContent := changes[activityIdx].Content
	newContent := ""
	newDate := time.Unix(int64(activityCreated), 0).String()
	oldDate := time.Unix(int64(changes[activityIdx-1].CreatedAt), 0).String()

	// Check all previous domains if something was removed in new domains
	for _, pCont := range prevProgramContent {
		old_id := pCont.Id

		found := false
		for _, newDom := range nextProgramContent {
			if newDom.Id == old_id {
				found = true
				break
			}
		}

		if found == false {
			newContent += fmt.Sprintf("\n`%s` (%s) was removed!\n", pCont.Endpoint, c.GetEndpointType(pCont.Type))
		}
	}

	// Check all new domains if something was added or updated
	for _, nCont := range nextProgramContent {
		new_id := nCont.Id

		found := false
		foundIdx := 0
		for idx, oldDom := range prevProgramContent {
			if oldDom.Id == new_id {
				found = true
				foundIdx = idx
			}
		}

		if found == false {
			newContent += fmt.Sprintf("\n`%s` (%s) was added with-in %s!\n", nCont.Endpoint, c.GetEndpointType(nCont.Type), c.GetEndpointTier(nCont.BountyTierId))

		} else {
			pCont := prevProgramContent[foundIdx]
			if pCont.Endpoint != nCont.Endpoint || pCont.Type != nCont.Type || pCont.BountyTierId != nCont.BountyTierId || pCont.Description != nCont.Description {
				newContent += fmt.Sprintf("\n`%s` (%s) was updated:\n", nCont.Endpoint, c.GetEndpointType(nCont.Type))
				if pCont.Endpoint != nCont.Endpoint {
					newContent += fmt.Sprintf(" - Endpoint: `%s` -> `%s`\n", pCont.Endpoint, nCont.Endpoint)
				}
				if pCont.Type != nCont.Type {
					newContent += fmt.Sprintf(" - Type: `%s` -> `%s`\n", c.GetEndpointType(pCont.Type), c.GetEndpointType(nCont.Type))
				}
				if pCont.BountyTierId != nCont.BountyTierId {
					newContent += fmt.Sprintf(" - Tier: `%s` -> `%s`\n", c.GetEndpointTier(pCont.BountyTierId), c.GetEndpointTier(nCont.BountyTierId))
				}
				if pCont.Description != nCont.Description {
					// Add newlines (gotextdiff complains about it)
					pCont.Description = fmt.Sprintf("%s\n", pCont.Description)
					nCont.Description = fmt.Sprintf("%s\n", nCont.Description)

					edits := myers.ComputeEdits(span.URIFromPath(oldDate), pCont.Description, nCont.Description)

					diff := fmt.Sprint(gotextdiff.ToUnified(oldDate, newDate, pCont.Description, edits))
					newContent += fmt.Sprintf(" - Description:\n```diff\n%s\n```\n", diff)
				}

			}
		}

	}
	if len(newContent) > 1800 {
		newContent = "Message too long"
	}
	return newContent

}
func (c *Client) GetEndpointType(typeId int) string {
	typeIds := []string{
		"Dummy",
		"URL",
		"Android",
		"iOS",
		"IpRange",
		"Device",
		"Other",
	}

	return typeIds[typeId]
}

func (c *Client) GetEndpointTier(tierId int) string {
	tierIds := []string{
		"Dummy",
		"No Bounty Tier",
		"Tier 3",
		"Tier 2",
		"Tier 1",
	}

	return tierIds[tierId]
}
