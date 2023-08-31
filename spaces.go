package kaiten

import (
	"fmt"
	"net/http"
	"strconv"
)

// IssueBoardsService handles communication with the spaces related
// methods of the Kaiten API.
type SpacesService struct {
	client *Client
}

type Space struct {
	ID                       uint     `json:"id"`
	UID                      string   `json:"uid"`
	Title                    string   `json:"title"`
	Updated                  string   `json:"updated"`
	Created                  string   `json:"created"`
	Archived                 bool     `json:"archived"`
	Access                   string   `json:"access"`
	ForeEveryoneAccessRoleID string   `json:"for_everyone_access_role_id"`
	EntityType               string   `json:"entity_type"`
	Path                     string   `json:"path"`
	SortOrder                uint     `json:"sort_order"`
	ParentEntityUID          string   `json:"parent_entity_uid"`
	CompanyID                uint     `json:"company_id"`
	AllowedCardTypeIDs       []string `json:"allowed_card_type_ids"`
	ExternalID               uint     `json:"external_id"`
	Settings                 *struct {
		Timeline *Timeline `json:"timeline"`
	} `json:"settings"`
}

type Timeline struct {
	StartHour            uint     `json:"startHour"`
	EndHour              uint     `json:"endHour"`
	WorkDays             []string `json:"workDays"`
	PlanningUnits        uint     `json:"planningUnits"`
	CalculateResourcesBy uint     `json:"calculateResourcesBy"`
}

func (s *SpacesService) RetrieveSpace(spaceID uint) (*Space, *Response, error) {
	path := fmt.Sprintf("/spaces/%v", strconv.FormatUint(uint64(spaceID), 10))

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	sp := new(Space)
	resp, err := s.client.Do(req, sp)
	if err != nil {
		return nil, resp, err
	}

	return sp, resp, nil

}
