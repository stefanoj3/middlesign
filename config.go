package middlesign

// MiddleSignConfig represents the middlesign configuration
type MiddleSignConfig struct {
	// TimestampKey is the key the the validator will read searching for the timestamp in the query parameters of the request
	TimestampKey string
	// TimestampFormat is the format used for the timestamps
	TimestampFormat string
	// SignatureKey is the key the the validator will read searching for the signature in the query parameters of the request
	SignatureKey string
	// Threshold is the amount of time for which we can consider the request valid expressed in seconds
	Threshold float64
	// Secret is the string used to sign/validate requests, only caller and server should know it
	// it should never be sent inside the request
	Secret string
}

// DefaultConfig returns a configuration with default values, only the secret must be defined by the user
func DefaultConfig(secret string) MiddleSignConfig {
	return MiddleSignConfig{
		TimestampKey:    "t",
		TimestampFormat: TimeFormatRFC3339,
		SignatureKey:    "sig",
		Threshold:       10,
		Secret:          secret,
	}
}
