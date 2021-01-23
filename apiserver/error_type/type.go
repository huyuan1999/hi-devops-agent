package error_type

const (
	FloatingPointError    = "FloatingPointError"    // 浮点计算错误
	AssertionError        = "AssertionError"        // 断言失败
	EOFError              = "EOFError"              // EOF
	IOError               = "IOError"               // IO 错误
	OSError               = "OSError"               // 操作系统错误
	SyscallError          = "SyscallError"          // 操作系统错误
	IndexError            = "IndexError"            // 序列中没有此索引
	KeyError              = "KeyError"              // 映射中没有这个键
	NameError             = "NameError"             // 名称错误
	NotImplementedError   = "NotImplementedError"   // 尚未实现的方法
	TypeError             = "TypeError"             // 类型错误
	ValueError            = "ValueError"            // 值错误
	UnicodeDecodeError    = "UnicodeDecodeError"    // 解码时错误
	UnicodeEncodeError    = "UnicodeEncodeError"    // 编码时错误
	UnicodeTranslateError = "UnicodeTranslateError" // 转码时错误
	ZeroDivisionError     = "ZeroDivisionError"     // 除(或取模)零 (所有数据类型)
	HttpRequestError      = "HttpRequestError"      // http 请求错误
	SerializationError    = "SerializationError"    // 序列化错误
	Deserialization       = "Deserialization"       // 反序列化错误
	CertGenerateKeyError  = "CertGenerateKeyError"  // 生成证书错误
	CertGenerateCsrError  = "CertGenerateCsrError"  // 生成证书错误
	ListenerError         = "listenerError"         // 监听错误
	UnknownError          = "UnknownError"          // 未知错误
)
