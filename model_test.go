package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	null "gopkg.in/guregu/null.v3"
)

func TestPostRecipeArg(t *testing.T) {
	testErrorCases := []struct {
		input PostRecipeArg
	}{
		{PostRecipeArg{null.StringFromPtr(nil), null.IntFrom(5), null.IntFrom(3), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFrom(5), null.IntFrom(3), null.BoolFromPtr(nil)}},
		{PostRecipeArg{null.StringFromPtr(nil), null.IntFrom(5), null.IntFrom(3), null.BoolFromPtr(nil)}},

		{PostRecipeArg{null.StringFromPtr(nil), null.IntFrom(5), null.IntFromPtr(nil), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFrom(5), null.IntFromPtr(nil), null.BoolFromPtr(nil)}},
		{PostRecipeArg{null.StringFromPtr(nil), null.IntFrom(5), null.IntFromPtr(nil), null.BoolFromPtr(nil)}},

		{PostRecipeArg{null.StringFromPtr(nil), null.IntFromPtr(nil), null.IntFrom(3), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFromPtr(nil), null.IntFrom(3), null.BoolFromPtr(nil)}},
		{PostRecipeArg{null.StringFromPtr(nil), null.IntFromPtr(nil), null.IntFrom(3), null.BoolFromPtr(nil)}},

		{PostRecipeArg{null.StringFromPtr(nil), null.IntFromPtr(nil), null.IntFromPtr(nil), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFromPtr(nil), null.IntFromPtr(nil), null.BoolFromPtr(nil)}},
		{PostRecipeArg{null.StringFromPtr(nil), null.IntFromPtr(nil), null.IntFromPtr(nil), null.BoolFromPtr(nil)}},

		{PostRecipeArg{null.StringFrom("name"), null.IntFrom(0), null.IntFrom(3), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFrom(-1), null.IntFrom(3), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFrom(-2), null.IntFrom(3), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFrom(5), null.IntFrom(0), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFrom(5), null.IntFrom(-1), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFrom(5), null.IntFrom(4), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFrom(5), null.IntFrom(5), null.BoolFrom(false)}},
	}
	for i, v := range testErrorCases {
		validate.Struct(v.input)
		assert.Error(t, validate.Struct(v.input), "Case [%d]", i)
	}

	testNoErrorCases := []struct {
		input PostRecipeArg
	}{
		{PostRecipeArg{null.StringFrom("name"), null.IntFrom(5), null.IntFrom(3), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFrom(5), null.IntFromPtr(nil), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFromPtr(nil), null.IntFrom(3), null.BoolFrom(false)}},
		{PostRecipeArg{null.StringFrom("name"), null.IntFromPtr(nil), null.IntFromPtr(nil), null.BoolFrom(false)}},
	}
	for i, v := range testNoErrorCases {
		validate.Struct(v.input)
		assert.NoError(t, validate.Struct(v.input), "Case [%d]", i)
	}
}

func TestPostRateRecipeArg(t *testing.T) {
	testErrorCases := []struct {
		input PostRateRecipeArg
	}{
		{PostRateRecipeArg{null.IntFromPtr(nil)}},
		{PostRateRecipeArg{null.IntFrom(-1)}},
		{PostRateRecipeArg{null.IntFrom(-0)}},
		{PostRateRecipeArg{null.IntFrom(6)}},
		{PostRateRecipeArg{null.IntFrom(7)}},
	}
	for i, v := range testErrorCases {
		validate.Struct(v.input)
		assert.Error(t, validate.Struct(v.input), "Case [%d]", i)
	}

	testNoErrorCases := []struct {
		input PostRateRecipeArg
	}{
		{PostRateRecipeArg{null.IntFrom(1)}},
		{PostRateRecipeArg{null.IntFrom(2)}},
		{PostRateRecipeArg{null.IntFrom(3)}},
		{PostRateRecipeArg{null.IntFrom(4)}},
		{PostRateRecipeArg{null.IntFrom(5)}},
	}
	for i, v := range testNoErrorCases {
		validate.Struct(v.input)
		assert.NoError(t, validate.Struct(v.input), "Case [%d]", i)
	}
}
