// Code generated by pegomock. DO NOT EDIT.
package matchers

import (
	"github.com/petergtz/pegomock"
	"reflect"

	github "github.com/google/go-github/v48/github"
)

func AnyPtrToGithubRepository() *github.Repository {
	pegomock.RegisterMatcher(pegomock.NewAnyMatcher(reflect.TypeOf((*(*github.Repository))(nil)).Elem()))
	var nullValue *github.Repository
	return nullValue
}

func EqPtrToGithubRepository(value *github.Repository) *github.Repository {
	pegomock.RegisterMatcher(&pegomock.EqMatcher{Value: value})
	var nullValue *github.Repository
	return nullValue
}

func NotEqPtrToGithubRepository(value *github.Repository) *github.Repository {
	pegomock.RegisterMatcher(&pegomock.NotEqMatcher{Value: value})
	var nullValue *github.Repository
	return nullValue
}

func PtrToGithubRepositoryThat(matcher pegomock.ArgumentMatcher) *github.Repository {
	pegomock.RegisterMatcher(matcher)
	var nullValue *github.Repository
	return nullValue
}
