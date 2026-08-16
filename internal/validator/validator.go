package validator

type Validator struct {
    Errors map[string][]string
}

func New() *Validator {
    return &Validator{Errors: map[string][]string{}}
}

func (v *Validator) AddError(key, message string) {
    v.Errors[key] = append(v.Errors[key], message)
}

func (v *Validator) Check(ok bool, key, message string) {
    if !ok {
        v.AddError(key, message)
    }
}

func (v *Validator) Valid() bool {
    return len(v.Errors) == 0
}

// helper validators
func ValidateContactEmail(v *Validator, email string) {
    if email == "" {
        v.AddError("email", "must be provided")
    }
}

func ValidateConsumer(v *Validator, c interface{}) bool {
    // placeholder: caller should populate errors via Check
    return v.Valid()
}
