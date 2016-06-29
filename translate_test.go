package mstranslator

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/core"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"testing"
)

func TestNewState(t *testing.T) {
	Convey("Given a context", t, func() {
		ctx := &core.Context{}
		Convey("When create a state with only required parameters", func() {
			params := data.Map{
				"client_id":     data.String("baz"),
				"client_secret": data.String("boo"),
			}
			s, err := NewState(ctx, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Terminate(ctx)
			})
			Convey("Then the state should create with default parameters", func() {
				token, ok := s.(*accessToken)
				So(ok, ShouldBeTrue)
				So(token.clientID, ShouldEqual, "baz")
				So(token.clientSecret, ShouldEqual, "boo")
				So(token.scope, ShouldEqual, msTranslatorScope)
				So(token.grantType, ShouldEqual, msTranslatorGrantType)
				So(token.accessTokenURL, ShouldEqual, msTranslatorAccessTokenURL)
				So(token.translatorURL, ShouldEqual, msTranslatorURL)
			})
		})

		Convey("When create a state with customized parameters", func() {
			params := data.Map{
				"client_id":        data.String("baz"),
				"client_secret":    data.String("boo"),
				"scope":            data.String("scope"),
				"access_token_url": data.String("access"),
				"grant_type":       data.String("cli"),
				"translator_url":   data.String("translation"),
			}
			s, err := NewState(ctx, params)
			So(err, ShouldBeNil)
			Reset(func() {
				s.Terminate(ctx)
			})
			Convey("Then the state should create with default parameters", func() {
				token, ok := s.(*accessToken)
				So(ok, ShouldBeTrue)
				So(token.clientID, ShouldEqual, "baz")
				So(token.clientSecret, ShouldEqual, "boo")
				So(token.scope, ShouldEqual, "scope")
				So(token.grantType, ShouldEqual, "cli")
				So(token.accessTokenURL, ShouldEqual, "access")
				So(token.translatorURL, ShouldEqual, "translation")
			})
		})
	})
}

func TestNewStateWithError(t *testing.T) {
	Convey("Given a context", t, func() {
		ctx := &core.Context{}
		Convey("When create a state without client ID", func() {
			params := data.Map{
				"client_secret": data.String("boo"),
			}
			Convey("Then an error should be occurred", func() {
				_, err := NewState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When create a state without client secret", func() {
			params := data.Map{
				"client_id": data.String("baz"),
			}
			Convey("Then an error should be occurred", func() {
				_, err := NewState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When create a state with error scope", func() {
			params := data.Map{
				"client_id":     data.String("baz"),
				"client_secret": data.String("boo"),
				"scope":         data.True,
			}
			Convey("Then an error should be occurred", func() {
				_, err := NewState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When create a state with error access_token_url", func() {
			params := data.Map{
				"client_id":        data.String("baz"),
				"client_secret":    data.String("boo"),
				"access_token_url": data.Int(55),
			}
			Convey("Then an error should be occurred", func() {
				_, err := NewState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When create a state with error grant_type", func() {
			params := data.Map{
				"client_id":     data.String("baz"),
				"client_secret": data.String("boo"),
				"grant_type":    data.Blob([]byte("byte")),
			}
			Convey("Then an error should be occurred", func() {
				_, err := NewState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When create a state with error translator_url", func() {
			params := data.Map{
				"client_id":      data.String("baz"),
				"client_secret":  data.String("boo"),
				"translator_url": data.Float(0.1),
			}
			Convey("Then an error should be occurred", func() {
				_, err := NewState(ctx, params)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
