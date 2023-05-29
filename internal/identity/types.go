package identity

type Identity[T any, I comparable] func(t T) I
