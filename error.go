package tokbox

import "fmt"

var (
	//FIXME maybe there are more
	codeErrDict = map[int]string{
		403: "api_key or JWT token error ",
		404: "session does not exist ",
		409: "session not use Media Router or session is already being recorded ",
		500: "OpenTok server error ",
	}
)

type ResponseError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

func (this ResponseError) Error() string {
	return fmt.Sprintf("message:%s description:%s", this.Message, this.Description)
}
