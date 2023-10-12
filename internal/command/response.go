package command

// https://tex2e.github.io/rfc-translater/html/rfc5321.html#4-2--SMTP-Replies
const (
	// 正常系
	CodeHelp  = 214
	CodeGreet = 220
	CodeQuit  = 221
	CodeOk    = 250

	CodeStartInput = 354

	// Temporary Error
	CodeServiceNotAvailable = 421

	// Permanent Error
	CodeSyntaxError                = 500
	CodeArgumentSyntaxError        = 501
	CodeCommandNotImplemented      = 502
	CodeBadSequence                = 503
	CodeCommandParamNotImplemented = 504
	CodeAborted                    = 552
	CodeTransactionFail            = 554
	CodeOptionParamNotRecognized   = 555
)

const (
	// 正常系
	MsgHelp  = "Following commands are implemented"
	MsgGreet = "Service ready"
	MsgQuit  = "Service closing transmission channel"
	MsgOk    = "OK"

	MsgStartInput = "Start mail input; end with <CRLF>.<CRLF>"

	// Temporary Error
	MsgServiceNotAvailable = "Service not available, closing transmission channel"

	// Permanent Error
	MsgSyntaxError                = "Syntax error, command unrecognized"
	MsgArgumentSyntaxError        = "Syntax error in parameters or arguments"
	MsgBadSequence                = "Bad sequence of commands"
	MsgCommandParamNotImplemented = "Command parameter not implemented"
	MsgCommandNotImplemented      = "Command not implemented"
	MsgGreetFail                  = "No SMTP service here"
	MsgAborted                    = "Requested mail action aborted"
	MsgTransactionFail            = "Transaction failed"
	MsgOptionParamNotRecognized   = "Message size exceeds limit"
)
