// Copyright 2022 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package github

import (
	"context"
	"fmt"
)

func (s *OrganizationsService) GetSecurityManagerRole(ctx context.Context, org string) (*CustomOrgRoles, error) {
	roles, _, err := s.ListRoles(ctx, org)
	if err != nil {
		return nil, err
	}

	for _, role := range roles.CustomRepoRoles {
		if *role.Name == "security_manager" {
			return role, nil
		}
	}

	return nil, fmt.Errorf("security manager role not found")
}

// ListSecurityManagerTeams lists all security manager teams for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/security-managers#list-security-manager-teams
//
//meta:operation GET /orgs/{org}/security-managers
func (s *OrganizationsService) ListSecurityManagerTeams(ctx context.Context, org string) ([]*Team, *Response, error) {
	securityManagerRole, err := s.GetSecurityManagerRole(ctx, org)
	if err != nil {
		return nil, nil, err
	}

	options := &ListOptions{PerPage: 100}
	securityManagerTeams := make([]*Team, 0)
	for {
		teams, resp, err := s.ListTeamsAssignedToOrgRole(ctx, org, securityManagerRole.GetID(), options)
		if err != nil {
			return nil, nil, err
		}

		securityManagerTeams = append(securityManagerTeams, teams...)
		if resp.NextPage == 0 {
			return securityManagerTeams, resp, nil
		}

		options.Page = resp.NextPage
	}
}

// AddSecurityManagerTeam adds a team to the list of security managers for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/security-managers#add-a-security-manager-team
//
//meta:operation PUT /orgs/{org}/security-managers/teams/{team_slug}
func (s *OrganizationsService) AddSecurityManagerTeam(ctx context.Context, org, team string) (*Response, error) {
	securityManagerRole, err := s.GetSecurityManagerRole(ctx, org)
	if err != nil {
		return nil, err
	}

	return s.AssignOrgRoleToTeam(ctx, org, team, securityManagerRole.GetID())
}

// RemoveSecurityManagerTeam removes a team from the list of security managers for an organization.
//
// GitHub API docs: https://docs.github.com/rest/orgs/security-managers#remove-a-security-manager-team
//
//meta:operation DELETE /orgs/{org}/security-managers/teams/{team_slug}
func (s *OrganizationsService) RemoveSecurityManagerTeam(ctx context.Context, org, team string) (*Response, error) {
	securityManagerRole, err := s.GetSecurityManagerRole(ctx, org)
	if err != nil {
		return nil, err
	}

	return s.RemoveOrgRoleFromTeam(ctx, org, team, securityManagerRole.GetID())
}
