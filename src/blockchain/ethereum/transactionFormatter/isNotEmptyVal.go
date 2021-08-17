package transactionFormatter

func IsNotEmptyVal(ethereumValue string) bool {
	return len(ethereumValue) >= 3 && ethereumValue != "0x0"
}
