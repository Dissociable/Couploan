package msg

import (
	"testing"
)

func TestMsg(t *testing.T) {
	// a := fiber.New()
	// sessStore := session.New()
	// sessStore.RegisterType(map[string][]string{})
	// a.Use(middleware.SessionToContext(sessStore))
	// a.Get(
	// 	"/*", func(ctx fiber.Ctx) error {
	// 		assertMsg := func(typ Type, message string) {
	// 			ret := Get(ctx, nil, typ)
	// 			require.Len(t, ret, 1)
	// 			assert.Equal(t, message, ret[0])
	// 			ret = Get(ctx, nil, typ)
	// 			require.Len(t, ret, 0)
	// 		}
	// 		{
	// 			// Set the session id to request header so session.Get() can function correctly
	// 			// This is only required in testing
	// 			cookieKey := sessStore.KeyLookup[strings.Index(sessStore.KeyLookup, ":")+1:]
	// 			ctx.Request().Header.SetCookie(cookieKey, util.GetSessionIDFromFiberCtx(ctx))
	// 		}
	//
	// 		text := "aaa"
	// 		Success(ctx, nil, text)
	// 		assertMsg(TypeSuccess, text)
	//
	// 		text = "bbb"
	// 		Info(ctx, nil, text)
	// 		assertMsg(TypeInfo, text)
	//
	// 		text = "ccc"
	// 		Danger(ctx, nil, text)
	// 		assertMsg(TypeDanger, text)
	//
	// 		text = "ddd"
	// 		Warning(ctx, nil, text)
	// 		assertMsg(TypeWarning, text)
	//
	// 		text = "eee"
	// 		Set(ctx, nil, TypeSuccess, text)
	// 		assertMsg(TypeSuccess, text)
	// 		return nil
	// 	},
	// )
	// _, _ = tests.NewContextTest(a, "/", nil)
}
