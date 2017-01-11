package utils

func NewResponse(status bool, message string, response interface{}) map[string]interface{} {
	rsp := make(map[string]interface{})
	rsp["status"] = status
	rsp["message"] = message
	rsp["response"] = response
	//fmt.Printf("Redirect json response : status =%v\tmessage=%s\tresponse = %s\n", status, message, response)
	return rsp
}
