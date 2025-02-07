// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package services

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/itential/ipctl/pkg/client"
	"github.com/itential/ipctl/pkg/logger"
)

type GroupInheritedGroup struct {
	GroupId string `json:"roleId"`
}

type GroupInheritedRole struct {
	RoleId string `json:"roleId"`
}

type GroupMemberOf struct {
	AAAManaged bool   `json:"aaaManaged"`
	GroupId    string `json:"groupId"`
}

// Group represents an authorization group in Platform.
type Group struct {
	Id              string                   `json:"_id,omitempty"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description"`
	Provenance      string                   `json:"provenance"`
	Meta            map[string]interface{}   `json:"_meta,omitempty"`
	AssignedRoles   []map[string]interface{} `json:"assignedRoles"`
	Inactive        bool                     `json:"inactive"`
	InheritedGroups []GroupInheritedGroup    `json:"inheritedGroups,omitempty"`
	InheritedRoles  []GroupInheritedRole     `json:"inheritedRoles,omitempty"`
	MemberOf        []GroupMemberOf          `json:"memberOf"`
}

type GroupService struct {
	client *ServiceClient
}

func NewGroupService(iapClient client.Client) *GroupService {
	return &GroupService{client: NewServiceClient(iapClient)}
}

func NewGroup(name, desc string) Group {
	logger.Trace()

	return Group{
		Name:          name,
		Description:   desc,
		Provenance:    "Pronghorn",
		Inactive:      false,
		AssignedRoles: []map[string]interface{}{},
		MemberOf:      []GroupMemberOf{},
	}
}

// GetAll will retrieve all of the configured authorization groupds and return
// them to the calling function as an array of Group instances.  If there are
// no configured groups on the server, this function will return an empty
// array.
func (svc *GroupService) GetAll() ([]Group, error) {
	logger.Trace()

	type Response struct {
		Results []Group `json:"results"`
		Total   int     `json:"total"`
	}

	var res *Response

	if err := svc.client.Get("/authorization/groups", &res); err != nil {
		return nil, err
	}

	logger.Info("Found %v group(s)", res.Total)

	return res.Results, nil
}

// Ger will attempt to retrieve the group as specified by the id argument.  The
// id argument is the 12 digest hex unique identifier for the authorization
// group.  If the group does not exist, this function will return an error.
func (svc *GroupService) Get(id string) (*Group, error) {
	logger.Trace()

	var res *Group
	var uri = fmt.Sprintf("/authorization/groups/%s", id)

	if err := svc.client.Get(uri, &res); err != nil {
		return nil, err
	}

	return res, nil

}

// GetByName attempts to find a group based on the group name.  The requirement
// argument `name` should specify the group name to find.  This function will
// iterate through all of the configured authorization gropus and return one
// that matches the same name.  If the group does not exist, this function will
// return an error
func (svc *GroupService) GetByName(name string) (*Group, error) {
	logger.Trace()

	groups, err := svc.GetAll()
	if err != nil {
		return nil, err
	}

	var res *Group

	for _, ele := range groups {
		if ele.Name == name {
			res = &ele
			break
		}
	}

	if res == nil {
		return nil, errors.New("group does not exist")
	}

	return res, nil
}

// Create will create a new authorization group.  This function does not check
// if the group already exists.  If it does, this function will return an
// error.
func (svc *GroupService) Create(in Group) (*Group, error) {
	logger.Trace()

	body := map[string]interface{}{
		"name":          in.Name,
		"description":   in.Description,
		"provenance":    in.Provenance,
		"assignedRoles": in.AssignedRoles,
		"inactive":      in.Inactive,
		"memberOf":      in.MemberOf,
	}

	type Response struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Data    Group  `json:"data"`
	}

	var res Response

	if err := svc.client.PostRequest(&Request{
		uri:                "/authorization/groups",
		body:               map[string]interface{}{"group": body},
		expectedStatusCode: http.StatusOK,
	}, &res); err != nil {
		return nil, err
	}

	logger.Info(res.Message)

	return &res.Data, nil
}

// Delete accepts the unique identifier and will delete the group from the
// system.  If the specified group identifier does not exist on the system,
// this function will return an error.
func (svc *GroupService) Delete(id string) error {
	logger.Trace()
	return svc.client.Delete(fmt.Sprintf("/authorization/groups/%s", id))
}
