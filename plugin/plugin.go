package plugin

import (
	mstranslator "github.com/disktnk/sb-mstranslator"
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
)

func init() {
	udf.MustRegisterGlobalUDSCreator("mstranslate",
		udf.UDSCreatorFunc(mstranslator.NewState))
	udf.MustRegisterGlobalUDF("mstranslate",
		udf.MustConvertGeneric(mstranslator.Translate))
}
