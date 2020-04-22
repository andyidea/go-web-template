package constant

func TemplateCode(smsType string) string {
	switch smsType {
	case SmsTypeLogin:
		return "SMS_182925837"
	}

	return ""
}
