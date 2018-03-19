package main

import (
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

const NO_VALID_TAG_FOUND = `
Error occurred computing tag that best satisfies version contraint expression:
%s
`

func (o *fetchOpts) tagToGet(tags []string) (tag string, err error) {
	specific, tag := isTagConstraintOrExactTag(o.tagConstraint)
	if !specific {
		// Find the specific release that matches the latest version constraint
		latestTag, err := bestFitTag(o.tagConstraint, tags)
		if err != nil {
			return tag, fmt.Errorf(NO_VALID_TAG_FOUND, err)
		}
		tag = latestTag

		fmt.Printf("Most suitable tag for constraint %s is %s\n", o.tagConstraint, tag)
	}

	return
}

func isTagConstraintOrExactTag(tagConstraint string) (bool, string) {
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

func bestFitTag(tagConstraint string, tags []string) (string, error) {
	var latestTag string

	if len(tags) == 0 {
		return latestTag, nil
	}

	// We use Hashicorp go-versions comparison to find tags that are
	// semantic version strings.
	var versions []*version.Version
	for _, tag := range tags {
		v, err := version.NewVersion(tag)
		if err != nil {
			if strings.Contains(err.Error(), "Malformed version") {
				continue // ignore tags that do not fit expected semver
			}
			return latestTag, err
		}

		versions = append(versions, v)
	}

	if len(versions) == 0 {
		return latestTag, fmt.Errorf("No valid git tags found")
	}
	// Sort all tags so that last is latest.
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

	bestFitVersion := versions[0]
	for _, version := range versions {
		if constraints.Check(version) && version.GreaterThan(bestFitVersion) {
			bestFitVersion = version
		}
	}

	// check constraint against latest acceptable version
	if !constraints.Check(bestFitVersion) {
		return latestTag, fmt.Errorf("No tag met constraint.\n")
	}

	// The tag name may have started with a "v". If so, re-apply that string now
	for _, originalTagName := range tags {
		if strings.Contains(originalTagName, bestFitVersion.String()) {
			latestTag = originalTagName
		}
	}

	return latestTag, nil
}
