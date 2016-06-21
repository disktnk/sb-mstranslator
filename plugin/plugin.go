package plugin

import (
	"gopkg.in/sensorbee/sensorbee.v0/bql/udf"
	mstranslator "pfi/tanakad/sb-mstranslator"
)

func init() {
	udf.MustRegisterGlobalUDSCreator("mstranslate",
		udf.UDSCreatorFunc(mstranslator.NewState))
	udf.MustRegisterGlobalUDF("mstranslate",
		udf.MustConvertGeneric(mstranslator.Translate))
}
