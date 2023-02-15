package response

import "net/http"

const (
	InvalidParametersErrorCode ErrorSectionCode = iota
	AlreadyExistsErrCode
	NotFoundErrorCode
	GenerateModelErrorCode
	GetFormFileErrorCode
	SaveFileErrorCode
	CheckDirErrorCode
)

var commonErrorList = []ServiceErrorInfo{
	{ErrorSectionCode: InvalidParametersErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusBadRequest, ErrorMessageFmt: "Invalid parameters!"}},
	{ErrorSectionCode: AlreadyExistsErrCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusConflict, ErrorMessageFmt: "The object already exists!"}},
	{ErrorSectionCode: NotFoundErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusBadRequest, ErrorMessageFmt: "The object is not found!"}},
	{ErrorSectionCode: GenerateModelErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusInternalServerError, ErrorMessageFmt: "Generate model error!"}},
	{ErrorSectionCode: GetFormFileErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusBadRequest, ErrorMessageFmt: "Get form file error!"}},
	{ErrorSectionCode: SaveFileErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusInternalServerError, ErrorMessageFmt: "Save file error!"}},
	{ErrorSectionCode: CheckDirErrorCode, ErrorBaseInfo: ErrorBaseInfo{Status: http.StatusInternalServerError, ErrorMessageFmt: "Check dir error!"}},
}

func init() {
	RegisterErrorInfo(CommonErrorSectionCode, commonErrorList...)
}
