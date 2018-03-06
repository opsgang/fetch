package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-version"
	"sort"
	"strings"
)

const INVALID_TAG_CONSTRAINT = `
The --tag value you entered is not a valid constraint expression.
See https://github.com/opsgang/fetch#version-constraint-operators for examples.

Underlying error message:
%s
`

func isTagConstraintSpecificTag(tagConstraint string) (bool, string) {
	if len(tagConstraint) > 0 {
		switch tagConstraint[0] {
		// Check for a tagConstraint '='
		case '=':
			return true, strings.TrimSpace(tagConstraint[1:])

		// Check for a tagConstraint without constraint specifier
		// Neither of '!=', '>', '>=', '<', '<=', '~>' is prefixed before tag
		case '>', '<', '!', '~':
			return false, tagConstraint

		default:
			return true, strings.TrimSpace(tagConstraint)
		}
	}
	return false, tagConstraint
}

func getLatestAcceptableTag(tagConstraint string, tags []string) (string, error) {
	var latestTag string

	if len(tags) == 0 {
		return latestTag, nil
	}

	// Sort all tags
	// Our use of the library go-version means that each tag will each be represented as a *version.Version
	versions := make([]*version.Version, len(tags))
	for i, tag := range tags {
		v, err := version.NewVersion(tag)
		if err != nil {
			return latestTag, err
		}

		versions[i] = v
	}
	sort.Sort(version.Collection(versions))

	// If the tag constraint is empty, set it to the latest tag
	if tagConstraint == "" {
		tagConstraint = versions[len(versions)-1].String()
	}

	// Find the latest version that matches the given tag constraint
	constraints, err := version.NewConstraint(tagConstraint)
	if err != nil {
		// more useful err msg if hashicorp consider a constraint invalid.
		if strings.Contains(err.Error(), "Malformed constraint") {
			err = fmt.Errorf(INVALID_TAG_CONSTRAINT, err.Error())
		}
		return latestTag, err
	}

	latestAcceptableVersion := versions[0]
	for _, version := range versions {
		if constraints.Check(version) && version.GreaterThan(latestAcceptableVersion) {
			latestAcceptableVersion = version
		}
	}

	// check constraint against latest acceptable version
	if ! constraints.Check(latestAcceptableVersion) {
		return latestTag, errors.New("Tag does not exist")
	}

	// The tag name may have started with a "v". If so, re-apply that string now
	for _, originalTagName := range tags {
		if strings.Contains(originalTagName, latestAcceptableVersion.String()) {
			latestTag = originalTagName
		}
	}

	return latestTag, nil
}
