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
	SortOrder                float64  `json:"sort_order"`
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

func (s *SpacesService) RetrieveSpace(id uint) (*Space, *Response, error) {
	path := fmt.Sprintf("spaces/%v", strconv.FormatUint(uint64(id), 10))

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

func (s *SpacesService) RetrieveListOfSpaces() ([]*Space, *Response, error) {
	path := "spaces"

	req, err := s.client.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	// TODO: additional fields in space object - boards, entity_uid, user_id, access_mod
	spaces := []*Space{}
	resp, err := s.client.Do(req, &spaces)
	if err != nil {
		return nil, resp, err
	}

	return spaces, resp, nil

}

type CreateSpaceOptions struct {
	Title      *string `json:"title,omitempty"`
	ExternalID *uint   `json:"external_id,omitempty"`
}

func (s *SpacesService) CreateSpace(opt *CreateSpaceOptions) (*Space, *Response, error) {
	path := "spaces"

	req, err := s.client.NewRequest(http.MethodPost, path, opt)
	if err != nil {
		return nil, nil, err
	}

	// TODO: additional fields in space object - users
	sp := new(Space)
	resp, err := s.client.Do(req, sp)
	if err != nil {
		return nil, resp, err
	}

	return sp, resp, nil
}

func (s *SpacesService) RemoveSpace(id uint) (*Response, error) {
	path := fmt.Sprintf("spaces/%v", strconv.FormatUint(uint64(id), 10))

	req, err := s.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

type UpdateSpaceOptions struct {
	Title      *string `json:"title,omitempty"`
	ExternalID *uint   `json:"external_id,omitempty"`
	// TODO: allowed_card_type_ids
}

func (s *SpacesService) UpdateSpace(id uint, opt *UpdateSpaceOptions) (*Space, *Response, error) {
	path := fmt.Sprintf("spaces/%v", strconv.FormatUint(uint64(id), 10))

	req, err := s.client.NewRequest(http.MethodPatch, path, opt)
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
