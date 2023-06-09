package environment

// Deferred allows to delay intensive tasks after response has been sent through the API, to lower user waiting time
// (for example when sending emails).
type Deferred func() error
