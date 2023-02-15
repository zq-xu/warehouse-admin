package supplier

//import (
//	"net/http"
//	"zq-xu/warehouse-admin/pkg/restapi/response"
//)
//
//const (
//	AlreadyExistsErrCode = iota + 1
//
//)
//
//var errorList = []response.ServiceErrorInfo{
//	{ErrorSectionCode: AlreadyExistsErrCode, ErrorBaseInfo: response.ErrorBaseInfo{Status: http.StatusConflict, ErrorMessageFmt: "The supplier [%s] already exists!"}},
//}
//
//func init() {
//	response.RegisterServiceErrorInfo(response.SupplierErrorSectionCode, errorList...)
//}
//
//func NewSupplierError(errCode response.ErrorSectionCode, msg ...interface{}) *response.ErrorInfo {
//	return response.NewServiceError(response.TransServiceCode(response.CommonErrorSectionCode, errCode), msg...)
//}
